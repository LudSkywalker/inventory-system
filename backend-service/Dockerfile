FROM golang:1.25.1-alpine

WORKDIR /app

# Install swag
RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Ensure docs directory exists
RUN mkdir -p backend-service/docs

# Generate Swagger docs
RUN swag init --dir ./backend-service -g ./cmd/main.go -o ./backend-service/docs

# Copy generated swagger.json to docs.json
RUN cp backend-service/docs/swagger.json backend-service/docs/docs.json

# Build the application
RUN go build -o main ./backend-service/cmd/main.go

# Run the application
CMD ["./main"]