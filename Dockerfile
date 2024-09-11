# Use the official Go image as a build stage
FROM golang:1.23 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app for Windows
RUN go build -o logservice main.go

# Expose port 8080 to the outside world
EXPOSE 8080

# Accept build-time arguments (optional)
ARG DB_USER
ARG DB_PASSWORD
ARG DB_NAME
ARG HTTP_HOST

# Set environment variables for runtime
ENV DB_USER=${DB_USER}
ENV DB_PASSWORD=${DB_PASSWORD}
ENV DB_NAME=${DB_NAME}
ENV HTTP_HOST=${HTTP_HOST}
# Run the binary program produced by `go build`
CMD ["./logservice"]
