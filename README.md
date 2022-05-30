Virtual MQTT client generating pressure, humidity and temperature data for testing purpose.

## Run MQTT broker

```bash
mkdir -p ./mosquitto/config
cat <<EOF | ./mosquitto/config/mosquitto.conf
> allow_anonymous true
> listener 1883
> EOF
docker network create mqtt
docker run -d --rm -p 1883:1883 --name mosquitto --network mqtt --volume ./mosquitto/config:/mosquitto/config:ro eclipse-mosquitto:latest
```

## Subscribe to broker
```bash
mosquitto_sub -t '#' -h localhost -p 1883
```

## Run PHT sensor
```bash
./pht_sensor
```

## Run PHT Sensor container
```bash
docker run --network mqtt pht_sensor:0.7 -broker tcp://mosquitto:1883
```
