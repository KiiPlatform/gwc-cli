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
	"github.com/surgemq/surgemq/service"
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

func inputJSONString() (string, error) {
	var s string
	fmt.Scanf("%s\n", &s)
	var v interface{}
	err := json.Unmarshal([]byte(s), &v)
	if err != nil {
		return "", fmt.Errorf("inputs is not json format:%s\n", err)
	}
	_, err = dproxy.New(v).Map()
	if err != nil {
		return "", fmt.Errorf("input is not json format:%s\n", err)
	}
	return s, nil
}

func readCommandResults() ([]byte, error) {
	b, err := ioutil.ReadFile("commandResult.json")
	if err != nil {
		return nil, err
	}

	// Validate contents.
	var v interface{}
	err = json.Unmarshal(b, &v)
	if err != nil {
		return nil, err
	}
	_, err = dproxy.New(v).M("commandID").String()
	if err != nil {
		return nil, errors.New("no commandID included in command result")
	}
	_, err = dproxy.New(v).M("actionResults").ProxySet().MapArray()
	if err != nil {
		return nil, errors.New("no actionsResults included in command result")
	}
	return b, nil
}

func onboardEndnode(c *service.Client) error {

	app := cc.Apps[*appName]
	topic := app.Site + "/" + app.ID + "/e/" + *endnodeVid + "/states"
	payload := `{}`

	err := publishTopic(c, topic, payload)
	if err != nil {
		return err
	}
	fmt.Printf("publish endnode state for onboarding:%s:%s\n", topic, payload)
	// wait 2 second for publishing to success
	time.Sleep(2 * time.Second)

	return nil
}

func updateEndnodeState(c *service.Client) error {

	app := cc.Apps[*appName]
	topic := app.Site + "/" + app.ID + "/e/" + *endnodeVid + "/states"

	es, err := inputEndnodeState()
	if err != nil {
		return err
	}

	err = publishTopic(c, topic, es)
	if err != nil {
		return err
	}
	fmt.Printf("publish endnode state:%s:%s\n", topic, es)
	// wait 2 second for publishing to success
	time.Sleep(2 * time.Second)
	return nil
}

func updateMultipleTraitState(c *service.Client) error {
	app := cc.Apps[*appName]
	topic := app.Site + "/" + app.ID + "/e/" + *endnodeVid + "/traitState"

	es, err := inputEndnodeState()
	if err != nil {
		return err
	}

	err = publishTopic(c, topic, es)
	if err != nil {
		return err
	}
	fmt.Printf("publish endnode trait state:%s:%s\n", topic, es)
	// wait 2 second for publishing to success
	time.Sleep(2 * time.Second)
	return nil
}

func updateSingleTraitState(c *service.Client) error {
	var alias string
	fmt.Println("Input alias string:")
	fmt.Scanf("%s\n", &alias)

	app := cc.Apps[*appName]
	topic := app.Site + "/" + app.ID + "/e/" + *endnodeVid + "/traitState/" + alias

	es, err := inputEndnodeState()
	if err != nil {
		return err
	}

	err = publishTopic(c, topic, es)
	if err != nil {
		return err
	}
	fmt.Printf("publish endnode trait state:%s:%s\n", topic, es)
	// wait 2 second for publishing to success
	time.Sleep(2 * time.Second)
	return nil
}

func subscribToReceiveCommand(c *service.Client) error {

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

	if err := c.Subscribe(sub, nil, onRecv); err != nil {
		return err
	}
	return nil
}

func publishCommandResults(c *service.Client) error {
	b, err := readCommandResults()
	if err != nil {
		return err
	}

	app := cc.Apps[*appName]
	topic := app.Site + "/" + app.ID + "/e/" + *endnodeVid + "/commandResults"
	err = publishTopic(c, topic, string(b))
	if err != nil {
		return err
	}
	fmt.Printf("publish endnode state:%s:%s\n", topic, string(b))
	// wait 2 second for publishing to success
	time.Sleep(2 * time.Second)
	return nil

}

func publishTraitCommandResults(c *service.Client) error {
	b, err := readCommandResults()
	if err != nil {
		return err
	}

	app := cc.Apps[*appName]
	topic := app.Site + "/" + app.ID + "/e/" + *endnodeVid + "/traitCmdResults"
	err = publishTopic(c, topic, string(b))
	if err != nil {
		return err
	}
	fmt.Printf("publish endnode state:%s:%s\n", topic, string(b))
	// wait 2 second for publishing to success
	time.Sleep(2 * time.Second)
	return nil

}

func reportConnectStatus(c *service.Client) error {
	app := cc.Apps[*appName]
	fmt.Println("Input thing properties for endnode(should be json format):")
	info, err := inputJSONString()
	if err != nil {
		return err
	}
	en := app.Site + "/" + app.ID + "/e/" + *endnodeVid
	return publishTopic(c, en+"/connect", info)
}

func reportDisconnectStatus(c *service.Client) error {
	app := cc.Apps[*appName]
	en := app.Site + "/" + app.ID + "/e/" + *endnodeVid
	return publishTopic(c, en+"/disconnect", "{}")
}
