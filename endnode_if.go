package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/koron/go-dproxy"
	"github.com/surgemq/message"
)

func inputEndnodeState() (string, error) {
	var s string
	fmt.Println("Input state for endnode(should be json format):")
	fmt.Scanf("%s\n", &s)
	var v interface{}
	err := json.Unmarshal([]byte(s), &v)
	if err != nil {
		return "", fmt.Errorf("input state is not json format:%s\n", err)
	}
	_, err = dproxy.New(v).Map()
	if err != nil {
		return "", fmt.Errorf("input state is not json format:%s\n", err)
	}
	return s, nil
}

func inputCommandResults() (string, error) {
	var cr string
	fmt.Println("Input command result (should be json format, consist of commandID and actionResults):")
	fmt.Scanf("%s\n", &cr)

	var v interface{}
	err := json.Unmarshal([]byte(cr), &v)
	if err != nil {
		return "", err
	}
	_, err = dproxy.New(v).M("commandID").String()
	if err != nil {
		return "", errors.New("no commandID included in command result")
	}
	_, err = dproxy.New(v).M("actionResults").ProxySet().MapArray()
	if err != nil {
		return "", errors.New("no actionsResults included in command result")
	}
	return cr, nil
}

func onboardEndnode() error {

	app := cc.Apps[*appName]
	topic := app.Site + "/" + app.ID + "/e/" + *endnodeVid + "/states"
	payload := `{}`

	err := publishTopic(lc, topic, payload)
	if err != nil {
		return err
	}
	fmt.Printf("publish endnode state for onboarding:%s:%s\n", topic, payload)
	// wait 2 second for publishing to success
	time.Sleep(2 * time.Second)

	return nil
}

func updateEndnodeState() error {

	app := cc.Apps[*appName]
	topic := app.Site + "/" + app.ID + "/e/" + *endnodeVid + "/states"

	es, err := inputEndnodeState()
	if err != nil {
		return err
	}

	err = publishTopic(lc, topic, es)
	if err != nil {
		return err
	}
	fmt.Printf("publish endnode state:%s:%s\n", topic, es)
	// wait 2 second for publishing to success
	time.Sleep(2 * time.Second)
	return nil
}

func subscribToReceiveCommand() error {

	app := cc.Apps[*appName]
	topic := fmt.Sprintf("%s/%s/e/%s/commands", app.Site, app.ID, *endnodeVid)

	// subscrible a topic.
	sub := message.NewSubscribeMessage()
	sub.AddTopic([]byte(topic), 0)

	onRecv := func(m *message.PublishMessage) error {
		p := m.Payload()
		fmt.Printf("received command: \n%s\n", string(p))
		var v interface{}
		err := json.Unmarshal(p, &v)
		if err != nil {
			log.Println("failed to parse message :", err)
			return nil
		}
		id, err := dproxy.New(v).M("commandID").String()
		if err != nil {
			log.Println("failed to parse message :", err)
			return nil
		}
		if _, err := os.Stat("commands"); os.IsNotExist(err) {
			err = os.Mkdir("commands", 0777)
			if err != nil {
				log.Println("failed to create dir:", err)
				return nil
			}
		}
		ioutil.WriteFile("commands/"+id+".json", p, 0600)
		return nil
	}

	if err := lc.Subscribe(sub, nil, onRecv); err != nil {
		return err
	}
	return nil
}

func publishCommandResults() error {
	cr, err := inputCommandResults()
	if err != nil {
		return err
	}

	app := cc.Apps[*appName]
	topic := app.Site + "/" + app.ID + "/e/" + *endnodeVid + "/commandResults"
	err = publishTopic(lc, topic, cr)
	if err != nil {
		return err
	}
	fmt.Printf("publish endnode state:%s:%s\n", topic, cr)
	// wait 2 second for publishing to success
	time.Sleep(2 * time.Second)
	return nil

}

func reportConnectionStatus(online bool) error {
	app := cc.Apps[*appName]
	en := app.Site + "/" + app.ID + "/e/" + *endnodeVid
	if online {
		return publishTopic(lc, en+"/connect", "{}")
	}
	return publishTopic(lc, en+"/disconnect", "{}")
}
