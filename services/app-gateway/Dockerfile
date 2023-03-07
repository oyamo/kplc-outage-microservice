# Initial stage: download modules
FROM golang:1.18-alpine as golang-builder

RUN apk --update add git ca-certificates
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64;


FROM golang-builder AS app-builder

COPY ./gateway-service/go.mod ./gateway-service/go.sum /app/gateway-service/
COPY ./gateway-service/ /app/gateway-service/
COPY .  /app

WORKDIR /app
RUN go mod download

RUN go build -o app-gateway-bin app-gateway/cmd/main/main.go

FROM app-builder AS prepare-bin

COPY --from=app-builder /app/gateway-service-bin /usr/bin/gateway-service

ENTRYPOINT [" /usr/bin/gateway-service"]