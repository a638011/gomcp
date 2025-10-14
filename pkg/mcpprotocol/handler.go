package mcpprotocol

import (
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog"
)

// ToolHandler is a function that handles tool execution
type ToolHandler func(args map[string]interface{}) (interface{}, error)

// Server represents an MCP protocol server
type Server struct {
	name         string
	version      string
	tools        map[string]*Tool
	capabilities map[string]interface{}
	logger       *zerolog.Logger
}

// Tool represents a registered MCP tool
type Tool struct {
	Definition ToolDefinition
	Handler    ToolHandler
}

// NewServer creates a new MCP server
func NewServer(name, version string, logger *zerolog.Logger) *Server {
	return &Server{
		name:    name,
		version: version,
		tools:   make(map[string]*Tool),
		capabilities: map[string]interface{}{
			"tools": map[string]interface{}{
				"completion": map[string]interface{}{}, // Structured outputs
			},
			"prompts":     map[string]interface{}{},
			"resources":   map[string]interface{}{},
			"roots":       map[string]interface{}{},                // Protocol 2025-06-18
			"sampling":    map[string]interface{}{},                // Protocol 2025-06-18
			"elicitation": map[string]interface{}{},                // Protocol 2025-06-18
			"logging":     map[string]interface{}{"level": "info"}, // Protocol 2025-06-18
			"pagination": map[string]interface{}{
				"support":     true,
				"maxPageSize": 100,
			}, // Protocol 2025-06-18
		},
		logger: logger,
	}
}

// RegisterTool registers a new tool with the server
func (s *Server) RegisterTool(def ToolDefinition, handler ToolHandler) {
	s.tools[def.Name] = &Tool{
		Definition: def,
		Handler:    handler,
	}
}

// HandleRequest processes an MCP JSON-RPC request
func (s *Server) HandleRequest(reqData []byte) ([]byte, error) {
	var req JSONRPCRequest
	if err := json.Unmarshal(reqData, &req); err != nil {
		errResp := NewErrorResponse(nil, ParseError, "Parse error", err.Error())
		return json.Marshal(errResp)
	}

	var resp *JSONRPCResponse

	switch req.Method {
	case MethodInitialize:
		resp = s.handleInitialize(req)
	case MethodToolsList:
		resp = s.handleToolsList(req)
	case MethodToolsCall:
		resp = s.handleToolsCall(req)
	case MethodPing:
		resp = s.handlePing(req)
	default:
		resp = NewErrorResponse(req.ID, MethodNotFound, "Method not found", req.Method)
	}

	return json.Marshal(resp)
}

func (s *Server) handleInitialize(req JSONRPCRequest) *JSONRPCResponse {
	var params InitializeParams
	if len(req.Params) > 0 {
		if err := json.Unmarshal(req.Params, &params); err != nil {
			return NewErrorResponse(req.ID, InvalidParams, "Invalid params", err.Error())
		}
	}

	result := InitializeResult{
		ProtocolVersion: "2024-11-05",
		ServerInfo: ServerInfo{
			Name:    s.name,
			Version: s.version,
		},
		Capabilities: s.capabilities,
	}

	return NewSuccessResponse(req.ID, result)
}

func (s *Server) handlePing(req JSONRPCRequest) *JSONRPCResponse {
	// Ping method returns an empty object as per MCP spec
	result := map[string]interface{}{}

	if s.logger != nil {
		s.logger.Debug().
			Interface("request_id", req.ID).
			Msg("Ping request received")
	}

	return NewSuccessResponse(req.ID, result)
}

func (s *Server) handleToolsList(req JSONRPCRequest) *JSONRPCResponse {
	tools := make([]ToolDefinition, 0, len(s.tools))
	for _, tool := range s.tools {
		tools = append(tools, tool.Definition)
	}

	result := ToolsListResult{
		Tools: tools,
	}

	return NewSuccessResponse(req.ID, result)
}

func (s *Server) handleToolsCall(req JSONRPCRequest) *JSONRPCResponse {
	var params ToolsCallParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return NewErrorResponse(req.ID, InvalidParams, "Invalid params", err.Error())
	}

	tool, ok := s.tools[params.Name]
	if !ok {
		return NewErrorResponse(req.ID, MethodNotFound, fmt.Sprintf("Tool not found: %s", params.Name), nil)
	}

	// Log tool call with request ID
	if s.logger != nil {
		s.logger.Info().
			Interface("tool_call_id", req.ID).
			Str("tool_name", params.Name).
			Interface("arguments", params.Arguments).
			Msg("Tool call invoked")
	}

	// Execute the tool handler
	result, err := tool.Handler(params.Arguments)
	if err != nil {
		return NewErrorResponse(req.ID, InternalError, err.Error(), nil)
	}

	// Convert result to ToolsCallResult format
	toolResult := formatToolResult(result, false)
	return NewSuccessResponse(req.ID, toolResult)
}

// formatToolResult converts a tool result into MCP ToolsCallResult format
func formatToolResult(result interface{}, isError bool) ToolsCallResult {
	// Convert result to JSON string for text content
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return ToolsCallResult{
			Content: []ContentItem{
				NewTextContent(fmt.Sprintf("Error formatting result: %v", err)),
			},
			IsError: true,
		}
	}

	return ToolsCallResult{
		Content: []ContentItem{
			NewTextContent(string(jsonBytes)),
		},
		IsError: isError,
	}
}

// GetTools returns all registered tools
func (s *Server) GetTools() []ToolDefinition {
	tools := make([]ToolDefinition, 0, len(s.tools))
	for _, tool := range s.tools {
		tools = append(tools, tool.Definition)
	}
	return tools
}
