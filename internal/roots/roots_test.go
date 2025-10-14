package roots

import (
	"context"
	"testing"
)

func TestListRoots(t *testing.T) {
	ctx := context.Background()

	result, err := ListRoots(ctx, nil)
	if err != nil {
		t.Fatalf("ListRoots failed: %v", err)
	}

	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	if len(result.Roots) == 0 {
		t.Error("Expected at least one root")
	}

	// Check that roots have required fields
	for i, root := range result.Roots {
		if root.URI == "" {
			t.Errorf("Root %d has empty URI", i)
		}

		if root.Name == "" {
			t.Errorf("Root %d has empty Name", i)
		}
	}
}

func TestListRootsURIFormat(t *testing.T) {
	ctx := context.Background()

	result, err := ListRoots(ctx, nil)
	if err != nil {
		t.Fatalf("ListRoots failed: %v", err)
	}

	// All URIs should start with file://
	for i, root := range result.Roots {
		if len(root.URI) < 7 || root.URI[:7] != "file://" {
			t.Errorf("Root %d URI should start with 'file://', got '%s'", i, root.URI)
		}
	}
}

func TestListRootsExpectedRoots(t *testing.T) {
	ctx := context.Background()

	result, err := ListRoots(ctx, nil)
	if err != nil {
		t.Fatalf("ListRoots failed: %v", err)
	}

	expectedURIs := []string{
		"file:///tmp",
		"file:///var/log",
		"file:///app",
	}

	if len(result.Roots) != len(expectedURIs) {
		t.Errorf("Expected %d roots, got %d", len(expectedURIs), len(result.Roots))
	}

	// Check for expected URIs
	foundURIs := make(map[string]bool)
	for _, root := range result.Roots {
		foundURIs[root.URI] = true
	}

	for _, expectedURI := range expectedURIs {
		if !foundURIs[expectedURI] {
			t.Errorf("Expected to find root with URI '%s'", expectedURI)
		}
	}
}

func TestListRootsNames(t *testing.T) {
	ctx := context.Background()

	result, err := ListRoots(ctx, nil)
	if err != nil {
		t.Fatalf("ListRoots failed: %v", err)
	}

	expectedNames := []string{
		"Temporary Directory",
		"System Logs",
		"Server Working Directory",
	}

	foundNames := make(map[string]bool)
	for _, root := range result.Roots {
		foundNames[root.Name] = true
	}

	for _, expectedName := range expectedNames {
		if !foundNames[expectedName] {
			t.Errorf("Expected to find root with name '%s'", expectedName)
		}
	}
}
