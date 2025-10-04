FROM golang:1.25.1

WORKDIR /app

# Install SQLite and build dependencies
RUN apt-get update && apt-get install -y gcc libc6-dev sqlite3 libsqlite3-dev && rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Run tests
RUN go test ./...

# Build the application
RUN go build -o main ./operator-service/cmd/main.go

# Run the application
CMD ["./main"]