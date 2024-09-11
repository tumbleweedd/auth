FROM golang:1.22.2-alpine3.19 as builder

# Install Git
RUN apk update && apk add --no-cache git

WORKDIR /go/src/app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/app ./cmd/authService/

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/migrator ./cmd/migrator/

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage to the final image
COPY --from=builder /go/bin/app /app/app
COPY --from=builder /go/bin/migrator /app/migrator

COPY config/config.yaml /app/config.yaml
COPY migrations /app/migrations

# Set the environment variable for the configuration path
ENV CONFIG_PATH="/app/config.yaml"

# Command to run the application
CMD ["/app/app"]
