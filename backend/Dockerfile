FROM golang:alpine AS builder

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

FROM alpine
COPY --from=builder /app/main /app/main
COPY /config/config.yaml app/config.yaml
CMD ["/app/main"]