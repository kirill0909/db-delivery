package bot

import (
	"context"
	models "db-delivery/internal/models/bot"
)

type PgRepo interface {
	UserActivation(context.Context, models.UserActivationParams) error
}
