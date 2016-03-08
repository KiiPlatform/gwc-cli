package main

import (
	"fmt"
	"github.com/surgemq/message"
	"github.com/surgemq/surgemq/service"
)

func connectToLocalBroker() (*service.Client, error) {
	if lc != nil {
		return lc, nil
	}

	app := cc.Apps[*appName]
	lc = &service.Client{}
	msg := message.NewConnectMessage()
	msg.SetVersion(3)
	msg.SetCleanSession(true)
	msg.SetClientId([]byte(app.Site + "/" + app.ID + "/" + *endnodeVid))
	msg.SetKeepAlive(300)

	url := fmt.Sprintf("tcp://%s:%d", cc.GatewayAddress.Host, cc.GatewayAddress.Port)
	if err := lc.Connect(url, msg); err != nil {
		panic(err)
	}
	return lc, nil
}

func publishTopic(c *service.Client, topic string, payload string) error {
	pub := message.NewPublishMessage()
	pub.SetTopic([]byte(topic))
	pub.SetQoS(0)
	pub.SetPayload([]byte(payload))

	if err := c.Publish(pub, nil); err != nil {
		panic(err)
	}

	return nil
}
