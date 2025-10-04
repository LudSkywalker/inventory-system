#!/bin/bash

# E2E test for store-service

STORE_URL="http://localhost:8080"

echo "Testing store-service health"
response=$(curl -s -o /dev/null -w "%{http_code}" "$STORE_URL/health")
if [ "$response" -ne 200 ]; then
    echo "❌ Store health failed: $response"
    exit 1
fi
echo "✅ Store health OK"

echo "Testing POST /inventory"
data='{"item_id":"test-item","item_name":"Test Item","store_id":"test-store","quantity":10}'
response=$(curl -s -X POST -H "Content-Type: application/json" -d "$data" -o /dev/null -w "%{http_code}" "$STORE_URL/inventory")
if [ "$response" -ne 200 ]; then
    echo "❌ Store POST inventory failed: $response"
    exit 1
fi
echo "✅ Store POST inventory OK"

echo "Testing GET /inventory"
response=$(curl -s -o /dev/null -w "%{http_code}" "$STORE_URL/inventory")
if [ "$response" -ne 200 ]; then
    echo "❌ Store GET inventory failed: $response"
    exit 1
fi
echo "✅ Store GET inventory OK"