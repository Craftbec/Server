FROM golang:latest

ENV PORT="8080"

ENV DB_HOST="host.docker.internal"
ENV DB_USER="postgres"
ENV DB_PASSWORD="postgres"
ENV DB_NAME="postgres"
ENV DB_PORT="5432"


RUN mkdir -p /app

COPY . /app

WORKDIR /app

RUN go build cmd/main.go

CMD ["./main"]