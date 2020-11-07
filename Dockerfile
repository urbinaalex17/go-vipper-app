# Use the golang version 1.15.2 as base image
FROM golang:1.15.2 AS build
# Set the working directory
WORKDIR /app
# Copy specific files and folders from source to the dest path in the image's filesystem.
COPY go.mod go.mod
COPY src src
# Compile source code
RUN CGO_ENABLED=0 GOOS=linux go build src/main.go


# Use the current stable alpine image for production stage
FROM alpine:latest AS production
# Set the working directory
WORKDIR /app
# Copy cli binary from build stage
COPY --from=build /app/main .
# Create the app user
RUN adduser -S -D -H -h /app appuser
RUN chown -R appuser /app
USER appuser
# Define entrypoint and command to execute
CMD ["./main"]
