package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"math/rand"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Temperature struct {
	Value float64 `json:"temperature"`
	ID    string  `json:"device_id"`
}

type Humidity struct {
	Value float64 `json:"humidity"`
	ID    string  `json:"device_id"`
}

type Pressure struct {
	Value float64 `json:"pressure"`
	ID    string  `json:"device_id"`
}

type ValueConfig struct {
	Mean      float64 `json:"mean"`
	Deviation float64 `json:"deviation"`
}

func (v ValueConfig) Value() float64 {
	return v.Mean + rand.NormFloat64()*v.Deviation
}

type Options struct {
	Interval    time.Duration `json:"interval"`
	Count       int64         `json:"count"`
	Temperature ValueConfig   `json:"temperature"`
	Humidity    ValueConfig   `json:"humidity"`
	Pressure    ValueConfig   `json:"pressure"`
	Broker      string        `json:"broker"`
	DeviceId    string        `json:"device_id"`
	TempTopic   string        `json:"temperature_topic"`
	HumiTopic   string        `json:"humidity_topic"`
	PresTopic   string        `json:"pressure_topic"`
}

func publisher(options *Options) {
	rand.Seed(time.Now().UnixNano())
	opts := mqtt.NewClientOptions().AddBroker(options.Broker).SetClientID(options.DeviceId)
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	log.Println("Connected & Start publishing")
	var i int64
	c := options.Count
	i = 0
	for {
		if c > 0 {
			if i == c {
				break
			}
			i = i + 1
		}
		log.Println("Publishing......")
		x := options.Temperature.Value()
		t := &Temperature{x, options.DeviceId}
		buf, err := json.Marshal(t)
		if err == nil {
			token := client.Publish(options.TempTopic, byte(0), false, buf)
			token.Wait()
		}

		y := options.Humidity.Value()
		h := &Humidity{y, options.DeviceId}
		buf, err = json.Marshal(h)
		if err == nil {
			token := client.Publish(options.HumiTopic, byte(0), false, buf)
			token.Wait()
		}

		z := options.Pressure.Value()
		p := &Pressure{z, options.DeviceId}
		buf, err = json.Marshal(p)
		if err == nil {
			token := client.Publish(options.PresTopic, byte(0), false, buf)
			token.Wait()
		}
		time.Sleep((time.Millisecond * options.Interval))
	}

	client.Disconnect(250)
	log.Println("Done and Disconnected")
}

func parseOptions() (Options, error) {
	var opts Options

	opt_buf, err := ioutil.ReadFile("config/config.json")
	if err != nil {
		//log.Fatalln("fail to read config file")
		return opts, err
	}

	if err := json.Unmarshal(opt_buf, &opts); err != nil {
		//log.Fatalln("fail to parse config file")
		return opts, err
	}
	return opts, nil
}

func main() {
	log.SetPrefix("[PHT]")
	broker := flag.String("broker", "tcp://localhost:1883", "The broker URI. ex: tcp://10.10.1.1:1883")
	flag.Parse()

	var wg sync.WaitGroup

	options, err := parseOptions()
	if err != nil {
		log.Fatalln("Failed to get options: %s", err.Error())
		return
	}
	options.Broker = *broker

	wg.Add(1)
	go func() {
		defer wg.Done()
		publisher(&options)
	}()
	wg.Wait()
}
