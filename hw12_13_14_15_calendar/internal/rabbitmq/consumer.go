package rabbitmq

import (
	"github.com/EkaterinaShamanaeva/otus-go/hw12_13_14_15_calendar/internal/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	*Session
}

func NewConsumer(addr, queue string, logger *logger.Logger) *Consumer {
	return &Consumer{New(addr, queue, logger)}
}

func (c *Consumer) Consume() (<-chan amqp.Delivery, error) {
	return c.channel.Consume(
		c.queue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
}
