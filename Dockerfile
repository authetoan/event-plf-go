# Build stage
FROM golang:1.22 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o booking-service ./cmd/main.go
###
FROM --platform=linux/amd64 alpine:3.11
WORKDIR /app
COPY --from=0 /app/booking-service .
COPY --from=0 /app/configs/config.yaml ./configs/config.yaml
EXPOSE 8080
CMD ["./booking-service"]
