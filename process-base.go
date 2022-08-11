package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	rmq             *Rmq
	QueueOptions    QueueOptions
	ConsumerOptions ConsumeOptions
	BindingOptions  BindingOptions
	Consume         func(msg Any, delivery amqp.Delivery)
	CreateQueue     bool `default:"true"`
}
