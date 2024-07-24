package event

import (
	"commits-monitor-service/internal/constants"

	amqp "github.com/rabbitmq/amqp091-go"
)

func declareExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		constants.COMMITS_TOPIC, // name
		"topic",                 // type
		true,                    // durable?
		false,                   // auto-deleted?
		false,                   // internal?
		false,                   // no-wait?
		nil,                     // arguements?
	)
}
