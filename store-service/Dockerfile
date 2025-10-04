FROM golang:1.25.1-alpine

WORKDIR /app

# Enable CGO
ENV CGO_ENABLED=1
ENV CGO_CFLAGS="-D_LARGEFILE64_SOURCE"

# Install SQLite and build dependencies
RUN apk add --no-cache gcc musl-dev sqlite-dev libc6-compat

# Install swag
RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Generate Swagger docs
RUN swag init --dir ./store-service -g ./cmd/main.go -o ./docs

# Build the application
RUN go build -o main ./store-service/cmd/main.go

# Run the application
CMD ["./main"]