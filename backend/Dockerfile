# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install gcc and required dependencies
RUN apk add --no-cache gcc musl-dev

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/main .

# Create .env file from build args
ENV HUGGINGFACE_API_KEY=""
ENV PORT="8080"

EXPOSE 8080

CMD ["./main"]