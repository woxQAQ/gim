package mq

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/woxQAQ/gim/pkg/logger"
	"github.com/woxQAQ/gim/pkg/workerpool"
)

type RabbitMqConsumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	tag     string
	done    chan error
}

type RabbitMqConfig struct {
	URI          string
	ConnName     string
	Exchange     string
	ExchangeType string
	QueueName    string
	CTag         string
}

func NewRabbitMqConsumer(
	uri, connName, exchange, exchangeType, queueName, ctag string,
	l logger.Logger,
) (*RabbitMqConsumer, error) {
	c := &RabbitMqConsumer{
		conn:    nil,
		channel: nil,
		tag:     "",
		done:    make(chan error),
	}

	var err error
	amqp.SetLogger(l.With(logger.String("domain", "rabbitMq")))

	config := amqp.Config{Properties: amqp.NewConnectionProperties()}
	config.Properties.SetClientConnectionName(connName)
	c.conn, err = amqp.DialConfig(uri, config)
	if err != nil {
		return nil, err
	}

	c.channel, err = c.conn.Channel()
	if err != nil {
		return nil, err
	}
	if err = c.channel.ExchangeDeclare(
		exchange,
		exchangeType,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return nil, err
	}
	queue, err := c.channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	if err = c.channel.QueueBind(
		queue.Name,
		"",
		exchange,
		false,
		nil,
	); err != nil {
		return nil, err
	}
	l.Debug("RabbitMqConsumer: Binding queue to exchange",
		logger.String("queue", queue.Name),
		logger.String("exchange", exchange),
	)

	deliveries, err := c.channel.Consume(
		queue.Name,
		ctag,
		false,
		false,
		false, // noLocal
		false, // noWait
		nil,
	)
	if err != nil {
		return nil, err
	}

	workerpool.GetInstance().Submit(func() {
		defer func() {
			l.Debug("RabbitMqConsumer: Closing channel",
				logger.String("exchange", exchange),
				logger.String("queue", queue.Name),
				logger.String("ctag", ctag),
			)
			c.done <- nil
		}()
		for d := range deliveries {
			l.Debug("RabbitMqConsumer: Received a message",
				logger.String("exchange", exchange),
				logger.String("queue", queue.Name),
				logger.String("ctag", ctag),
				logger.String("message", string(d.Body)),
			)
			d.Ack(false)
		}
	})

	return c, nil
}
