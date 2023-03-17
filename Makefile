configure:
	protoc -I=proto  --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/power.proto
up:
	docker-compose -p kenya-power-outage up

build:
	docker-compose -p kenya-power-outage build

up-prod:
	docker-compose -p kenya-power-outage up