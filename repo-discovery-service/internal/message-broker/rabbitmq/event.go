package event

import (
	"repos-discovery-service/internal/constants"

	amqp "github.com/rabbitmq/amqp091-go"
)

func declareExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		constants.GITHUB_API_TOPIC, // name
		"topic",                    // type
		true,                       // durable?
		false,                      // auto-deleted?
		false,                      // internal?
		false,                      // no-wait?
		nil,                        // arguements?
	)
}

type Payload struct {
	Name string `json:"name"`
	Data any    `json:"data"`
}
