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
		log.Printf("Dial: %s", err)
		return err
	}
	defer connection.Close()

	channel, err := connection.Channel()
	if err != nil {
		log.Printf("Channel: %s", err)
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

	if _, err := channel.QueueDeclare(
		"measurement", //name
		true,          // durable
		false,         //autodelete
		false,         //exclusive
		false,         //noWait
		nil,           //arguments
	); err != nil {
		log.Printf("Measurement Queue Declare: %s", err)
		return err
	}

	if err = channel.QueueBind(
		"measurement", //name
		"measurement", //routing key
		exchange,      //exchange
		false,         //noWait
		nil,           //arguments
	); err != nil {
		log.Printf("Measurement Queue Binding: %s", err)
		return err
	}

	if _, err = channel.QueueDeclare(
		"control", //name
		true,      // durable
		false,     //autodelete
		false,     //exclusive
		false,     //noWait
		nil,       //arguments
	); err != nil {
		log.Printf("Control Queue Declare: %s", err)
		return err
	}

	if err = channel.QueueBind(
		"control", //name
		"control", //routing key
		exchange,  //exchange
		false,     //noWait
		nil,       //arguments
	); err != nil {
		log.Printf("Control Queue Binding: %s", err)
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
			DeliveryMode:    amqp.Persistent, // 1=non-persistent, 2=persistent
			Priority:        0,               // 0-9
			// a bunch of application/implementation-specific fields
		},
	); err != nil {
		log.Printf("Exchange Publish: %s", err)
		return err
	}
	return nil
}
