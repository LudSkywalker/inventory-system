FROM golang:1.21-alpine

WORKDIR /app

# Install required system dependencies
RUN apk add --no-cache gcc musl-dev sqlite

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