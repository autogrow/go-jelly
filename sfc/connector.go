package sfc

import (
	"fmt"

	"github.com/nats-io/go-nats"
)

type intelli interface {
	Serial() string
	Update([]byte) error
}

// ConnectToNATS will connect an intelli object to the NATS server running on an IntelliLink device
func ConnectToNATS(i intelli, url string, timeout int) error {
	client, err := nats.Connect(url)

	if err != nil {
		return err
	}

	topic := fmt.Sprintf("intelli/%s", i.Serial())

	_, err = client.Subscribe(topic, func(msg *nats.Msg) {
		i.Update(msg.Data)
	})

	if err != nil {
		return err
	}

	return nil
}
