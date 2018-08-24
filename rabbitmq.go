package golang_commons

import (
	"fmt"
	"github.com/streadway/amqp"
	"gogs.kopilka.kz/backend/golang_commons"
)

// FailOnError обоснован? или нет? с одной стороны - не может достучаться до рэббита - значит, не работает
func Publish(queue string, b []byte, rabbitMQUser, rabbitMQPassword, rabbitMQUrl string) {
	addr := fmt.Sprint("amqp://", rabbitMQUser, ":", rabbitMQPassword, "@", rabbitMQUrl)
	conn, err := amqp.Dial(addr)
	golang_commons.FailOnError(err, "Failed to connect to RabbitMQ")

	defer conn.Close()

	ch, err := conn.Channel()
	golang_commons.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(queue, true, false, false, false, nil)
	golang_commons.FailOnError(err, "Failed to declare a queue")

	err = ch.Publish("", q.Name, false, false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        b,
		})
	golang_commons.FailOnError(err, "Failed to publish a message")
}
