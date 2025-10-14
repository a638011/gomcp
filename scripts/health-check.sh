#!/bin/bash
# Health check script for gomcp server
# Usage: ./scripts/health-check.sh [host] [port]

HOST="${1:-localhost}"
PORT="${2:-8080}"
URL="http://${HOST}:${PORT}/health"

echo "Checking health at ${URL}..."

RESPONSE=$(curl -s -w "\n%{http_code}" "${URL}")
HTTP_CODE=$(echo "${RESPONSE}" | tail -n1)
BODY=$(echo "${RESPONSE}" | head -n-1)

if [ "${HTTP_CODE}" -eq 200 ]; then
    echo "✓ Server is healthy (HTTP ${HTTP_CODE})"
    echo "${BODY}" | jq '.' 2>/dev/null || echo "${BODY}"
    exit 0
else
    echo "✗ Server is unhealthy (HTTP ${HTTP_CODE})"
    echo "${BODY}"
    exit 1
fi

