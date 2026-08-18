package postgresql

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/yourname/smartmeal/service/handler"
	"github.com/yourname/smartmeal/storages/postgresql/schema/sqlc"
)

func (s *Storage) CreateDailyMeal(ctx context.Context, name string, calories int32) (handler.Meal, error) {
	row, err := s.q.CreateDailyMeal(ctx, sqlc.CreateDailyMealParams{
		Name:     name,
		Calories: calories,
	})
	if err != nil {
		return handler.Meal{}, err
	}
	return toMeal(row), nil
}

func (s *Storage) GetDailyMeal(ctx context.Context, id int64) (handler.Meal, error) {
	row, err := s.q.GetDailyMeal(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return handler.Meal{}, handler.ErrNotFound
		}
		return handler.Meal{}, err
	}
	return toMeal(row), nil
}

func (s *Storage) ListDailyMeals(ctx context.Context) ([]handler.Meal, error) {
	rows, err := s.q.ListDailyMeals(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]handler.Meal, 0, len(rows))
	for _, row := range rows {
		out = append(out, toMeal(row))
	}
	return out, nil
}

func toMeal(row sqlc.Meal) handler.Meal {
	return handler.Meal{
		ID:        row.ID,
		Name:      row.Name,
		Calories:  row.Calories,
		CreatedAt: row.CreatedAt.Time,
	}
}
