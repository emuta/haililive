package watcher

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/emuta/haililive/watcher/handler"
	"github.com/streadway/amqp"
)

type Watcher struct {
	url          string
	exchangeName string
	queueName    string
	handler      *handler.Handler
}

func NewWatcher(url, exchangeName, queueName string, handler *handler.Handler) *Watcher {
	return &Watcher{
		url:          url,
		exchangeName: exchangeName,
		queueName:    queueName,
		handler:      handler,
	}
}

func (w *Watcher) handleRabbitConnection(conn *amqp.Connection) {
	// Config
	log.WithFields(log.Fields{
		"Vhost":      conn.Config.Vhost,
		"ChannelMax": conn.Config.ChannelMax,
		"FrameSize":  conn.Config.FrameSize,
		"Heartbeat":  conn.Config.Heartbeat,
	}).Info("RabbitMQ Connection.Config")

	// Properties
	log.WithFields(log.Fields{
		"ClusterName":   conn.Properties["cluster_name"],
		"ServerVersion": conn.Properties["version"],
	}).Info("RabbitMQ Connection.Properties")

	// Properties.capabilities
	capabilities, _ := conn.Properties["capabilities"].(amqp.Table)
	log.WithFields(log.Fields(capabilities)).Info("RabbitMQ Connection.Properties[capabilities]")

	// monitor connection is close
	go func() {
		for {
			if conn.IsClosed() {
				log.Error("RabbitMQ Server closed, missing connection...")
			} else {
				log.Info("Heartbeat check RabbitMQ server connection: connected OK")
			}
			time.Sleep(30 * time.Second)
		}
	}()
}

func (w *Watcher) Run(ctx context.Context) {
	conn, err := amqp.Dial(w.url)
	if err != nil {
		log.WithError(err).Fatal("unable to create connection")
	}
	defer conn.Close()

	go w.handleRabbitConnection(conn)

	ch, err := conn.Channel()
	if err != nil {
		log.WithError(err).Fatal("unable to open new channel")
	}
	defer ch.Close()

	// declare exchange
	err = ch.ExchangeDeclare(
		w.exchangeName,      // exchange name
		amqp.ExchangeFanout, // exchange kind
		true,                // durable
		false,               // auto delete
		false,               // internal
		true,                // nowait
		nil,                 // args
	)
	if err != nil {
		log.WithFields(log.Fields{
			"exchange": w.exchangeName,
		}).WithError(err).Fatal("unable to open new channel")
	}

	queue, err := ch.QueueDeclare(
		w.queueName, // queue name
		false,       // durable
		true,        // delete when unused
		false,       // exclusive
		true,        // no wait
		nil,         // arguments
	)

	if err != nil {
		log.Fatal(err)
	}

	err = ch.QueueBind(
		queue.Name,     // queue name
		"",             // routing key
		w.exchangeName, // exchange name
		true,           // nowait
		nil,            // arguments
	)

	if err != nil {
		log.Fatal(err)
	}

	msgs, err := ch.Consume(
		queue.Name, // queue
		"",         // consumer
		true,       // auto ack
		false,      // exclusive
		false,      // no local
		false,      // no wait
		nil,        // arguments
	)

	for msg := range msgs {
		go w.handler.OnMessage(ctx, msg)
	}
}
