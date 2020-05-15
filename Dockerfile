FROM golang:1.14.2 as builder

ENV LANG C.UTF-8
ENV TZ Asia/Tokyo
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
ENV APP_ROOT /usr/src/network-monitoring
RUN mkdir -p ${APP_ROOT}
WORKDIR ${APP_ROOT}

COPY go.mod ${APP_ROOT}
COPY go.sum ${APP_ROOT}
RUN go mod download
COPY . ${APP_ROOT}
RUN go build -o "bin/network-monitoring"

FROM alpine
RUN echo $APP_ROOT
COPY --from=builder /usr/src/network-monitoring/bin/network-monitoring /network-monitoring

