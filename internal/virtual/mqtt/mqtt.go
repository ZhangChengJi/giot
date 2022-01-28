package mqtt

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Broker struct {
	Client mqtt.Client
	topic  string
}

var (
	Client mqtt.Client
)

func InitMqtt(c mqtt.Client) {
	Client = c
}

//Subscribe 订阅
//func (b *Broker) Subscribe(topic string) {
//	if token := b.Client.Subscribe(topic, 0, func(client mqtt.Client, message mqtt.Message) {
//
//	}); token.Wait() && token.Error() != nil {
//		fmt.Println(token.Error())
//		log.Error("subscribe topic:%s failed", topic)
//		return
//
//	}
//}

// Publish 发布
func (b *Broker) Publish(topic string, data interface{}) {
	token := b.Client.Publish(topic, 0, false, data)
	token.Wait()
}
