FROM golang:1.14.2 as builder

ENV LANG C.UTF-8
ENV TZ Asia/Tokyo
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
ENV APP_ROOT /usr/src/gote
RUN mkdir -p ${APP_ROOT}
WORKDIR ${APP_ROOT}

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o "bin/gote"

FROM alpine
RUN apk --no-cache add tzdata && \
  cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime && \
  apk del tzdata
COPY --from=builder /usr/src/gote/bin/gote /gote

