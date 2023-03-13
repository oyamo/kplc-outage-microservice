# Initial stage: download modules
FROM golang:1.18-alpine as golang-builder

RUN apk add build-base poppler-utils
RUN apk --update add git ca-certificates
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64;


FROM golang-builder AS app-builder


WORKDIR /app
# Copy go mod files
COPY  services/app-scrapper/go.mod services/app-scrapper/go.sum \
      /app/services/app-scrapper/

COPY services/app-gateway/go.mod services/app-gateway/go.sum \
      /app/services/app-gateway/

COPY services/app-notifications/go.mod services/app-notifications/go.sum \
      /app/services/app-notifications/

COPY proto/go.mod proto/go.sum \
     /app/proto/

COPY pkg/model/go.mod pkg/model/go.sum \
     /app/pkg/model/

COPY ./go.work /app


RUN go mod download

COPY .  /app


RUN go build -o app-scrapper-bin /app/services/app-scrapper/cmd/main/main.go

FROM app-builder AS prepare-bin

COPY --from=app-builder /app/app-scrapper-bin /usr/bin/scrapping-service

ENTRYPOINT ["/usr/bin/scrapping-service"]