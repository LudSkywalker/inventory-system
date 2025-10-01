# Distributed Inventory Management System

## System Architecture

This system implements a distributed inventory management solution with the following components:

### 1. Store Service
- Local inventory management for each store
- In-memory SQLite database for fast operations
- Fiber HTTP server for API endpoints
- Kafka producer for event publishing
- Hexagonal architecture implementation

### 2. Operator Service
- Processes inventory events from all stores
- Maintains global inventory state
- Kafka consumer for event processing
- RESTful API for global inventory queries
- In-memory SQLite for state management

### 3. Backend Service
- Public API for inventory queries
- Aggregates data from operator service
- Provides filtered views of inventory data

## Technical Stack

- **Language:** Go 1.21+
- **Web Framework:** Fiber
- **Message Broker:** Apache Kafka
- **Database:** SQLite (in-memory)
- **Architecture:** Hexagonal (Ports and Adapters)

## Project Structure

```
inventory-system/
├── store-service/
│   ├── domain/         # Domain entities and business logic
│   ├── app/            # Application services and DTOs
│   └── infra/         # Infrastructure adapters
├── operator-service/
│   ├── domain/
│   ├── app/
│   └── infra/
└── backend-service/
    ├── app/
    └── infra/
```

## Event Flow

1. Store updates inventory locally
2. Event published to Kafka topic
3. Operator service consumes event
4. Global state updated
5. Backend service queries updated state

## API Endpoints

### Store Service
- POST /api/v1/inventory - Update inventory
- GET /api/v1/inventory/:storeId/:itemId - Get item inventory
- DELETE /api/v1/inventory/:storeId/:itemId - Delete inventory

### Operator Service
- GET /api/v1/global-inventory - Get all inventories
- GET /api/v1/global-inventory/:storeId/:itemId - Get specific inventory

### Backend Service
- GET /api/v1/inventory - Get global inventory
- GET /api/v1/inventory/store/:storeId - Get store inventory
- GET /api/v1/inventory/item/:itemId - Get item inventory

## Setup and Deployment

1. Start Kafka:
```bash
docker-compose up kafka
```

2. Start all services:
```bash
docker-compose up
```

## Configuration

Environment variables:
- `KAFKA_BROKERS` - Comma-separated list of Kafka brokers
- `PORT` - HTTP server port

## Testing

Run tests:
```bash
go test ./...
```

## Monitoring

- Prometheus metrics available at `/metrics`
- Logging using structured logging
- Kafka lag monitoring for consumer health

## Security

- HTTPS for all services
- Rate limiting implemented
- Error handling and validation