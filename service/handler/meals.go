package handler

import (
	"context"
	"errors"
	"time"

	"github.com/rs/zerolog"
	api "github.com/yourname/smartmeal/api/openapi/ogen"
)

var ErrNotFound = errors.New("meal not found")

type Meal struct {
	ID        int64
	Name      string
	Calories  int32
	CreatedAt time.Time
}

type MealStore interface {
	CreateDailyMeal(ctx context.Context, name string, calories int32) (Meal, error)
	GetDailyMeal(ctx context.Context, id int64) (Meal, error)
	ListDailyMeals(ctx context.Context) ([]Meal, error)
}

type EventPublisher interface {
	PublishMealCreated(ctx context.Context, meal Meal) error
}

type MealHandler struct {
	store     MealStore
	publisher EventPublisher
}

func NewMealHandler(store MealStore, publisher EventPublisher) *MealHandler {
	return &MealHandler{store: store, publisher: publisher}
}

func toOgenMeal(m Meal) *api.Meal {
	return &api.Meal{
		ID:        m.ID,
		Name:      m.Name,
		Calories:  int(m.Calories),
		CreatedAt: m.CreatedAt,
	}
}

func (h *MealHandler) CreateDailyMeal(ctx context.Context, req *api.CreateDailyMealRequest) (api.CreateDailyMealRes, error) {
	log := zerolog.Ctx(ctx)

	if req.Name == "" || req.Calories <= 0 {
		log.Warn().Msg("invalid meal data")
		return &api.CreateDailyMealBadRequest{}, nil
	}

	meal, err := h.store.CreateDailyMeal(ctx, req.Name, int32(req.Calories))
	if err != nil {
		log.Error().Err(err).Msg("failed to create meal")
		return nil, err
	}
	if err := h.publisher.PublishMealCreated(ctx, meal); err != nil {
		log.Error().Err(err).Msg("failed to publish meal event")
		return nil, err
	}
	return toOgenMeal(meal), nil
}

func (h *MealHandler) GetDailyMeal(ctx context.Context, params api.GetDailyMealParams) (api.GetDailyMealRes, error) {
	log := zerolog.Ctx(ctx)

	meal, err := h.store.GetDailyMeal(ctx, params.MealId)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			log.Info().Int64("meal_id", params.MealId).Msg("meal not found")
			return &api.Error{
				Message: "meal not found",
			}, nil
		}
		log.Error().Err(err).Int64("meal_id", params.MealId).Msg("failed to get meal")
		return nil, err
	}
	return toOgenMeal(meal), nil
}

func (h *MealHandler) ListDailyMeals(ctx context.Context) (api.ListDailyMealsRes, error) {
	log := zerolog.Ctx(ctx)

	meals, err := h.store.ListDailyMeals(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to list meals")
		return nil, err
	}
	list := make(api.ListDailyMealsOKApplicationJSON, 0, len(meals))
	for _, m := range meals {
		list = append(list, *toOgenMeal(m))
	}
	return &list, nil
}
