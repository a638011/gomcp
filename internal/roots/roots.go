package roots

import (
	"context"

	mcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

// ListRoots returns the list of filesystem roots the server exposes
func ListRoots(ctx context.Context, req *mcp.ListRootsRequest) (*mcp.ListRootsResult, error) {
	// Define the roots that clients can access
	roots := []*mcp.Root{
		{
			URI:  "file:///tmp",
			Name: "Temporary Directory",
		},
		{
			URI:  "file:///var/log",
			Name: "System Logs",
		},
		{
			URI:  "file://" + getCurrentWorkingDir(),
			Name: "Server Working Directory",
		},
	}

	return &mcp.ListRootsResult{
		Roots: roots,
	}, nil
}

// getCurrentWorkingDir returns the current working directory
func getCurrentWorkingDir() string {
	// This would normally use os.Getwd()
	// For now, return a placeholder
	return "/app"
}
