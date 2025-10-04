FROM golang:1.25.1

WORKDIR /app

# Install required system dependencies
RUN apt-get update && apt-get install -y gcc libc6-dev sqlite3 && rm -rf /var/lib/apt/lists/*

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the localstore service
RUN CGO_ENABLED=1 GOOS=linux go build -o /localstore ./cmd/localstore/main.go

# Expose port for database connections
EXPOSE 9000

CMD ["/localstore"]