package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"math/rand"
	"time"
)

type Tempature struct {
	ID    int     `json:"device_id"`
	Value float64 `json:"temperature"`
}

type Humidity struct {
	ID    int     `json:"device_id"`
	Value float64 `json:"humidity"`
}

type Pressure struct {
	ID    int     `json:"device_id"`
	Value float64 `json:"pressure"`
}

func publisher(opts *mqtt.ClientOptions) {
	rand.Seed(time.Now().UnixNano())
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	fmt.Println("Sample Publisher Started")
	for {
		//fmt.Println("---- doing publish ----")
		x := 26 + rand.NormFloat64()*1.0
		t := &Tempature{1, x}
		buf, err := json.Marshal(t)
		if err == nil {
			token := client.Publish("temperature/simulated/0", byte(0), false, buf)
			token.Wait()
		}

		y := 40 + rand.NormFloat64()*20.0
		h := &Humidity{1, y}
		buf, err = json.Marshal(h)
		if err == nil {
			token := client.Publish("humidity/simulated/0", byte(0), false, buf)
			token.Wait()
		}

		z := 1 + rand.NormFloat64()*.1
		p := &Pressure{1, z}
		buf, err = json.Marshal(p)
		if err == nil {
			token := client.Publish("pressure/simulated/0", byte(0), false, buf)
			token.Wait()
		}

		time.Sleep(time.Second * 5)
	}

	client.Disconnect(250)
	fmt.Println("Sample Publisher Disconnected")
}

func main() {
	broker := flag.String("broker", "tcp://ia_mqtt_broker:1883", "The broker URI. ex: tcp://10.10.1.1:1883")
	id := flag.String("id", "virtual_device_001", "The ClientID (optional)")
	flag.Parse()

	opts := mqtt.NewClientOptions().AddBroker(*broker).SetClientID(*id)

	go publisher(opts)
	select {}
}
