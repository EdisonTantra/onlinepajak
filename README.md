# Project App

## Prerequisite

1. Golang  v1.23.0
2. Docker Compose v2.21.0

## Get Started

1. Using Docker Compose

```bash
# update postgresql host to `db` in config.yaml
docker compose up
```

2. Using Golang Binary

```bash
docker compose up -d db

# update postgresql host to `localhost` in config.yaml 
go run main.go http
```

or you could build the binary first with these commands

```bash

go build .
./lemonPajak http
```

### docker-compose command helper

This docker-compose load database schema from `./pkg/infra/psql/migrations/database.sql` 
If you want to reload new database schema, you need to execute this command below 
to remove the volume.
```bash      
docker-compose down --volumes 
```


## Unit Test

```bash
go test -v ./...
```

## TODO

1. accept JPG extension
2. QR Code Extraction
3. validation with more clear requirements
4. unittest
5. testing with more sample input
6. tracer
7. telemetry
8. remove psql