package transfer

import (
	"fmt"
	broker "giot/internal/scheduler/mqtt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Transfer struct {
	mqtt     *broker.Broker
	dataChan chan string
}

func SetupTransfer() {
	t := &Transfer{
		mqtt: &broker.Broker{
			Client: broker.Client,
		},
	}
	t.ListenMqtt()
}

func (t *Transfer) ListenMqtt() {
	t.mqtt.Client.Subscribe("transfer/data/#", 0, func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("TOPIC: %s\n", msg.Topic())
		fmt.Printf("MSG: %s\n", msg.Payload())
	})
	t.mqtt.Client.Subscribe("transfer/alarm/#", 0, func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("TOPIC: %s\n", msg.Topic())
		fmt.Printf("MSG: %s\n", msg.Payload())
	})
}
