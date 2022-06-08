#!/usr/bin/env bash
ocker run -d --rm -p 1883:1883 --name mosquitto --network mqtt --volume $PWD/mosquitto/config:/mosquitto/config:ro eclipse-mosquitto:latest
docker run -d --rm --network mqtt -v $PWD/config:/app/config:ro pht_sensor:0.9
