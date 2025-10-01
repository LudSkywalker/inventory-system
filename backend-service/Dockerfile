FROM golang:1.21-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the application
RUN go build -o main ./backend-service/cmd/main.go

# Run the application
CMD ["./main"]