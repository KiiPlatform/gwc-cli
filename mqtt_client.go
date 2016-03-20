package main

import (
	"fmt"
	"github.com/surgemq/message"
	"github.com/surgemq/surgemq/service"
)

func connectToLocalBroker(c *service.Client, app App, converterID string, keepAlive uint16) error {
	msg := message.NewConnectMessage()
	msg.SetVersion(3)
	msg.SetCleanSession(true)
	msg.SetClientId([]byte(app.Site + "/" + app.ID + "/c/" + converterID))
	msg.SetKeepAlive(keepAlive)

	url := fmt.Sprintf("tcp://%s:%d", cc.GatewayAddress.Host, cc.GatewayAddress.Port)
	err := c.Connect(url, msg)
	if err != nil {
		return err
	}
	return nil
}

func publishTopic(c *service.Client, topic string, payload string) error {
	pub := message.NewPublishMessage()
	pub.SetTopic([]byte(topic))
	pub.SetQoS(0)
	pub.SetPayload([]byte(payload))

	err := c.Publish(pub, nil)
	if err != nil {
		return err
	}
	return nil
}
