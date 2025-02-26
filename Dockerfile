FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS builder

WORKDIR /app

ARG TARGETPLATFORM
ARG BUILDPLATFORM
RUN echo "Building on $BUILDPLATFORM for $TARGETPLATFORM"

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download
COPY ./src ./src

RUN if [ "$TARGETPLATFORM" = "linux/amd64" ]; then \
    GOARCH=amd64; \
    elif [ "$TARGETPLATFORM" = "linux/arm64" ]; then \
    GOARCH=arm64; \
    else \
    GOARCH=amd64; \
    fi && \
    CGO_ENABLED=0 GOOS=linux GOARCH=$GOARCH go build -o pocketbase ./src/main.go

FROM --platform=$TARGETPLATFORM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/pocketbase /app/pocketbase

EXPOSE 8090

CMD ["/app/pocketbase", "serve", "--http=0.0.0.0:8090"]