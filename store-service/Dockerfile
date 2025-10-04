FROM golang:1.25.1

WORKDIR /app

# Enable CGO
ENV CGO_ENABLED=1
ENV CGO_CFLAGS="-D_LARGEFILE64_SOURCE"

# Install SQLite and build dependencies
RUN apt-get update && apt-get install -y gcc libc6-dev sqlite3 libsqlite3-dev && rm -rf /var/lib/apt/lists/*

# Install swag
RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Run tests
RUN go test ./...

# Generate Swagger docs
RUN swag init --dir ./store-service -g ./cmd/main.go -o ./docs

# Build the application
RUN go build -o main ./store-service/cmd/main.go

# Run the application
CMD ["./main"]