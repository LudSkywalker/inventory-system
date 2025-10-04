#!/bin/bash

# E2E test for db-service

DB_URL="http://localhost:8083"

echo "Testing db-service health"
response=$(curl -s -o /dev/null -w "%{http_code}" "$DB_URL/health")
if [ "$response" -ne 200 ]; then
    echo "❌ DB health failed: $response"
    exit 1
fi
echo "✅ DB health OK"

echo "Testing GET /inventories"
response=$(curl -s -o /dev/null -w "%{http_code}" "$DB_URL/inventories")
if [ "$response" -ne 200 ]; then
    echo "❌ DB GET inventories failed: $response"
    exit 1
fi
echo "✅ DB GET inventories OK"