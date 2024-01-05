package repository

import (
	"context"
	"db-delivery/internal/bot"
	models "db-delivery/internal/models/bot"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type BotPGRepo struct {
	db *sqlx.DB
}

func NewBotPGRepo(db *sqlx.DB) bot.PgRepo {
	return &BotPGRepo{db: db}
}

func (r *BotPGRepo) UserActivation(ctx context.Context, params models.UserActivationParams) error {
	res, err := r.db.ExecContext(ctx, queryUserActivation, params.ChatID, params.UUID, params.TgName)
	if err != nil {
		return errors.Wrapf(err, "BotPGRepo.UserActivation.queryUserActivation. params(%+v)", params)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return errors.Wrapf(err, "BotPGRepo.UserActivation.RowsAffected. params(%+v)", params)

	}

	if rowsAffected != 1 {
		return fmt.Errorf("BotPGRepo.UserActivation rowsAffected != 1. params(%+v)", params)
	}

	return nil
}
