package completion

import (
	"encoding/json"
	"testing"
)

func TestNewStructuredOutput(t *testing.T) {
	text := "Operation completed successfully"
	data := map[string]interface{}{
		"status": "success",
		"count":  42,
	}

	output := NewStructuredOutput(text, data)

	if len(output.Content) != 1 {
		t.Errorf("Expected 1 content item, got %d", len(output.Content))
	}

	if output.Content[0].Type != "text" {
		t.Errorf("Expected content type 'text', got '%s'", output.Content[0].Type)
	}

	if output.Content[0].Text != text {
		t.Errorf("Expected text '%s', got '%s'", text, output.Content[0].Text)
	}

	if output.StructuredContent["status"] != "success" {
		t.Errorf("Expected status 'success', got '%v'", output.StructuredContent["status"])
	}

	if output.StructuredContent["count"] != 42 {
		t.Errorf("Expected count 42, got '%v'", output.StructuredContent["count"])
	}

	if output.IsError {
		t.Error("Expected IsError to be false")
	}
}

func TestNewTextOutput(t *testing.T) {
	text := "Simple text output"
	output := NewTextOutput(text)

	if len(output.Content) != 1 {
		t.Errorf("Expected 1 content item, got %d", len(output.Content))
	}

	if output.Content[0].Text != text {
		t.Errorf("Expected text '%s', got '%s'", text, output.Content[0].Text)
	}

	if output.StructuredContent != nil {
		t.Error("Expected StructuredContent to be nil")
	}

	if output.IsError {
		t.Error("Expected IsError to be false")
	}
}

func TestNewErrorOutput(t *testing.T) {
	errorMsg := "Something went wrong"
	output := NewErrorOutput(errorMsg)

	if !output.IsError {
		t.Error("Expected IsError to be true")
	}

	expectedText := "Error: " + errorMsg
	if output.Content[0].Text != expectedText {
		t.Errorf("Expected text '%s', got '%s'", expectedText, output.Content[0].Text)
	}
}

func TestAddTextContent(t *testing.T) {
	output := NewTextOutput("Initial")
	output.AddTextContent("Additional text")

	if len(output.Content) != 2 {
		t.Errorf("Expected 2 content items, got %d", len(output.Content))
	}

	if output.Content[1].Text != "Additional text" {
		t.Errorf("Expected 'Additional text', got '%s'", output.Content[1].Text)
	}
}

func TestAddImageContent(t *testing.T) {
	output := NewTextOutput("Test")
	output.AddImageContent("base64data", "image/png")

	if len(output.Content) != 2 {
		t.Errorf("Expected 2 content items, got %d", len(output.Content))
	}

	if output.Content[1].Type != "image" {
		t.Errorf("Expected type 'image', got '%s'", output.Content[1].Type)
	}

	if output.Content[1].Data != "base64data" {
		t.Errorf("Expected data 'base64data', got '%s'", output.Content[1].Data)
	}
}

func TestAddResourceContent(t *testing.T) {
	output := NewTextOutput("Test")
	output.AddResourceContent("file:///path/to/resource")

	if len(output.Content) != 2 {
		t.Errorf("Expected 2 content items, got %d", len(output.Content))
	}

	if output.Content[1].Type != "resource" {
		t.Errorf("Expected type 'resource', got '%s'", output.Content[1].Type)
	}

	if output.Content[1].URI != "file:///path/to/resource" {
		t.Errorf("Expected URI 'file:///path/to/resource', got '%s'", output.Content[1].URI)
	}
}

func TestSetStructuredData(t *testing.T) {
	output := NewTextOutput("Test")
	data := map[string]interface{}{
		"key": "value",
	}
	output.SetStructuredData(data)

	if output.StructuredContent["key"] != "value" {
		t.Errorf("Expected key 'value', got '%v'", output.StructuredContent["key"])
	}
}

func TestToJSON(t *testing.T) {
	output := NewStructuredOutput("Test", map[string]interface{}{
		"status": "success",
	})

	jsonBytes, err := output.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON failed: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &parsed); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if parsed["isError"].(bool) != false {
		t.Error("Expected isError to be false")
	}
}

func TestToJSONString(t *testing.T) {
	output := NewStructuredOutput("Test", map[string]interface{}{
		"status": "success",
	})

	jsonStr := output.ToJSONString()
	if jsonStr == "" {
		t.Error("Expected non-empty JSON string")
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &parsed); err != nil {
		t.Fatalf("Failed to parse JSON string: %v", err)
	}
}

func TestBuildOutputSchema(t *testing.T) {
	schema := BuildOutputSchema(
		"Test schema",
		map[string]PropertySchema{
			"name": {
				Type:        "string",
				Description: "User name",
			},
			"age": {
				Type:        "integer",
				Description: "User age",
			},
		},
		[]string{"name"},
	)

	if schema.Type != "object" {
		t.Errorf("Expected type 'object', got '%s'", schema.Type)
	}

	if schema.Description != "Test schema" {
		t.Errorf("Expected description 'Test schema', got '%s'", schema.Description)
	}

	if len(schema.Properties) != 2 {
		t.Errorf("Expected 2 properties, got %d", len(schema.Properties))
	}

	if schema.Properties["name"]["type"] != "string" {
		t.Errorf("Expected name type 'string', got '%v'", schema.Properties["name"]["type"])
	}

	if len(schema.Required) != 1 || schema.Required[0] != "name" {
		t.Errorf("Expected required ['name'], got %v", schema.Required)
	}
}

func TestPropertySchemaWithEnum(t *testing.T) {
	schema := BuildOutputSchema(
		"Test",
		map[string]PropertySchema{
			"status": {
				Type:        "string",
				Description: "Status",
				Enum:        []string{"active", "inactive"},
			},
		},
		[]string{},
	)

	statusProp := schema.Properties["status"]
	enumVal, ok := statusProp["enum"]
	if !ok {
		t.Error("Expected enum property")
	}

	enumSlice, ok := enumVal.([]string)
	if !ok || len(enumSlice) != 2 {
		t.Errorf("Expected enum with 2 values, got %v", enumVal)
	}
}

func TestPropertySchemaWithFormat(t *testing.T) {
	schema := BuildOutputSchema(
		"Test",
		map[string]PropertySchema{
			"email": {
				Type:        "string",
				Description: "Email address",
				Format:      "email",
			},
		},
		[]string{},
	)

	emailProp := schema.Properties["email"]
	if emailProp["format"] != "email" {
		t.Errorf("Expected format 'email', got '%v'", emailProp["format"])
	}
}

func TestExampleSchemas(t *testing.T) {
	// Test SimpleResult
	if ExampleSchemas.SimpleResult == nil {
		t.Error("SimpleResult schema is nil")
	}

	if ExampleSchemas.SimpleResult.Type != "object" {
		t.Errorf("Expected SimpleResult type 'object', got '%s'", ExampleSchemas.SimpleResult.Type)
	}

	// Test ListResult
	if ExampleSchemas.ListResult == nil {
		t.Error("ListResult schema is nil")
	}

	// Test OperationResult
	if ExampleSchemas.OperationResult == nil {
		t.Error("OperationResult schema is nil")
	}

	if len(ExampleSchemas.OperationResult.Properties) < 2 {
		t.Error("OperationResult should have at least 2 properties")
	}
}

func TestComplexNestedData(t *testing.T) {
	output := NewStructuredOutput(
		"Complex data",
		map[string]interface{}{
			"user": map[string]interface{}{
				"id":   123,
				"name": "John",
			},
			"tags": []string{"go", "mcp", "test"},
		},
	)

	if output.StructuredContent == nil {
		t.Error("Expected structured content")
	}

	user, ok := output.StructuredContent["user"].(map[string]interface{})
	if !ok {
		t.Error("Expected user to be a map")
	}

	if user["id"] != 123 {
		t.Errorf("Expected user id 123, got %v", user["id"])
	}

	tags, ok := output.StructuredContent["tags"].([]string)
	if !ok || len(tags) != 3 {
		t.Errorf("Expected 3 tags, got %v", tags)
	}
}
