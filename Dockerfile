# Dockerfile definition for Backend application service.

# From which image we want to build. This is basically our environment.
FROM golang:1.23-alpine as Build

WORKDIR /app

# This will copy all the files in our repo to the inside the container at root location.
COPY . .

# Build our binary at root location.
RUN GOPATH= go build -o ./build/main main.go

####################################################################
# This is the actual image that we will be using in production.
FROM alpine:latest

WORKDIR /app

# We need to copy the binary from the build image to the production image.
COPY --from=Build /app/build/main .
COPY --from=Build /app/config ./config


# This is the port that our application will be listening on.
EXPOSE 8080

# This is the command that will be executed when the container is started.
ENTRYPOINT [ "/app/main" ]

CMD [ "http" ]
