# Initial stage: download modules
FROM golang:1.18-alpine as golang-builder

RUN apk --update add git ca-certificates
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64;


FROM golang-builder AS app-builder

COPY ./notification-service/go.mod ./notification-service/go.sum /app/notification-service/
COPY ./notification-service/ /app/notification-service/
COPY . /app

WORKDIR /app/notification-service

RUN go mod download && go mod verify

RUN go build -o app-notifications cmd/main/main.go

FROM app-builder AS prepare-bin

COPY --from=app-builder /app/notification-service/notification-service /usr/bin/notification-service

ENTRYPOINT [" /usr/bin/notification-service"]