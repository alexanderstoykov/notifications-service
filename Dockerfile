FROM golang:1.20-alpine3.16 as base

ENV BASE_DIR /go/src/github.com/alexanderstoykov/notifications-service

WORKDIR ${BASE_DIR}

COPY go.mod go.sum ${BASE_DIR}/

RUN go mod download

COPY cmd ${BASE_DIR}/cmd
COPY internal ${BASE_DIR}/internal
COPY config ${BASE_DIR}/config

RUN CGO_ENABLED=0 GOOS=linux go build -v -o /dist/server ./cmd/server/main.go \
    && CGO_ENABLED=0 GOOS=linux go build -v -o /dist/cron ./cmd/cron/main.go

FROM alpine:3.16

RUN apk update \
    && apk add tzdata

COPY --from=base /dist .

EXPOSE 8080

CMD ["/server"]