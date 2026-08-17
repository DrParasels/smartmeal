package stats

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	eventsv1 "github.com/yourname/smartmeal/internal/pb/events/v1"
	"github.com/yourname/smartmeal/internal/storages/sqlc"
)

type MealStats struct {
	q *sqlc.Queries
}

func ConvertUnixToPgDate(unixSec int64) pgtype.Date {
	t := time.Unix(unixSec, 0).UTC()
	return pgtype.Date{
		Time: t,
		Valid: true,
	}
}

func NewMealStats(q *sqlc.Queries) *MealStats {
	return &MealStats{q: q}
}

func (h *MealStats) HandleMealCreated(ctx context.Context, ev *eventsv1.MealCreated) (sqlc.Stat, error) {
	row, err := h.q.UpsertDailyStat(ctx, sqlc.UpsertDailyStatParams{
		StatDate: ConvertUnixToPgDate(ev.CreatedAtUnix),
		TotalCalories: ev.Calories,
	})
	if err != nil {
		return sqlc.Stat{}, err
	}
	return row, nil
}