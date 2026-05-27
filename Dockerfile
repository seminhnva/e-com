# Stage 1: Build
FROM golang:1.25.6-alpine AS builder

WORKDIR /app

# Install swag to generate docs
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build
COPY . .
RUN swag init -g ./cmd/api/main.go -o ./docs
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api ./cmd/api

# Stage 2: Run
FROM scratch

WORKDIR /app

# Copy CA certificates for HTTPS
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/api .

EXPOSE 8080

CMD ["./api"]
