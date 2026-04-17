FROM golang:1.26.2

RUN apt-get update && apt-get install -y vim && rm -rf /var/lib/apt/lists/*

RUN go install entgo.io/ent/cmd/ent@v0.14.6

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
