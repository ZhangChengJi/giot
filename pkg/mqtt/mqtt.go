package mqtt

import (
	"fmt"
	"giot/conf"
	"giot/pkg/log"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"time"
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("接收消息: 从话题[ %s ] 发来的内容: %s \n", msg.Topic(), msg.Payload())
}

func New(conf *conf.Mqtt) (mqtt.Client, error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%v", conf.Host, conf.Port))
	opts.SetClientID("mqttx_2e81508d")
	opts.SetKeepAlive(5 * time.Second)
	//opts.SetDefaultPublishHandler(f)
	opts.SetPingTimeout(1 * time.Second)
	opts.SetCleanSession(false)
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.SetUsername(conf.Username)
	opts.SetPassword(conf.Password)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	opts.SetWill("offline", "go_mqtt_client offline", 1, false)

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		log.Error("mqtt connect failed:%s", token.Error())
		return nil, token.Error()
	}
	return c, nil

}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("mqtt连接成功...")
}

// 连接丢失的回调
var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("链接丢失: %v", err)
}
