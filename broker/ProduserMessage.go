package broker

import (
	"FucknGO/config"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

type MailMessage struct {
	Name     string
	Email    string
	Password string
}

const sendMessage = "sendMessage"

func PublishMessage(message MailMessage) (err error) {

	conf, err := config.GetConfig()

	if err != nil {
		return err
	}
	//"amqp://guest:guest@localhost:5672"
	addressMQ := fmt.Sprintf("amqp://%s:%s@%s:%d",
		conf.JsonStr.RabbitMQ.User,
		conf.JsonStr.RabbitMQ.Password,
		conf.JsonStr.RabbitMQ.Address,
		conf.JsonStr.RabbitMQ.Port,
	)
	println(addressMQ)
	dial, err := amqp.Dial(addressMQ)

	if err != nil {
		return err
	}

	defer dial.Close()

	channel, err := dial.Channel()

	if err != nil {
		return err
	}

	defer channel.Close()

	queue, err := channel.QueueDeclare(sendMessage, true, false, false, false, nil)

	if err != nil {
		return err
	}

	body, err := json.Marshal(message)

	if err != nil {
		return err
	}

	err = channel.Publish(
		"", queue.Name,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         body,
		})

	if err != nil {
		return err
	}

	log.Printf("send mesage: %s + %s", message.Name, message.Email)
	return nil

}
