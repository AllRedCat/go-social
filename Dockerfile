# Build
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git=2.54.0-r0 ca-certificates=20260611-r0 tzdata=2026b-r0

# Copy dependencies (enficienty cache)
COPY go.mod go.sum ./
RUN go mod download

# Copy code and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# Minimal Runtime
FROM scratch

# CA certificates
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

WORKDIR /root/

# Copy only binary
COPY --from=builder /app/main .

# Expose port
EXPOSE 8080

CMD ["./main"]
