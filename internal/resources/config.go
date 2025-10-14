package resources

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// GetConfigTemplate returns a configuration template
func GetConfigTemplate(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	configTemplate := map[string]interface{}{
		"server": map[string]interface{}{
			"host":              "localhost",
			"port":              8081,
			"transport":         "http",
			"enable_auth":       false,
			"log_level":         "INFO",
			"read_timeout_sec":  30,
			"write_timeout_sec": 30,
			"idle_timeout_sec":  120,
		},
		"database": map[string]interface{}{
			"enabled":         false,
			"host":            "localhost",
			"port":            5432,
			"database":        "mcpdb",
			"user":            "mcpuser",
			"password":        "changeme",
			"pool_size":       10,
			"max_connections": 20,
			"ssl_mode":        "disable",
		},
		"oauth": map[string]interface{}{
			"enabled":                false,
			"authorization_code_ttl": 600,
			"access_token_ttl":       3600,
			"refresh_token_ttl":      2592000,
		},
		"cors": map[string]interface{}{
			"enabled":     true,
			"origins":     []string{"*"},
			"credentials": true,
			"methods":     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			"headers":     []string{"Accept", "Authorization", "Content-Type"},
		},
		"session": map[string]interface{}{
			"secret":      "change-this-secret-key",
			"cookie_name": "gomcp_session",
			"max_age_sec": 86400,
			"http_only":   true,
			"secure":      false,
			"same_site":   "lax",
		},
		"notes": map[string]string{
			"env_vars":      "All config can be set via environment variables (uppercase, underscore-separated)",
			"example":       "MCP_HOST, MCP_PORT, ENABLE_AUTH, POSTGRES_HOST, etc.",
			"precedence":    "Environment variables > .env file > defaults",
			"documentation": "See README.md for complete configuration guide",
		},
	}

	jsonData, err := json.MarshalIndent(configTemplate, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal config template: %w", err)
	}

	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      "config://template",
				MIMEType: "application/json",
				Text:     string(jsonData),
			},
		},
	}, nil
}

// GetEnvVarsReference returns environment variables reference
func GetEnvVarsReference(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	envVars := map[string]interface{}{
		"server_configuration": map[string]string{
			"MCP_HOST":               "Server bind address (default: localhost)",
			"MCP_PORT":               "Server port (default: 8081)",
			"MCP_TRANSPORT_PROTOCOL": "Transport protocol: http, sse (default: http)",
			"MCP_HOST_ENDPOINT":      "External endpoint URL for OAuth redirects",
			"MCP_SSL_CERTFILE":       "Path to SSL certificate file",
			"MCP_SSL_KEYFILE":        "Path to SSL private key file",
		},
		"database_configuration": map[string]string{
			"POSTGRES_HOST":            "PostgreSQL host (default: localhost)",
			"POSTGRES_PORT":            "PostgreSQL port (default: 5432)",
			"POSTGRES_DB":              "Database name",
			"POSTGRES_USER":            "Database user",
			"POSTGRES_PASSWORD":        "Database password",
			"POSTGRES_SSLMODE":         "SSL mode: disable, require, verify-ca, verify-full",
			"POSTGRES_POOL_SIZE":       "Connection pool size (default: 10)",
			"POSTGRES_MAX_CONNECTIONS": "Max connections (default: 20)",
		},
		"authentication": map[string]string{
			"ENABLE_AUTH":                "Enable OAuth2 authentication (default: false)",
			"SESSION_SECRET":             "Session encryption secret",
			"SESSION_COOKIE_NAME":        "Session cookie name (default: gomcp_session)",
			"SESSION_MAX_AGE":            "Session max age in seconds (default: 86400)",
			"AUTHORIZATION_CODE_TIMEOUT": "Auth code TTL in seconds (default: 600)",
			"ACCESS_TOKEN_TIMEOUT":       "Access token TTL in seconds (default: 3600)",
			"REFRESH_TOKEN_TIMEOUT":      "Refresh token TTL in seconds (default: 2592000)",
		},
		"cors_configuration": map[string]string{
			"CORS_ENABLED":     "Enable CORS (default: true)",
			"CORS_ORIGINS":     "Allowed origins (comma-separated)",
			"CORS_CREDENTIALS": "Allow credentials (default: true)",
			"CORS_METHODS":     "Allowed methods (comma-separated)",
			"CORS_HEADERS":     "Allowed headers (comma-separated)",
		},
		"logging": map[string]string{
			"LOG_LEVEL": "Log level: DEBUG, INFO, WARNING, ERROR (default: INFO)",
		},
		"examples": []string{
			"export MCP_PORT=9000",
			"export ENABLE_AUTH=true",
			"export POSTGRES_HOST=db.example.com",
			"export LOG_LEVEL=DEBUG",
		},
	}

	jsonData, err := json.MarshalIndent(envVars, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal env vars: %w", err)
	}

	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      "config://env-vars",
				MIMEType: "application/json",
				Text:     string(jsonData),
			},
		},
	}, nil
}
