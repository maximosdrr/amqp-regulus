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

func (this *Rmq) Connect(url string, exchanges []ExchangeOptions) *Rmq {

	conn, err := amqp.Dial(url)
	this.catchError(err, "Failed to connect to Rabbitmq")
	this.Connection = conn

	ch, err := this.Connection.Channel()
	this.catchError(err, "Failed to open channel")

	this.Channel = ch
	this.registerExchanges(exchanges)

	return this
}

func (this *Rmq) CloseConnection() {
	this.Connection.Close()
}

func (this *Rmq) CloseChannel() {
	this.Channel.Close()
}

func RegisterConsumer[T any](opts Consumer[T], connection *Rmq) {

	defaults.SetDefaults(&opts)

	if opts.CreateQueue {
		connection.createQueue(opts.QueueOptions)
	}

	if opts.BindingOptions.RoutingKey != "__NO_BIND__" && opts.BindingOptions.Exchange != "__NO_BIND__" {
		connection.Channel.QueueBind(opts.QueueOptions.QueueName, opts.BindingOptions.RoutingKey, opts.BindingOptions.Exchange, opts.BindingOptions.NoWait, nil)
	}

	consume(connection, opts.ConsumerOptions, opts.Consume)
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

func consume[T any](connection *Rmq, opts ConsumeOptions, callback func(msg T, delivery amqp.Delivery)) {
	defaults.SetDefaults(&opts)

	msgs, err := connection.Channel.Consume(opts.queueName, opts.Consumer, opts.AutoAck, opts.Exclusive, opts.NoLocal, opts.NoWait, nil)
	connection.catchError(err, "failed to register a consumer")

	go func() {
		for d := range msgs {
			var body T

			json.Unmarshal(d.Body, &body)
			callback(body, d)
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
