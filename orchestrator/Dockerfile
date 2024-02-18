FROM golang:1.21

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .


# Server
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/app


CMD ["./server"]