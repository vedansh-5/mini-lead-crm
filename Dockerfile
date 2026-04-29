# Start from the official Go image
FROM golang:alpine

# Set the working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go application from our main file
RUN go build -o main ./cmd/api

# Expose port 8080
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
