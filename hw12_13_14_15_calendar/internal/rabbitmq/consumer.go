package rabbitmq

import (
	"github.com/EkaterinaShamanaeva/otus-go/hw12_13_14_15_calendar/internal/logger"
	"github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	*Session
}

func NewConsumer(addr, queue string, logger *logger.Logger) *Consumer {
	return &Consumer{New(addr, queue, logger)}
}

func (c *Consumer) Consume() (<-chan amqp091.Delivery, error) {
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
