package golang_commons

import (
	"github.com/streadway/amqp"
	"log"
)

type RMQHandlerFunc func(delivery amqp.Delivery) error

func Publish(queue string, b []byte, rabbitMQUser, rabbitMQPassword, rabbitMQUrl string) error {
	conn, err := amqp.Dial("amqp://" + rabbitMQUser + ":" + rabbitMQPassword + "@" + rabbitMQUrl)
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		return err
	}

	err = ch.Publish("", q.Name, false, false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        b,
		})
	if err != nil {
		return err
	}
	return nil
}

func ListenAndRecieve(queue string, handler RMQHandlerFunc, rabbitMQUser, rabbitMQPassword, rabbitMQUrl string) error {
	conn, err := amqp.Dial("amqp://" + rabbitMQUser + ":" + rabbitMQPassword + " + @" + rabbitMQUrl + "/")
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		queue, // name
		true,  // durable
		false, // delete when usused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			err = handler(d)
			if err != nil {
				log.Printf("Error: %s", err)
				ch.Reject(d.DeliveryTag, false)
			}

			ch.Ack(d.DeliveryTag, false)
		}
	}()
}
