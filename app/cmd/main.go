package main

import (
	"context"
	"db-delivery/config"
	"db-delivery/internal/bot/repository"
	"db-delivery/internal/bot/usecase"
	"db-delivery/internal/models"
	"db-delivery/pkg/logger"
	"db-delivery/pkg/storage/postgres"
	"db-delivery/pkg/storage/rabbit"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
)

func main() {
	viper, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := config.ParseConfig(viper)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Config loaded")

	ctx := context.Background()
	dependencies, err := initDependencies(ctx, cfg)
	if err != nil {
		log.Println(err.Error())
	}

	defer func(dependencies models.Dependencies) {
		if err := closeDependencies(dependencies); err != nil {
			dependencies.Logger.Errorf(err.Error())
		}
	}(dependencies)

	if err := mapHandler(ctx, cfg, dependencies); err != nil {
		dependencies.Logger.Errorf(err.Error())
	}

	exitCh := make(chan os.Signal, 1)
	signal.Notify(exitCh, os.Interrupt, syscall.SIGTERM)
	<-exitCh
}

func mapHandler(ctx context.Context, cfg *config.Config, dep models.Dependencies) (err error) {
	// repository
	botRepo := repository.NewBotPGRepo(dep.PgDB)

	// usecase
	botUC := usecase.NewBotUC(dep.Logger, cfg, botRepo, dep.RabbitMQ)

	go func() {
		if err := botUC.Consume(); err != nil {
			dep.Logger.Infof(err.Error())
		}
	}()

	return nil
}

func initDependencies(ctx context.Context, cfg *config.Config) (models.Dependencies, error) {
	logger := logger.InitLogger()

	pgDB, err := postgres.InitPgDB(ctx, cfg)
	if err != nil {
		return models.Dependencies{}, err
	} else {
		logger.Infof("PostgreSQL successful connection")
	}

	rabbitMQ, err := rabbit.InitRabbit(cfg)
	if err != nil {
		return models.Dependencies{}, err
	} else {
		logger.Infof("RabbitMQ successfil initialization")
	}

	return models.Dependencies{PgDB: pgDB, Logger: logger, RabbitMQ: rabbitMQ}, nil
}

func closeDependencies(dep models.Dependencies) error {
	if err := dep.PgDB.Close(); err != nil {
		return errors.Wrap(err, "PostgreSQL error close connection")
	} else {
		dep.Logger.Infof("PostgreSQL successful close connection")
	}

	if err := dep.RabbitMQ.Chann.Close(); err != nil {
		return errors.Wrap(err, "RabbitMQ error close chann")
	} else {
		dep.Logger.Infof("RabbitMQ successful close chann")
	}

	if err := dep.RabbitMQ.Conn.Close(); err != nil {
		return errors.Wrap(err, "RabbitMQ error close connection")
	} else {
		dep.Logger.Infof("RabbitMQ successful close connection")
	}

	return nil
}
