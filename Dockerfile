# Use the official Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
# https://hub.docker.com/_/golang
FROM golang:1.20 as builder

# Copy the local package files to the container's workspace.
WORKDIR /app
COPY . .


# Build the Go app for a Linux environment
RUN CGO_ENABLED=0 GOOS=linux go build -v -o tiny-url

# Use the official lightweight Alpine image.
# https://hub.docker.com/_/alpine
# This is for the smaller image size.
FROM alpine:latest

EXPOSE 8181
# Copy the binary from builder stage.
COPY --from=builder /app/tiny-url /tiny-url

# Run the application binary.
CMD ["/tiny-url"]