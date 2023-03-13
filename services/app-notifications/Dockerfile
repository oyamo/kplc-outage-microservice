# Initial stage: download modules
FROM golang:1.18-alpine as golang-builder

RUN apk --update add git ca-certificates
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64;


FROM golang-builder AS app-builder
WORKDIR /app/

COPY go.mod go.sum /app/
RUN go mod download && go mod verify

# Compile application
COPY . /app

RUN go build -o app-notifications-bin cmd/main/main.go

FROM app-builder AS prepare-bin

COPY --from=app-builder /app/app-notifications-bin /usr/bin/notification-service

ENTRYPOINT ["/usr/bin/notification-service"]