package main

import (
	"encoding/json"
	"flag"
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

type PublishOptions struct {
	Interval  time.Duration //int64
	Count     int64
	TempBase  float64
	HumiBase  float64
	PresBase  float64
	TempDev   float64
	HumiDev   float64
	PresDev   float64
	DeviceId  string
	TempTopic string
	HumiTopic string
	PresTopic string
}

func publisher(opts *mqtt.ClientOptions, options *PublishOptions) {
	rand.Seed(time.Now().UnixNano())
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
		x := options.TempBase + rand.NormFloat64()*options.TempDev
		t := &Temperature{x, options.DeviceId}
		buf, err := json.Marshal(t)
		if err == nil {
			token := client.Publish(options.TempTopic, byte(0), false, buf)
			token.Wait()
		}

		y := options.HumiBase + rand.NormFloat64()*options.HumiDev
		h := &Humidity{y, options.DeviceId}
		buf, err = json.Marshal(h)
		if err == nil {
			token := client.Publish(options.HumiTopic, byte(0), false, buf)
			token.Wait()
		}

		z := options.PresBase + rand.NormFloat64()*options.PresDev
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

func main() {
	log.SetPrefix("[PHT]")
	broker := flag.String("broker", "tcp://localhost:1883", "The broker URI. ex: tcp://10.10.1.1:1883")
	id := flag.String("id", "virtual_device_001", "The ClientID (optional)")
	i := flag.Int64("i", 1000, "Publish data every # milliseconds")
	c := flag.Int64("count", 10, "Exit after publish #count times, 0 for infinite")
	tt := flag.String("tt", "temperature/simulated/0", "Tempature topic")
	ht := flag.String("ht", "humidity/simulated/0", "Humidity topic")
	pt := flag.String("pt", "pressure/simulated/0", "Pressure topic")
	tb := flag.Float64("tb", 26.0, "Temperature base value")
	td := flag.Float64("td", 2.0, "Temperature deviation value")
	hb := flag.Float64("hb", 40.0, "Humidity base value")
	hd := flag.Float64("hd", 2.0, "Humidity deviation value")
	pb := flag.Float64("pb", 1.0, "Pressure base value")
	pd := flag.Float64("pd", .1, "Pressure deviation value")
	flag.Parse()

	opts := mqtt.NewClientOptions().AddBroker(*broker).SetClientID(*id)

	var wg sync.WaitGroup

	options := PublishOptions{
		Interval:  time.Duration(*i),
		Count:     *c,
		TempBase:  *tb,
		HumiBase:  *hb,
		PresBase:  *pb,
		TempDev:   *td,
		HumiDev:   *hd,
		PresDev:   *pd,
		DeviceId:  *id,
		TempTopic: *tt,
		HumiTopic: *ht,
		PresTopic: *pt,
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		publisher(opts, &options)
	}()
	wg.Wait()
}
