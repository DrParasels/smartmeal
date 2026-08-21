package natsout

import (
	"context"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	eventsv1 "github.com/yourname/smartmeal/events/proto/proto_gen"
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

func (p *Publisher) PublishMealCreated(ctx context.Context, meal handler.Meal) error {
	log := zerolog.Ctx(ctx)
	event := &eventsv1.MealCreated{
		MealId:        meal.ID,
		Calories:      meal.Calories,
		CreatedAtUnix: meal.CreatedAt.Unix(),
	}
	data, err := proto.Marshal(event)
	if err != nil {
		return err
	}
	if err := p.nc.Publish(p.subject, data); err != nil {
		return err
	}
	log.Debug().Int64("meal_id", meal.ID).Msg("meal event published")
	return p.nc.Publish(p.subject, data)
}
