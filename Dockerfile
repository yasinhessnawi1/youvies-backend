FROM golang:1.21 as builder

# Label the stage and the maintainer of the Dockerfile.
LABEL maintainer="yasinmh@stud.ntnu.no"
LABEL stage=builder

# Set the working directory inside the container where all commands will be run.
WORKDIR /go/src/youvies-backend

# Copy the go module files first to leverage Docker cache to save re-downloading the same dependencies.
COPY go.mod go.sum /go/src/youvies-backend/
# Download all the dependencies specified in go.mod and go.sum.
RUN go mod download

# Copy the rest of the application code to the container.
COPY api/ /go/src/youvies-backend/api
COPY database/ /go/src/youvies-backend/database
COPY models/ /go/src/youvies-backend/models
COPY cmd/ /go/src/youvies-backend/cmd
COPY utils/ /go/src/youvies-backend/utils
COPY scraper/ /go/src/youvies-backend/scraper
COPY .env /go/src/youvies-backend/
# Compile the application to an executable named 'dashboard'.
# Specify the directory of the main package if it's not in the root directory.
# CGO_ENABLED=0 is required for building a statically linked binary.
# GOOS=linux specifies that the binary is for Linux OS.
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o youvies-backend ./cmd

# Start a new stage from scratch to keep the final image clean and small.
FROM alpine:latest

# Copy only the built executable from the builder stage into this lightweight image.
COPY --from=builder /go/src/youvies-backend/youvies-backend .
COPY --from=builder /go/src/youvies-backend/.env .
COPY --from=builder /go/src/youvies-backend/utils/ ./utils
# Inform Docker that the container listens on port 5000 at runtime.
EXPOSE 5000

# Define a health check for the application.
# This will help Docker know how to test that the application is working.
HEALTHCHECK --interval=30s --timeout=30s --start-period=5s --retries=3 \
  CMD [ "curl", "-f", "http://localhost:5000/health" ]

# Set the container's default executable which is the application binary.
ENTRYPOINT ["./youvies-backend"]
