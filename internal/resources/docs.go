package resources

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// GetQuickStartGuide returns the quick start documentation
func GetQuickStartGuide(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	guide := `# gomcp Quick Start Guide

## Installation

1. Clone the repository
2. Build: make build
3. Run: ./bin/gomcp

## Basic Usage

### Test the server
curl http://localhost:8081/health

### Call a tool
curl -X POST http://localhost:8081/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": "1",
    "method": "tools/call",
    "params": {
      "name": "multiply_numbers",
      "arguments": {"a": 13, "b": 34}
    }
  }'

## Adding Custom Tools

Create a file in internal/tools/:

` + "```go" + `
package tools

import (
    "context"
    "github.com/modelcontextprotocol/go-sdk/mcp"
)

type MyToolInput struct {
    Message string ` + "`json:\"message\" jsonschema:\"required,input message\"`" + `
}

type MyToolOutput struct {
    Result string ` + "`json:\"result\"`" + `
}

func MyToolSDK(
    ctx context.Context,
    req *mcp.CallToolRequest,
    input MyToolInput,
) (*mcp.CallToolResult, MyToolOutput, error) {
    return nil, MyToolOutput{Result: "Processed: " + input.Message}, nil
}
` + "```" + `

Register in internal/mcp/server_sdk.go:

` + "```go" + `
mcp.AddTool(server, &mcp.Tool{Name: "my_tool"}, tools.MyToolSDK)
` + "```" + `

Rebuild and run!

## Configuration

Set environment variables:
- MCP_HOST (default: localhost)
- MCP_PORT (default: 8081)
- ENABLE_AUTH (default: false)
- LOG_LEVEL (default: INFO)

## Resources

- Full README: README.md
- Tools documentation: internal/tools/README.md
- Official MCP docs: https://modelcontextprotocol.io`

	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      "docs://quickstart",
				MIMEType: "text/markdown",
				Text:     guide,
			},
		},
	}, nil
}

// GetAPIReference returns API endpoint documentation
func GetAPIReference(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	reference := `# gomcp API Reference

## Endpoints

### Health Check
**GET /health**

Returns server health status.

Response:
` + "```json" + `
{
  "status": "healthy",
  "service": "template-mcp-server",
  "transport_protocol": "http",
  "version": "0.1.0"
}
` + "```" + `

### Metrics
**GET /metrics**

Returns comprehensive server metrics including MCP endpoint usage.

Response includes:
- Server version and uptime
- Memory usage
- MCP endpoint statistics (legacy vs SDK)
- Tool call counts per endpoint
- Client identification
- Response times

### MCP Protocol (Legacy)
**POST /mcp**

Legacy JSON-RPC 2.0 endpoint for MCP protocol.

Methods:
- tools/list - List available tools
- tools/call - Execute a tool
- prompts/list - List available prompts
- prompts/get - Get a prompt
- resources/list - List available resources
- resources/read - Read a resource

### MCP Protocol (SDK)
**GET/POST /mcp/sse**

Official SDK endpoint using Server-Sent Events transport.

Supports full MCP 2024-11-05 specification.

## OAuth2 Endpoints

(Available when ENABLE_AUTH=true)

- **POST /auth/register** - Register OAuth client
- **GET /auth/authorize** - Authorization request
- **POST /auth/token** - Token exchange
- **POST /auth/revoke** - Revoke token
- **GET /.well-known/oauth-protected-resource** - Resource metadata
- **GET /.well-known/oauth-authorization-server** - Server metadata

## Example: Call a Tool

` + "```bash" + `
curl -X POST http://localhost:8081/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": "request-1",
    "method": "tools/call",
    "params": {
      "name": "multiply_numbers",
      "arguments": {"a": 42, "b": 10}
    }
  }'
` + "```" + `

Response:
` + "```json" + `
{
  "jsonrpc": "2.0",
  "id": "request-1",
  "result": {
    "content": [{
      "type": "text",
      "text": "{\"result\": 420, \"status\": \"success\"}"
    }]
  }
}
` + "```" + ``

	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      "docs://api",
				MIMEType: "text/markdown",
				Text:     reference,
			},
		},
	}, nil
}
