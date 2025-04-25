FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

COPY .env ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/app/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]