# Initial stage: download modules
FROM golang:1.18-alpine as golang-builder

RUN apk add build-base poppler-utils
RUN apk --update add git ca-certificates
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64;


FROM golang-builder AS app-builder

COPY ./scrapping-service/go.mod ./scrapping-service/go.sum /app/scrapping-service/
COPY ./scrapping-service/ /app/scrapping-service/
COPY .  /app

WORKDIR /app
RUN go mod download

RUN go build -o app-scrapper-bin app-scrapper/cmd/main/main.go

FROM app-builder AS prepare-bin

COPY --from=app-builder /app/scrapping-service-bin /usr/bin/scrapping-service

ENTRYPOINT [" /usr/bin/scrapping-service"]