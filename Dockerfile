# Stage 1: Build
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Build all files into a single binary named 'daily-stox'
RUN go build -o daily-stox .

# Stage 2: Runtime
FROM alpine:latest
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/daily-stox /app/daily-stox

# The command to run your app
ENTRYPOINT ["/app/daily-stox"]