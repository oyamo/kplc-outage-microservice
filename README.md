# kplc-outage-microservice

A heavily concurrent and scalable microservice written in Golang with gRPC and JSON transport for notifying customers - users - about planned outages in their regions

## Services
 - Notification Service - Collect customer details, queues and  send notifications
 - scrapper service  - Collect outage information and save to a database
 - Gateway service   - HTTP/1.1 REST service 

## Databases
- Redis - provide a cached copy of outages 
- MongoDB- high availability, scalable and variety of data types

## Message Queueing
- RabbitMQ - A reliable, scalable and efficient queueing software based on Advanced Message Queuing Protocol (AMQP),

## Setup 
### GRPC and Protobuffer package dependencies
```shell
cd proto
make
```
NOTE: You should add the protoc-gen-go-grpc to your PATH

```shell
PATH="${PATH}:${HOME}/go/bin"
```

### Compile docker images
```shell
make build
```

### Execute
```shell
make up
```

## References
 - Mario Carion, Domain Driven Design - [Youtube](https://www.youtube.com/watch?v=LUvid5TJ81Y)