#!/bin/bash
# Graceful shutdown script for gomcp server

set -euo pipefail

PID_FILE="${1:-/tmp/gomcp.pid}"
TIMEOUT="${2:-30}"

if [ ! -f "${PID_FILE}" ]; then
    echo "PID file not found: ${PID_FILE}"
    exit 1
fi

PID=$(cat "${PID_FILE}")

if ! kill -0 "${PID}" 2>/dev/null; then
    echo "Process ${PID} is not running"
    rm -f "${PID_FILE}"
    exit 0
fi

echo "Sending SIGTERM to process ${PID}..."
kill -TERM "${PID}"

# Wait for graceful shutdown
for i in $(seq 1 "${TIMEOUT}"); do
    if ! kill -0 "${PID}" 2>/dev/null; then
        echo "Process ${PID} terminated gracefully"
        rm -f "${PID_FILE}"
        exit 0
    fi
    sleep 1
done

# Force kill if still running
echo "Process ${PID} did not terminate gracefully, sending SIGKILL..."
kill -KILL "${PID}"
rm -f "${PID_FILE}"
echo "Process ${PID} killed"

