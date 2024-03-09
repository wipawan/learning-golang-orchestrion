# Use the official Golang base image
FROM golang:1.21.5

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the working directory
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire application
COPY . .

# Build the Go application
RUN go build -o main .

# Expose port 8080 for the application
EXPOSE 8080

ARG DD_GIT_REPOSITORY_URL
ARG DD_GIT_COMMIT_SHA
ENV DD_GIT_REPOSITORY_URL=github.com/jon94/learn-golang
ENV DD_GIT_COMMIT_SHA=${DD_GIT_COMMIT_SHA}

# Command to run the executable
CMD ["./main"]