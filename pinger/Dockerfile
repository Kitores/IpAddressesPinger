FROM golang:alpine AS builder

WORKDIR /app
COPY . .
RUN apk add --no-cache make ca-certificates
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

FROM alpine
COPY --from=builder /app/main /usr/local/bin/main
RUN apk add docker
CMD ["main"]