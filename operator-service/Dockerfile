FROM golang:1.25.1-alpine

WORKDIR /app

# Install SQLite and build dependencies
RUN apk add --no-cache gcc musl-dev sqlite-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the application
RUN go build -o main ./operator-service/cmd/main.go

# Run the application
CMD ["./main"]