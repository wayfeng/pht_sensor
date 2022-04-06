From golang:alpine3.14 as builder
#Run go env -w GOPROXY="https://goproxy.io,direct"
Workdir /go/src/vdev
Add pht_sensor.go /go/src/vdev
Add go.mod /go/src/vdev
Run go mod tidy && go build -ldflags="-w -s" pht_sensor.go

From alpine:3.14
Copy --from=builder /go/src/vdev/pht_sensor /app/pht_sensor
ENTRYPOINT ["/app/pht_sensor"]
