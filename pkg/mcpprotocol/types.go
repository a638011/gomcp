package mcpprotocol

import "encoding/json"

// JSONRPCVersion is the JSON-RPC version used by MCP
const JSONRPCVersion = "2.0"

// Common MCP methods
const (
	MethodToolsList   = "tools/list"
	MethodToolsCall   = "tools/call"
	MethodInitialize  = "initialize"
	MethodInitialized = "initialized"
	MethodPing        = "ping"
)

// JSONRPCRequest represents a JSON-RPC 2.0 request
type JSONRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// JSONRPCResponse represents a JSON-RPC 2.0 response
type JSONRPCResponse struct {
	JSONRPC string        `json:"jsonrpc"`
	ID      interface{}   `json:"id,omitempty"`
	Result  interface{}   `json:"result,omitempty"`
	Error   *JSONRPCError `json:"error,omitempty"`
}

// JSONRPCError represents a JSON-RPC 2.0 error
type JSONRPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Error codes
const (
	ParseError     = -32700
	InvalidRequest = -32600
	MethodNotFound = -32601
	InvalidParams  = -32602
	InternalError  = -32603
)

// ToolDefinition represents an MCP tool definition
type ToolDefinition struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

// ToolsListResult represents the result of tools/list
type ToolsListResult struct {
	Tools []ToolDefinition `json:"tools"`
}

// ToolsCallParams represents parameters for tools/call
type ToolsCallParams struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
}

// ToolsCallResult represents the result of tools/call
type ToolsCallResult struct {
	Content []ContentItem `json:"content"`
	IsError bool          `json:"isError,omitempty"`
}

// ContentItem represents a piece of content in a tool result
type ContentItem struct {
	Type     string `json:"type"`
	Text     string `json:"text,omitempty"`
	Data     string `json:"data,omitempty"`
	MimeType string `json:"mimeType,omitempty"`
}

// InitializeParams represents parameters for initialize method
type InitializeParams struct {
	ProtocolVersion string                 `json:"protocolVersion"`
	ClientInfo      ClientInfo             `json:"clientInfo"`
	Capabilities    map[string]interface{} `json:"capabilities,omitempty"`
}

// ClientInfo represents information about the MCP client
type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// InitializeResult represents the result of initialize
type InitializeResult struct {
	ProtocolVersion string                 `json:"protocolVersion"`
	ServerInfo      ServerInfo             `json:"serverInfo"`
	Capabilities    map[string]interface{} `json:"capabilities"`
}

// ServerInfo represents information about the MCP server
type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// NewSuccessResponse creates a successful JSON-RPC response
func NewSuccessResponse(id interface{}, result interface{}) *JSONRPCResponse {
	return &JSONRPCResponse{
		JSONRPC: JSONRPCVersion,
		ID:      id,
		Result:  result,
	}
}

// NewErrorResponse creates an error JSON-RPC response
func NewErrorResponse(id interface{}, code int, message string, data interface{}) *JSONRPCResponse {
	return &JSONRPCResponse{
		JSONRPC: JSONRPCVersion,
		ID:      id,
		Error: &JSONRPCError{
			Code:    code,
			Message: message,
			Data:    data,
		},
	}
}

// NewTextContent creates a text content item
func NewTextContent(text string) ContentItem {
	return ContentItem{
		Type: "text",
		Text: text,
	}
}

// NewImageContent creates an image content item
func NewImageContent(data string, mimeType string) ContentItem {
	return ContentItem{
		Type:     "image",
		Data:     data,
		MimeType: mimeType,
	}
}
