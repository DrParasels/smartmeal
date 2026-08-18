package postgresql

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/yourname/smartmeal/storages/postgresql/schema/sqlc"
)

func (s *Storage) AddMealCalories(ctx context.Context, createdAt time.Time, calories int32) error {
	_, err := s.q.UpsertDailyStat(ctx, sqlc.UpsertDailyStatParams{
		StatDate: pgtype.Date{
			Time:  createdAt.UTC(),
			Valid: true,
		},
		TotalCalories: calories,
	})
	return err
}
