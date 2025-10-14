package pagination

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

// PageInfo contains pagination metadata
type PageInfo struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	TotalItems int `json:"total_items,omitempty"`
}

// Cursor represents an opaque pagination cursor
type Cursor string

// PaginatedResponse wraps any response with pagination info
type PaginatedResponse struct {
	Items      interface{} `json:"items"`                // The actual data
	NextCursor *string     `json:"nextCursor,omitempty"` // Pointer to allow null
}

// EncodeCursor creates an opaque cursor from page info
func EncodeCursor(pageInfo PageInfo) (Cursor, error) {
	jsonData, err := json.Marshal(pageInfo)
	if err != nil {
		return "", fmt.Errorf("failed to marshal cursor: %w", err)
	}

	// Base64 encode to make it opaque
	encoded := base64.StdEncoding.EncodeToString(jsonData)
	return Cursor(encoded), nil
}

// DecodeCursor decodes an opaque cursor back to page info
func DecodeCursor(cursor Cursor) (*PageInfo, error) {
	if cursor == "" {
		// Empty cursor means first page
		return &PageInfo{
			Page:    1,
			PerPage: 10, // Default
		}, nil
	}

	// Base64 decode
	jsonData, err := base64.StdEncoding.DecodeString(string(cursor))
	if err != nil {
		return nil, fmt.Errorf("invalid cursor format: %w", err)
	}

	// Unmarshal JSON
	var pageInfo PageInfo
	if err := json.Unmarshal(jsonData, &pageInfo); err != nil {
		return nil, fmt.Errorf("failed to parse cursor: %w", err)
	}

	return &pageInfo, nil
}

// Paginate creates a paginated response from a slice of items
func Paginate(items interface{}, cursor Cursor, perPage int) (*PaginatedResponse, error) {
	// Decode cursor
	pageInfo, err := DecodeCursor(cursor)
	if err != nil {
		return nil, err
	}

	// Update perPage if provided
	if perPage > 0 {
		pageInfo.PerPage = perPage
	}

	// Convert items to slice (we need reflection here)
	// For simplicity, assume items is already a slice
	itemsSlice, ok := items.([]interface{})
	if !ok {
		return nil, fmt.Errorf("items must be a slice")
	}

	totalItems := len(itemsSlice)
	startIdx := (pageInfo.Page - 1) * pageInfo.PerPage
	endIdx := startIdx + pageInfo.PerPage

	// Check if we're past the end
	if startIdx >= totalItems {
		return &PaginatedResponse{
			Items:      []interface{}{},
			NextCursor: nil, // No more pages
		}, nil
	}

	// Adjust end index if needed
	if endIdx > totalItems {
		endIdx = totalItems
	}

	// Get page of items
	pageItems := itemsSlice[startIdx:endIdx]

	// Create next cursor if there are more items
	var nextCursor *string
	if endIdx < totalItems {
		nextPageInfo := PageInfo{
			Page:       pageInfo.Page + 1,
			PerPage:    pageInfo.PerPage,
			TotalItems: totalItems,
		}
		nextCursorStr, err := EncodeCursor(nextPageInfo)
		if err == nil {
			cursorStr := string(nextCursorStr)
			nextCursor = &cursorStr
		}
	}

	return &PaginatedResponse{
		Items:      pageItems,
		NextCursor: nextCursor,
	}, nil
}

// PaginateSlice is a generic helper for paginating any slice
func PaginateSlice[T any](items []T, cursor Cursor, perPage int) (*PaginatedResponse, error) {
	// Decode cursor
	pageInfo, err := DecodeCursor(cursor)
	if err != nil {
		return nil, err
	}

	// Update perPage if provided
	if perPage > 0 {
		pageInfo.PerPage = perPage
	}

	// Default perPage if not set
	if pageInfo.PerPage == 0 {
		pageInfo.PerPage = 10
	}

	totalItems := len(items)
	startIdx := (pageInfo.Page - 1) * pageInfo.PerPage
	endIdx := startIdx + pageInfo.PerPage

	// Check if we're past the end
	if startIdx >= totalItems {
		return &PaginatedResponse{
			Items:      []T{},
			NextCursor: nil,
		}, nil
	}

	// Adjust end index if needed
	if endIdx > totalItems {
		endIdx = totalItems
	}

	// Get page of items
	pageItems := items[startIdx:endIdx]

	// Create next cursor if there are more items
	var nextCursor *string
	if endIdx < totalItems {
		nextPageInfo := PageInfo{
			Page:       pageInfo.Page + 1,
			PerPage:    pageInfo.PerPage,
			TotalItems: totalItems,
		}
		nextCursorStr, err := EncodeCursor(nextPageInfo)
		if err == nil {
			cursorStr := string(nextCursorStr)
			nextCursor = &cursorStr
		}
	}

	return &PaginatedResponse{
		Items:      pageItems,
		NextCursor: nextCursor,
	}, nil
}

// NewPaginatedResponse creates a paginated response manually
func NewPaginatedResponse(items interface{}, hasMore bool, nextPage int, perPage int) (*PaginatedResponse, error) {
	var nextCursor *string

	if hasMore {
		nextPageInfo := PageInfo{
			Page:    nextPage,
			PerPage: perPage,
		}
		cursorStr, err := EncodeCursor(nextPageInfo)
		if err != nil {
			return nil, err
		}
		cursor := string(cursorStr)
		nextCursor = &cursor
	}

	return &PaginatedResponse{
		Items:      items,
		NextCursor: nextCursor,
	}, nil
}

// ValidateCursor checks if a cursor is valid
func ValidateCursor(cursor Cursor) error {
	_, err := DecodeCursor(cursor)
	return err
}

// GetPageInfo extracts page information from a cursor
func GetPageInfo(cursor Cursor) (*PageInfo, error) {
	return DecodeCursor(cursor)
}

// DefaultPageSize is the default number of items per page
const DefaultPageSize = 10

// MaxPageSize is the maximum allowed page size
const MaxPageSize = 100

// ClampPageSize ensures page size is within bounds
func ClampPageSize(size int) int {
	if size <= 0 {
		return DefaultPageSize
	}
	if size > MaxPageSize {
		return MaxPageSize
	}
	return size
}
