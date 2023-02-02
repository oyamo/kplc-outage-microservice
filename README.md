# kplc-outage-microservice

A microservice written in Golang with gRPC and JSON transport for notifying customers - users - about planned outages in their regions

## Services
 - Notification Service - Collect customer details and send notifications
 - scrapper service  - Collect outage information and save to a database
 - Gateway service   - HTTP/1 REST interface for listing outage information

## Databases
- Redis - provide a cached copy of outages 
- MySQL - provide permanent persistance

## Setup 
### GRPC and Protobuffer package dependencies
```shell
go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
```
NOTE: You should add the protoc-gen-go-grpc to your PATH

```shell
PATH="${PATH}:${HOME}/go/bin"
```

## References
 - Mario Carion, Domain Driven Design - [Youtube](https://www.youtube.com/watch?v=LUvid5TJ81Y)