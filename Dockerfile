# Build Stage
FROM golang:1.22 as builder

# Label the stage and the maintainer of the Dockerfile
LABEL maintainer="yasinmh@stud.ntnu.no"
LABEL stage=builder

# Set the working directory
WORKDIR /go/src/youvies-backend

# Copy go.mod and go.sum for dependency management
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Compile the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o youvies-backend ./cmd

# Production Stage
FROM alpine:latest

# Install dependencies required for running the application and health checks
RUN apk --no-cache add ca-certificates curl

# Set the working directory inside the production image
WORKDIR /app

# Copy the executable and .env file from the builder stage
COPY --from=builder /go/src/youvies-backend/youvies-backend .
COPY --from=builder /go/src/youvies-backend/.env .

# Expose the application port
EXPOSE 5000

# Health check configuration
HEALTHCHECK --interval=30s --timeout=30s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:5000/health || exit 1

# Define the default command to run the application
ENTRYPOINT ["./youvies-backend"]