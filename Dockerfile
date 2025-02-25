FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download
COPY ./src ./src

RUN CGO_ENABLED=0 GOOS=linux go build -o pocketbase ./src/main.go

FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/pocketbase /app/pocketbase

EXPOSE 8090

CMD ["/app/pocketbase", "serve", "--http=0.0.0.0:8090"] 