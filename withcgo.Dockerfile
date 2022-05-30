FROM golang:alpine3.14 as builder
#Run go env -w GOPROXY="https://goproxy.io,direct"
WORKDIR /go/src/vdev
COPY pht_sensor.go go.mod /go/src/vdev/
RUN go mod tidy && go build -ldflags="-w -s" pht_sensor.go

FROM scratch
COPY --from=builder /go/src/vdev/pht_sensor /app/pht_sensor
ENTRYPOINT ["/app/pht_sensor"]
