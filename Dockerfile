# Build stage
FROM golang:1.7.3 AS builder
RUN go get github.com/prometheus/node_exporter
COPY ./node-exporter/ /go/src/github.com/prometheus/node_exporter/collector
WORKDIR /go/src/github.com/prometheus/node_exporter
RUN make build

# the final image
FROM alpine:3.7
RUN apk add pciutils
COPY --from=builder /go/src/github.com/prometheus/node_exporter/node_exporter /bin/node_exporter
EXPOSE      9100
ENTRYPOINT  [ "/bin/node_exporter" ]

