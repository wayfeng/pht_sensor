From golang:alpine3.14 as builder
#Run go env -w GOPROXY="https://goproxy.io,direct"
Workdir /go/vdev
Add pht_sensor.go /go/vdev
Add go.mod /go/vdev
Run go mod tidy && CGO_ENABLED=0 go build -ldflags="-w -s" pht_sensor.go

From scratch
Copy --from=builder /go/vdev/pht_sensor /app/pht_sensor
ENTRYPOINT ["/app/pht_sensor"]
