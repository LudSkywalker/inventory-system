#!/bin/bash

# E2E test for operator-service

OPERATOR_URL="http://localhost:8082"

echo "Testing operator-service health"
response=$(curl -s -o /dev/null -w "%{http_code}" "$OPERATOR_URL/health")
if [ "$response" -ne 200 ]; then
    echo "❌ Operator health failed: $response"
    exit 1
fi
echo "✅ Operator health OK"