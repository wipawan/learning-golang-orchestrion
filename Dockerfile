# Use the official Golang base image
FROM golang:1.18-alpine

# Set necessary environment variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the working directory
COPY go.mod go.sum ./

# Copy the entire application
COPY . .

#orchestrion instrumentation
RUN go install github.com/datadog/orchestrion@v0.6.0

# Activate Orchestrion
RUN $(go env GOPATH)/bin/orchestrion -w ./

# With the newly instrumented code, manage dependency
RUN go mod tidy

# Build the Go application
RUN go build -o main .

# Expose port 8080 for the application
EXPOSE 8080

# Command to run the executable
CMD ["./main"]