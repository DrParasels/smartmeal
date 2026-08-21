package handler

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	api "github.com/yourname/smartmeal/api/openapi/ogen"
)

type FakeStore struct {
	call bool
}

func (h *FakeStore) CreateDailyMeal(ctx context.Context, name string, calories int32) (Meal, error) {
	h.call = true
	return Meal{}, nil
}
func (h *FakeStore) GetDailyMeal(ctx context.Context, id int64) (Meal, error) {
	return Meal{}, nil
}
func (h *FakeStore) ListDailyMeals(ctx context.Context) ([]Meal, error) {
	return nil, nil
}

type FakePublisher struct {
	call bool
}

func (p *FakePublisher) PublishMealCreated(ctx context.Context, meal Meal) error {
	p.call = true
	return nil
}

func TestCreateDailyMeal_InputError(t *testing.T) {
	tests := []struct {
		name string
		req *api.CreateDailyMealRequest
	}{
		{
			name: "empty name",
			req: &api.CreateDailyMealRequest{Name: "", Calories: 100},
		},
		{
			name: "wrong calories",
			req: &api.CreateDailyMealRequest{Name: "Хлеб", Calories: 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &FakeStore{}
			publisher := &FakePublisher{}
			h := NewMealHandler(store, publisher)
			res, err := h.CreateDailyMeal(context.Background(), tt.req)
			require.NoError(t, err)
			require.False(t, store.call)
			require.False(t, publisher.call)
			require.IsType(t, &api.CreateDailyMealBadRequest{}, res)
		})
	}
	
}
