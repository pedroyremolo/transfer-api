FROM golang:alpine as builder

LABEL maintainer="Pedro Lopes <pedroyremolo@gmail.com>"

# Install git.
RUN apk update && apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

# Test before running
RUN CGO_ENABLED=0 GOOS=linux go test -v ./...

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/transfer-server

# Start a new stage from scratch
FROM alpine as runtime

RUN adduser -D worker
USER worker
WORKDIR /home/worker

# Copy the Pre-built binary file from the previous stage. Observe we also copied the .env file
COPY --from=builder --chown=worker:worker /app/main .

# Expose port to the outside world
EXPOSE $APP_PORT

#Command to run the executable
CMD ["./main"]