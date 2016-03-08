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
	lc *service.Client // client for local mqtt broker
	cc Config

	endnodeVid = flag.String("evid", "", "Vendor ThingID of Endnode(required).")
	appName    = flag.String("aname", "", "Name of te app configured in sample_confing.yml")
)

type Config struct {
	Apps           map[string]App `yaml:"apps"`
	GatewayAddress GatewayAddress `yaml:"gateway-address"`
}

type App struct {
	ID   string `yaml:"app-id"`
	Key  string `yaml:"app-key"`
	Site string `yaml:"app-site"`
}

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
	flag.Parse()

	if *endnodeVid == "" {
		fmt.Println("No vendorThingID of endnode found. Provide with -evid")
		return
	}
	if *appName == "" {
		fmt.Println("No appName is specified. Provide with -aname")
		return
	}
	_, err := connectToLocalBroker()
	if err != nil {
		fmt.Printf("fail to connect to local mqtt broker: %s\n", err)
		return
	}
	// subscribe to local mqtt broker after connect
	if err := subscribToReceiveCommand(); err != nil {
		fmt.Printf("fail to subscribe to local mqtt broker to receive command:%s\n", err)
		return
	}

	description :=
		` Please select a feature by input the following number:
		0. exit
		1. End Node Onboarding (end-node  → gateway)
		2. End Node State update (end-node → gateway ⇒ Kii Cloud)
		3. End Node Command result (end-node → gateway ⇒ Kii Cloud)
		4. Connect Endnode
		5. Disconnect Endnode
		`
	fmt.Printf("%s\n", description)

MainLoop:
	for {
		var n string
		fmt.Scanf("%s\n", &n)
		switch n {
		case "0":
			break MainLoop
		case "1":
			if err := onboardEndnode(); err != nil {
				fmt.Printf("fail to onboard endnode:%s\n", err)
			}
		case "2":
			if err := updateEndnodeState(); err != nil {
				fmt.Printf("fail to update endnode state:%s\n", err)
			}
		case "3":
			if err := publishCommandResults(); err != nil {
				fmt.Printf("fail to publish command results: %s\n", err)
			}
		case "4":
			if err := reportEndnodeStatus(true); err != nil {
				fmt.Println("fail to report online of endnode:", err)
			}
		case "5":
			if err := reportEndnodeStatus(false); err != nil {
				fmt.Println("fail to report offline of endnode:", err)
			}
		}
	}
}
