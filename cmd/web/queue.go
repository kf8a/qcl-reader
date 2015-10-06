package main

import (
	"github.com/streadway/amqp"
	"log"
)

var (
	uri          = "amqp://guest:guest@localhost:5672/"
	exchange     = "qcl"
	exchangeType = "direct"
	routingKey   = "measurement"
	reliable     = true
)

func publish(key string, message []byte) error {
	connection, err := amqp.Dial(uri)
	if err != nil {
		log.Println("Dial: %s", err)
		return err
	}
	defer connection.Close()

	channel, err := connection.Channel()
	if err != nil {
		log.Println("Channel: %s", err)
		return err
	}
	if err := channel.ExchangeDeclare(
		exchange,     // name
		exchangeType, // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // noWait
		nil,          // arguments
	); err != nil {
		log.Printf("Exchange Declare: %s", err)
		return err
	}

	if err = channel.Publish(
		exchange, // publish to an exchange
		key,
		// routingKey, // routing to 0 or more queues
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            message,
			DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
			Priority:        0,              // 0-9
			// a bunch of application/implementation-specific fields
		},
	); err != nil {
		log.Print("Exchange Publish: %s", err)
		return err
	}
	return nil
}
