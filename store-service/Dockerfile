FROM golang:1.25.1-alpine

WORKDIR /app

# Install SQLite and build dependencies
RUN apk add --no-cache gcc musl-dev sqlite-dev

# Install swag
RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Ensure docs directory exists and copy pre-generated docs
RUN mkdir -p store-service/docs

# Build the application
RUN go build -o main ./store-service/cmd/main.go

# Run the application
CMD ["./main"]