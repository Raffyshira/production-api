# --- BUILD STAGE ---
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o api ./cmd/api


# --- RUN STAGE ---
FROM alpine:latest
WORKDIR /app
RUN apk --no-cache add ca-certificates tzdata
# Copy CA certificates
COPY --from=builder /app/api .
EXPOSE 8080
CMD ["./api"]
