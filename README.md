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
docker compose up db

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

