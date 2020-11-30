# Start from golang base image
FROM golang:alpine as builder

# ENV GO111MODULE=on

# Add Maintainer info
LABEL maintainer="Fabien Bellanger <valentil@gmail.com>"

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

# Set the current working directory inside the container 
WORKDIR /app

# Copy go mod and sum files 
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and the go.sum files are not changed 
RUN go mod download 

# Copy the source from the current directory to the working Directory inside the container 
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o go-fiber .

# -----------------------------------------------------------------------------

# Start a new stage from scratch
FROM alpine:latest

RUN adduser -S -D -H -h /app appuser
USER appuser

WORKDIR /app

# Copy the Pre-built binary file from the previous stage. Observe we also copied the .env file
COPY --from=builder /app/go-fiber .
COPY --from=builder /app/config.toml .     
COPY --from=builder /app/projects.json .

EXPOSE 8888

CMD ["./go-fiber"]
