package completion

import (
	"encoding/json"
	"fmt"
)

// StructuredOutput represents a tool response with both text content and structured data
type StructuredOutput struct {
	Content           []ContentItem          `json:"content"`
	StructuredContent map[string]interface{} `json:"structuredContent,omitempty"`
	IsError           bool                   `json:"isError"`
}

// ContentItem represents a piece of content in MCP format
type ContentItem struct {
	Type string `json:"type"` // "text", "image", "resource"
	Text string `json:"text,omitempty"`
	Data string `json:"data,omitempty"`
	URI  string `json:"uri,omitempty"`
}

// ToolOutputSchema defines the schema for a tool's output
type ToolOutputSchema struct {
	Type        string                            `json:"type"`
	Properties  map[string]map[string]interface{} `json:"properties,omitempty"`
	Required    []string                          `json:"required,omitempty"`
	Description string                            `json:"description,omitempty"`
}

// NewStructuredOutput creates a structured output with both text and data
func NewStructuredOutput(text string, data map[string]interface{}) *StructuredOutput {
	return &StructuredOutput{
		Content: []ContentItem{
			{
				Type: "text",
				Text: text,
			},
		},
		StructuredContent: data,
		IsError:           false,
	}
}

// NewTextOutput creates a simple text-only output
func NewTextOutput(text string) *StructuredOutput {
	return &StructuredOutput{
		Content: []ContentItem{
			{
				Type: "text",
				Text: text,
			},
		},
		IsError: false,
	}
}

// NewErrorOutput creates an error output
func NewErrorOutput(errorMsg string) *StructuredOutput {
	return &StructuredOutput{
		Content: []ContentItem{
			{
				Type: "text",
				Text: fmt.Sprintf("Error: %s", errorMsg),
			},
		},
		IsError: true,
	}
}

// AddTextContent adds text content to the output
func (s *StructuredOutput) AddTextContent(text string) {
	s.Content = append(s.Content, ContentItem{
		Type: "text",
		Text: text,
	})
}

// AddImageContent adds image content to the output
func (s *StructuredOutput) AddImageContent(data, mimeType string) {
	s.Content = append(s.Content, ContentItem{
		Type: "image",
		Data: data,
	})
}

// AddResourceContent adds resource reference to the output
func (s *StructuredOutput) AddResourceContent(uri string) {
	s.Content = append(s.Content, ContentItem{
		Type: "resource",
		URI:  uri,
	})
}

// SetStructuredData sets the structured data portion
func (s *StructuredOutput) SetStructuredData(data map[string]interface{}) {
	s.StructuredContent = data
}

// ToJSON converts the output to JSON
func (s *StructuredOutput) ToJSON() ([]byte, error) {
	return json.Marshal(s)
}

// ToJSONString converts the output to a JSON string
func (s *StructuredOutput) ToJSONString() string {
	bytes, err := s.ToJSON()
	if err != nil {
		return fmt.Sprintf(`{"content":[{"type":"text","text":"Error serializing output: %s"}],"isError":true}`, err.Error())
	}
	return string(bytes)
}

// BuildOutputSchema creates a JSON schema for tool output
func BuildOutputSchema(description string, properties map[string]PropertySchema, required []string) *ToolOutputSchema {
	props := make(map[string]map[string]interface{})

	for name, prop := range properties {
		propMap := map[string]interface{}{
			"type": prop.Type,
		}

		if prop.Description != "" {
			propMap["description"] = prop.Description
		}
		if prop.Format != "" {
			propMap["format"] = prop.Format
		}
		if len(prop.Enum) > 0 {
			propMap["enum"] = prop.Enum
		}
		if prop.Items != nil {
			propMap["items"] = prop.Items
		}
		if prop.Properties != nil {
			propMap["properties"] = prop.Properties
		}

		props[name] = propMap
	}

	return &ToolOutputSchema{
		Type:        "object",
		Properties:  props,
		Required:    required,
		Description: description,
	}
}

// PropertySchema defines a property in an output schema
type PropertySchema struct {
	Type        string                 // "string", "number", "integer", "boolean", "array", "object"
	Description string                 // Property description
	Format      string                 // Format (e.g., "date-time", "email", "uri")
	Enum        []string               // Allowed values
	Items       map[string]interface{} // For array types
	Properties  map[string]interface{} // For object types
}

// ExampleSchemas provides common output schemas
var ExampleSchemas = struct {
	SimpleResult    *ToolOutputSchema
	ListResult      *ToolOutputSchema
	OperationResult *ToolOutputSchema
}{
	SimpleResult: BuildOutputSchema(
		"A simple result with a single value",
		map[string]PropertySchema{
			"result": {
				Type:        "string",
				Description: "The operation result",
			},
		},
		[]string{"result"},
	),
	ListResult: BuildOutputSchema(
		"A list of items",
		map[string]PropertySchema{
			"items": {
				Type:        "array",
				Description: "List of items",
				Items: map[string]interface{}{
					"type": "object",
				},
			},
			"count": {
				Type:        "integer",
				Description: "Number of items",
			},
		},
		[]string{"items", "count"},
	),
	OperationResult: BuildOutputSchema(
		"Result of an operation with status",
		map[string]PropertySchema{
			"success": {
				Type:        "boolean",
				Description: "Whether the operation succeeded",
			},
			"message": {
				Type:        "string",
				Description: "Human-readable message",
			},
			"data": {
				Type:        "object",
				Description: "Additional data",
			},
		},
		[]string{"success", "message"},
	),
}
