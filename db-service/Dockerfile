FROM golang:1.25-alpine

WORKDIR /app

# Install SQLite and build dependencies
RUN apk add --no-cache gcc musl-dev sqlite-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the application
RUN CGO_CFLAGS="-D_LARGEFILE64_SOURCE" go build -o main ./db-service/cmd/main.go

# Run the application
CMD ["./main"]