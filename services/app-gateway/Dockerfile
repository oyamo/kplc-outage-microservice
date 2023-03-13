# Initial stage: download modules
FROM golang:1.18-alpine as golang-builder

RUN apk --update add git ca-certificates
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64;


FROM golang-builder AS app-builder
WORKDIR /app


COPY go.mod go.sum /app/
# Download modules
RUN go mod download

# Compile the app
COPY .  /app
RUN go build -o app-gateway-bin cmd/main/main.go

FROM app-builder AS prepare-bin

COPY --from=app-builder /app/app-gateway-bin /usr/bin/gateway-service

ENTRYPOINT ["/usr/bin/gateway-service"]