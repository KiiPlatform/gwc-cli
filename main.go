package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"

	"github.com/surgemq/surgemq/service"
)

var (
	cc Config

	endnodeVid = flag.String("evid", "", "Vendor ThingID of Endnode(required).")
	appName    = flag.String("aname", "", "Name of te app configured in sample_confing.yml")
)

// Config object
type Config struct {
	Apps           map[string]App `yaml:"apps"`
	GatewayAddress GatewayAddress `yaml:"gateway-address"`
	ConverterID    string         `yaml:"converter-id"`
	KeepAlive      uint16         `yaml:"keep-alive"`
}

// App represents Kii Cloud App.
type App struct {
	ID   string `yaml:"app-id"`
	Key  string `yaml:"app-key"`
	Site string `yaml:"app-site"`
}

// GatewayAddress repressents address of the Gateway.
type GatewayAddress struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func init() {
	b, err := ioutil.ReadFile("./config.yml")
	if err != nil {
		log.Fatalln("can't read ./config.yml file.")
	}
	err = yaml.Unmarshal(b, &cc)
	if err != nil {
		log.Fatalln("can't unmarshal ./config.yml")
	}
}

func main() {
	c := &service.Client{}
	flag.Parse()

	if *endnodeVid == "" {
		log.Fatalln("No vendorThingID of endnode found. Provide with -evid")
	}
	if *appName == "" {
		log.Fatalln("No appName is specified. Provide with -aname")
	}
	app := cc.Apps[*appName]
	err := connectToLocalBroker(c, app, cc.ConverterID, cc.KeepAlive)
	if err != nil {
		log.Fatalln("fail to connect to local mqtt broker: ", err)
	}
	// subscribe to local mqtt broker after connect
	if err := subscribToReceiveCommand(c); err != nil {
		log.Fatalln("fail to subscribe to local mqtt broker to receive command: ", err)
	}
	if err := onboardEndnode(c); err != nil {
		log.Fatalln("fail to onboard endnode: ", err)
	}
	description :=
		` Please select a feature by input the following number:
		0. exit
		1. End Node State update (end-node → gateway ⇒ Kii Cloud)
		2. End Node Command result (end-node → gateway ⇒ Kii Cloud)
		3. Connect Endnode
		4. Disconnect Endnode
		5. End Node multiple trait format State udpate
		6. End Node single trait format state udpate
		7. Update trait command result
		`
	log.Println(description)

MainLoop:
	for {
		var n string
		fmt.Scanf("%s\n", &n)
		switch n {
		case "0":
			break MainLoop
		case "1":
			if err := updateEndnodeState(c); err != nil {
				log.Println("fail to update endnode state: ", err)
			}
		case "2":
			if err := publishCommandResults(c); err != nil {
				log.Println("fail to publish command results: ", err)
			}
		case "3":
			if err := reportConnectStatus(c); err != nil {
				log.Println("fail to report online of endnode: ", err)
			}
		case "4":
			if err := reportDisconnectStatus(c); err != nil {
				log.Println("fail to report offline of endnode: ", err)
			}
		case "5":
			if err := updateMultipleTraitState(c); err != nil {
				log.Println("fail to update endnode multiple trait state: ", err)
			}
		case "6":
			if err := updateSingleTraitState(c); err != nil {
				log.Println("fail to update endnode single trait state: ", err)
			}
		case "7":
			if err := publishTraitCommandResults(c); err != nil {
				log.Println("fail to publish command results: ", err)
			}
		}
	}
}
