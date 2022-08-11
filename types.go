package rabbitmq

import amqp "github.com/rabbitmq/amqp091-go"

type Any interface{}

type ConsumeOptions struct {
	queueName string
	Consumer  string `default:""`
	AutoAck   bool   `default:"false"`
	Exclusive bool   `default:"false"`
	NoLocal   bool   `default:"false"`
	NoWait    bool   `default:"false"`
}

type QueueOptions struct {
	QueueName        string
	Durable          bool `default:"false"`
	DeleteWhenUnused bool `default:"true"`
	Exclusive        bool `default:"false"`
	NoWait           bool `default:"false"`
}

type ExchangeOptions struct {
	Name             string
	Type             string `default:"direct"`
	Durable          bool   `default:"false"`
	DeleteWhenUnused bool   `default:"false"`
	Internal         bool   `default:"false"`
	NoWait           bool   `default:"false"`
}

type BindingOptions struct {
	RoutingKey string `default:"__NO_BIND__"`
	Exchange   string `default:"__NO_BIND__"`
	NoWait     bool   `default:"false"`
}

type Publishing struct {
	amqp.Publishing
	Content any `default:"??@@__NO_CUSTOM_CONTENT__@@??"`
}

type PublishOptions struct {
	Queue     string
	Mandatory bool `default:"false"`
	Immediate bool `default:"false"`
	Content   Any
}

type PublishForExchange struct {
	Exchange   string
	RoutingKey string
	Mandatory  bool `default:"false"`
	Immediate  bool `default:"false"`
	Content    Any
}
