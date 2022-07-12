FROM golang:1.18.3 as builder

ENV LANG C.UTF-8
ENV TZ Asia/Tokyo
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
ENV APP_ROOT /gote
RUN mkdir -p ${APP_ROOT}
WORKDIR ${APP_ROOT}

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN make build

FROM ubuntu:jammy
RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y ca-certificates tzdata && \
    ln -sf /usr/share/zoneinfo/Asia/Tokyo /etc/localtime

COPY --from=builder /gote/bin/gote /gote
# have error occured 'fatal   network/network.go:49   ListenPacket    {"error": "listen ip4:icmp 0.0.0.0: socket: operation not permitted"}''
# USER 10000:10000


CMD ["/gote"]

