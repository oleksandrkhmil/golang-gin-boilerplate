# Start from golang base image
FROM golang:1.20.1-alpine3.17 as builder

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git build-base

# Set the current working directory inside the container
WORKDIR /app

RUN go install github.com/githubnemo/CompileDaemon@latest
RUN go install github.com/pressly/goose/cmd/goose@latest
RUN go install github.com/swaggo/swag/cmd/swag@latest

ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.7.3/wait /wait
RUN chmod +x /wait

#Command to run the executable
CMD swag init -g cmd/main.go \
  && /wait \
  && goose -dir "./db/migrations" ${DB_DRIVER} "${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}" up \
  && CompileDaemon --build="go build cmd/main.go" --command="./main" --color
