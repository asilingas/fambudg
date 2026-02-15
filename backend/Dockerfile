# Stage 1: Build Go binaries
FROM golang:1.25-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o migrate ./cmd/migrate/main.go

# Stage 2: Runtime
FROM alpine:3.19

RUN apk add --no-cache ca-certificates wget && \
    adduser -D -g '' appuser

WORKDIR /app

COPY --from=builder /build/server .
COPY --from=builder /build/migrate .
COPY --from=builder /build/migrations/ ./migrations/

USER appuser

EXPOSE 8080

CMD ["./server"]
