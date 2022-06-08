package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"math/rand"
	"sync"
	"time"

	"gopkg.in/yaml.v3"

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
	Interval    int64       `json:"interval" yaml:"interval"`
	Count       int64       `json:"count" yaml:"count"`
	Temperature ValueConfig `json:"temperature" yaml:"temperature"`
	Humidity    ValueConfig `json:"humidity" yaml:"humidity"`
	Pressure    ValueConfig `json:"pressure" yaml:"pressure"`
	Broker      string      `json:"broker" yaml:"broker"`
	DeviceId    string      `json:"device_id" yaml:"device_id"`
	TempTopic   string      `json:"temperature_topic" yaml:"temperature_topic"`
	HumiTopic   string      `json:"humidity_topic" yaml:"humidity_topic"`
	PresTopic   string      `json:"pressure_topic" yaml:"pressure_topic"`
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
		time.Sleep(time.Millisecond * time.Duration(options.Interval))
	}

	client.Disconnect(250)
	log.Println("Done and Disconnected")
}

func parseOptions(filename string) (Options, error) {
	var opts Options

	opt_buf, err := ioutil.ReadFile(filename)
	if err != nil {
		//log.Fatalln("fail to read config file")
		return opts, err
	}

	if err := yaml.Unmarshal(opt_buf, &opts); err != nil {
		//log.Fatalln("fail to parse config file")
		return opts, err
	}
	return opts, nil
}

func main() {
	log.SetPrefix("[PHT]")
	broker := flag.String("broker", "", "The broker URI. ex: tcp://10.10.1.1:1883")
	conf_file := flag.String("f", "config/config.yaml", "Configure file")
	flag.Parse()

	var wg sync.WaitGroup

	options, err := parseOptions(*conf_file)
	if err != nil {
		log.Fatalln("Failed to get options: %s", err.Error())
		return
	}
	if *broker != "" {
		options.Broker = *broker
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		publisher(&options)
	}()
	wg.Wait()
}
