package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer[T any] struct {
	rmq             *Rmq
	QueueOptions    QueueOptions
	ConsumerOptions ConsumeOptions
	BindingOptions  BindingOptions
	Consume         func(msg T, delivery amqp.Delivery)
	CreateQueue     bool `default:"true"`
}
