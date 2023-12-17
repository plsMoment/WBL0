package nats_streaming

import (
	"WBL0/internal/model"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/nats-io/stan.go"
	"log"
)

const (
	stanChannel = "channel"
	durableName = "last-position"
)

var ch chan model.Order

func stanReceiverHandle(msg *stan.Msg) {
	var validOrder model.ValidOrderTemplate
	defer func() {
		if err := msg.Ack(); err != nil {
			log.Printf("message acknowledge failed: %v", err)
			log.Printf("message index: %d", msg.Sequence)
		}
	}()

	if err := json.Unmarshal(msg.Data, &validOrder); err != nil {
		log.Println("invalid JSON")
		log.Printf("message index: %d", msg.Sequence)
		return
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(validOrder)
	if err != nil {
		log.Printf("validate error, message index: %d", msg.Sequence)
		for _, err := range err.(validator.ValidationErrors) {
			log.Println(err.Field())
		}
	}

	ch <- model.Order{Id: validOrder.OrderUID, Data: msg.Data}
}

// NewSubscriber create new subscribe to NATS cluster
func NewSubscriber(sc stan.Conn, orderCh chan model.Order) (stan.Subscription, error) {
	sub, err := sc.Subscribe(
		stanChannel,
		stanReceiverHandle,
		stan.DurableName(durableName),
		stan.SetManualAckMode(),
	)
	if err != nil {
		return nil, err
	}

	ch = orderCh

	return sub, nil
}
