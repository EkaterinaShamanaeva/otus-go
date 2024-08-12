package rabbitmq

import (
	"context"
	"github.com/EkaterinaShamanaeva/otus-go/hw12_13_14_15_calendar/internal/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	*Session
}

func NewProducer(addr, queue string, logger *logger.Logger) *Producer {
	return &Producer{New(addr, queue, logger)}
}

func (p *Producer) Publish(ctx context.Context, body []byte) error {
	return p.channel.PublishWithContext(
		ctx,
		"",
		p.queue,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
}
