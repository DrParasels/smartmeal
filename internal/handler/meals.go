package handler

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/nats-io/nats.go"
	api "github.com/yourname/smartmeal/api/ogen"
	eventsv1 "github.com/yourname/smartmeal/internal/pb/events/v1"
	"github.com/yourname/smartmeal/internal/storages/sqlc"
	"google.golang.org/protobuf/proto"
)

type MealHandler struct {
	q           *sqlc.Queries
	nc          *nats.Conn
	natsSubject string
}

func toOgenMeal(row sqlc.Meal) (*api.Meal) {
	return &api.Meal{
		ID: row.ID,
		Name: row.Name,
		Calories: int(row.Calories),
		CreatedAt: row.CreatedAt.Time,
	}
}

func NewMealHandler(q *sqlc.Queries, nc *nats.Conn, natsSubject string) *MealHandler {
	return &MealHandler{q: q, nc: nc, natsSubject: natsSubject}
}

func ConvertPgTimestampToUnix(pgTs pgtype.Timestamptz) int64 {
	if !pgTs.Valid {
		return 0 
	}
	return pgTs.Time.Unix()
}

func mealForStats(row sqlc.Meal) (*eventsv1.MealCreated) {
	return &eventsv1.MealCreated {
		MealId: row.ID,
		Calories: row.Calories,
		CreatedAtUnix: ConvertPgTimestampToUnix(row.CreatedAt),
	}

}

func (h *MealHandler) CreateDailyMeal(ctx context.Context, req *api.CreateDailyMealRequest) (api.CreateDailyMealRes, error) {
	if req.Name == "" {
        return &api.CreateDailyMealBadRequest{}, nil
    }
    if req.Calories <= 0 {
        return &api.CreateDailyMealBadRequest{}, nil
    }
	
	row, err := h.q.CreateDailyMeal(ctx, sqlc.CreateDailyMealParams{
		Name: req.Name,
		Calories: int32(req.Calories),
	})
	if err != nil {
		return nil, err
	}

	event := mealForStats(row)
	
	data, err := proto.Marshal(event)
	if err != nil {
		return nil, err
	}
	if err := h.nc.Publish(h.natsSubject, data); err != nil {
		return nil, err
	}

	return toOgenMeal(row), nil
}

func (h *MealHandler) GetDailyMeal(ctx context.Context, params api.GetDailyMealParams) (api.GetDailyMealRes, error) {
	row, err := h.q.GetDailyMeal(ctx, params.MealId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &api.Error{
				Message: "meal not found",
			}, nil
        }
        return nil, err
	}
	return toOgenMeal(row), nil
}

func (h *MealHandler) ListDailyMeals(ctx context.Context) (api.ListDailyMealsRes, error) {
	rows, err := h.q.ListDailyMeals(ctx)
	if err != nil {
		return &api.ListDailyMealsInternalServerError{}, nil
	}
	list := make(api.ListDailyMealsOKApplicationJSON, 0, len(rows))
	for _, v := range rows {
		list = append(list, *toOgenMeal(v))
	}
	return &list, nil
}

