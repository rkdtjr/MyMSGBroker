# Build Stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o broker main.go

# Final Stage
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/broker .
COPY --from=builder /app/web ./web
EXPOSE 8080 80
CMD ["./broker"]