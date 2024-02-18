FROM golang:1.21

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .


# Daemon
RUN CGO_ENABLED=1 GOOS=linux go build -o /daemon ./cmd/app


CMD ["./daemon"]