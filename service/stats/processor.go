package stats

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	eventsv1 "github.com/yourname/smartmeal/events"
)

type StatsStore interface {
	AddMealCalories(ctx context.Context, createdAt time.Time, calories int32) error
}

type MealStats struct {
	store StatsStore
}

func NewMealStats(store StatsStore) *MealStats {
	return &MealStats{store: store}
}

func (h *MealStats) HandleMealCreated(ctx context.Context, ev *eventsv1.MealCreated) error {
	log := zerolog.Ctx(ctx)
	createdAt := time.Unix(ev.CreatedAtUnix, 0).UTC()
	if err := h.store.AddMealCalories(ctx, createdAt, ev.Calories); err != nil {
		log.Error().Err(err).Msg("failed to add info at stats")
		return err
	}
	return nil
}
