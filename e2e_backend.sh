#!/bin/bash

# E2E test for backend-service

BACKEND_URL="http://localhost:8081"

echo "Testing backend-service health"
response=$(curl -s -o /dev/null -w "%{http_code}" "$BACKEND_URL/health")
if [ "$response" -ne 200 ]; then
    echo "❌ Backend health failed: $response"
    exit 1
fi
echo "✅ Backend health OK"

echo "Testing GET /api/v1/inventory"
response=$(curl -s -o /dev/null -w "%{http_code}" "$BACKEND_URL/api/v1/inventory")
if [ "$response" -ne 200 ]; then
    echo "❌ Backend GET inventory failed: $response"
    exit 1
fi
echo "✅ Backend GET inventory OK"

echo "Testing GET /api/v1/inventory/grouped"
response=$(curl -s -o /dev/null -w "%{http_code}" "$BACKEND_URL/api/v1/inventory/grouped")
if [ "$response" -ne 200 ]; then
    echo "❌ Backend GET grouped inventory failed: $response"
    exit 1
fi
echo "✅ Backend GET grouped inventory OK"