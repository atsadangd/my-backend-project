# Stage 1: build
FROM golang:1.20-alpine AS builder
WORKDIR /app

# ensure certificates for HTTPS
RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
COPY . ./

ENV CGO_ENABLED=0
RUN go mod vendor
RUN go build -mod=vendor -o /server ./cmd/server

# Stage 2: runtime
FROM alpine:3.18
RUN apk add --no-cache ca-certificates
COPY --from=builder /server /server
EXPOSE 3000
ENTRYPOINT ["/server"]
