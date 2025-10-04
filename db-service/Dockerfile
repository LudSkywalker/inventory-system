FROM golang:1.25.1

WORKDIR /app

# Install SQLite and build dependencies
RUN apt-get update && apt-get install -y gcc libc6-dev sqlite3 libsqlite3-dev && rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the application
RUN CGO_CFLAGS="-D_LARGEFILE64_SOURCE" go build -o main ./db-service/cmd/main.go

# Run the application
CMD ["./main"]