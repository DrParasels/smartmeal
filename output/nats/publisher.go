package natsout

import (
	"context"

	"github.com/nats-io/nats.go"
	eventsv1 "github.com/yourname/smartmeal/events"
	"github.com/yourname/smartmeal/service/handler"
	"google.golang.org/protobuf/proto"
)

type Publisher struct {
	nc      *nats.Conn
	subject string
}

func NewPublisher(nc *nats.Conn, subject string) *Publisher {
	return &Publisher{nc: nc, subject: subject}
}

func (p *Publisher) PublishMealCreated(_ context.Context, meal handler.Meal) error {
	event := &eventsv1.MealCreated{
		MealId:        meal.ID,
		Calories:      meal.Calories,
		CreatedAtUnix: meal.CreatedAt.Unix(),
	}
	data, err := proto.Marshal(event)
	if err != nil {
		return err
	}
	return p.nc.Publish(p.subject, data)
}
