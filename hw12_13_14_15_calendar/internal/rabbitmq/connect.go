package rabbitmq

import (
	"github.com/EkaterinaShamanaeva/otus-go/hw12_13_14_15_calendar/internal/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Session struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	addr    string
	queue   string
	logger  *logger.Logger
}

func New(addr, queue string, logger *logger.Logger) *Session {
	return &Session{addr: addr, queue: queue, logger: logger}
}

func (s *Session) Connect() (err error) {
	// connection
	s.conn, err = amqp.Dial(s.addr)
	if err != nil {
		s.logger.Error("failed to connect RabbitMQ")
		return err
	}

	// open channel
	s.channel, err = s.conn.Channel()
	if err != nil {
		s.logger.Error("failed to open channel Rabbit MQ")
		return err
	}

	_, err = s.channel.QueueDeclare(
		s.queue,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		s.logger.Error("failed to queue declare")
		return err
	}
	return
}

func (s *Session) Close() error {
	err := s.channel.Close()
	if err != nil {
		s.logger.Error("failed to close channel RabbitMQ")
		return err
	}
	err = s.conn.Close()
	if err != nil {
		s.logger.Error("failed to close connection RabbitMQ")
		return err
	}
	return nil
}
