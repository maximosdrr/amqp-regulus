package rabbitmq

import (
	"encoding/json"
	"github.com/mcuadros/go-defaults"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type Rmq struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

func (this *Rmq) Connect(url string, consumers []Consumer, exchanges []ExchangeOptions) {

	conn, err := amqp.Dial(url)
	this.catchError(err, "Failed to connect to Rabbitmq")
	this.Connection = conn

	ch, err := this.Connection.Channel()
	this.catchError(err, "Failed to open channel")

	this.Channel = ch

	this.registerConsumers(consumers)
	this.registerExchanges(exchanges)
}

func (this *Rmq) CloseConnection() {
	this.Connection.Close()
}

func (this *Rmq) CloseChannel() {
	this.Channel.Close()
}

func (this *Rmq) registerConsumers(consumers []Consumer) {
	if consumers == nil {
		return
	}

	for _, consumer := range consumers {

		defaults.SetDefaults(&consumer)

		if consumer.CreateQueue {
			this.createQueue(consumer.QueueOptions)
		}

		if consumer.BindingOptions.RoutingKey != "__NO_BIND__" && consumer.BindingOptions.Exchange != "__NO_BIND__" {
			this.Channel.QueueBind(consumer.QueueOptions.QueueName, consumer.BindingOptions.RoutingKey, consumer.BindingOptions.Exchange, consumer.BindingOptions.NoWait, nil)
		}
		this.consume(consumer.ConsumerOptions, consumer.Consume)
	}
}

func (this *Rmq) catchError(err error, msg string) {
	if err != nil {
		log.Panicf("%s : %s", msg, err)
	}
}

func (this *Rmq) createQueue(opts QueueOptions) (amqp.Queue, error) {
	defaults.SetDefaults(&opts)

	msgs, err := this.Channel.QueueDeclare(opts.QueueName, opts.Durable, opts.DeleteWhenUnused, opts.Exclusive, opts.NoWait, nil)
	this.catchError(err, "Failed to declare a queue")

	return msgs, err
}

func (this *Rmq) consume(opts ConsumeOptions, callback func(msg Any, delivery amqp.Delivery)) {
	defaults.SetDefaults(&opts)

	msgs, err := this.Channel.Consume(opts.queueName, opts.Consumer, opts.AutoAck, opts.Exclusive, opts.NoLocal, opts.NoWait, nil)
	this.catchError(err, "failed to register a consumer")

	go func() {
		for d := range msgs {
			var body Any

			json.Unmarshal(d.Body, &body)
			callback(d.Body, d)
		}
	}()
}

func (this *Rmq) registerExchanges(opts []ExchangeOptions) {
	if opts == nil {
		return
	}

	for _, option := range opts {
		defaults.SetDefaults(&option)
		this.Channel.ExchangeDeclare(option.Name, option.Type, option.Durable, option.DeleteWhenUnused, option.Internal, option.NoWait, nil)
	}
}

func (this *Rmq) Publish(opts PublishOptions) {
	defaults.SetDefaults(&opts)
	data, err := json.Marshal(opts.Content)

	if err != nil {
		this.catchError(err, "Error caught when trying convert message to json")
	}

	content := amqp.Publishing{
		Body: data,
	}

	this.Channel.Publish("", opts.Queue, opts.Mandatory, opts.Immediate, content)
}

func (this *Rmq) PublishForExchange(opts PublishForExchange) {
	defaults.SetDefaults(&opts)
	data, err := json.Marshal(opts.Content)

	if err != nil {
		this.catchError(err, "Error caught when trying convert message to json")
	}

	content := amqp.Publishing{
		Body: data,
	}

	this.Channel.Publish(opts.Exchange, opts.RoutingKey, opts.Mandatory, opts.Immediate, content)
}
