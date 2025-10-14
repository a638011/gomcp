#!/bin/bash
# Production startup script for gomcp server
# Handles environment validation, graceful startup, and monitoring

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check required environment variables
check_env() {
    local required_vars=("MCP_HOST" "MCP_PORT")
    local missing=0

    for var in "${required_vars[@]}"; do
        if [ -z "${!var:-}" ]; then
            log_error "Required environment variable ${var} is not set"
            missing=1
        fi
    done

    if [ $missing -eq 1 ]; then
        log_error "Please set all required environment variables"
        exit 1
    fi
}

# Validate configuration
validate_config() {
    log_info "Validating configuration..."
    
    # Check if port is valid
    if [ "${MCP_PORT}" -lt 1024 ] || [ "${MCP_PORT}" -gt 65535 ]; then
        log_error "Invalid port: ${MCP_PORT}. Must be between 1024-65535"
        exit 1
    fi

    # Check if PostgreSQL is configured properly
    if [ "${ENABLE_AUTH:-false}" == "true" ]; then
        if [ -n "${POSTGRES_HOST:-}" ] && [ -n "${POSTGRES_DB:-}" ]; then
            log_info "PostgreSQL storage configured"
        else
            log_warn "In-memory storage will be used (not recommended for production)"
        fi
    fi

    log_info "Configuration valid"
}

# Wait for dependencies (PostgreSQL)
wait_for_postgres() {
    if [ -z "${POSTGRES_HOST:-}" ]; then
        return 0
    fi

    log_info "Waiting for PostgreSQL at ${POSTGRES_HOST}:${POSTGRES_PORT:-5432}..."
    
    local max_attempts=30
    local attempt=0

    while [ $attempt -lt $max_attempts ]; do
        if pg_isready -h "${POSTGRES_HOST}" -p "${POSTGRES_PORT:-5432}" -U "${POSTGRES_USER:-postgres}" > /dev/null 2>&1; then
            log_info "PostgreSQL is ready"
            return 0
        fi
        
        attempt=$((attempt + 1))
        log_info "Attempt ${attempt}/${max_attempts}: PostgreSQL not ready yet..."
        sleep 2
    done

    log_error "PostgreSQL failed to become ready after ${max_attempts} attempts"
    exit 1
}

# Start the server
start_server() {
    log_info "Starting gomcp server..."
    log_info "Host: ${MCP_HOST}:${MCP_PORT}"
    log_info "Transport: ${MCP_TRANSPORT_PROTOCOL:-http}"
    log_info "Log Level: ${LOG_LEVEL:-INFO}"
    
    # Execute the binary
    exec ./bin/gomcp
}

# Main execution
main() {
    log_info "gomcp startup script"
    
    # Load .env file if present
    if [ -f .env ]; then
        log_info "Loading .env file..."
        set -a
        source .env
        set +a
    fi

    # Check and validate
    check_env
    validate_config
    
    # Wait for dependencies if needed
    if [ "${ENABLE_AUTH:-false}" == "true" ] && [ -n "${POSTGRES_HOST:-}" ]; then
        wait_for_postgres
    fi

    # Start the server
    start_server
}

# Handle signals for graceful shutdown
trap 'log_info "Received shutdown signal, exiting..."; exit 0' SIGTERM SIGINT

# Run main
main "$@"

