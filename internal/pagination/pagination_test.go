package pagination

import (
	"testing"
)

func TestEncodeDecodeCursor(t *testing.T) {
	pageInfo := PageInfo{
		Page:       2,
		PerPage:    10,
		TotalItems: 100,
	}

	// Encode
	cursor, err := EncodeCursor(pageInfo)
	if err != nil {
		t.Fatalf("EncodeCursor failed: %v", err)
	}

	if cursor == "" {
		t.Error("Expected non-empty cursor")
	}

	// Decode
	decoded, err := DecodeCursor(cursor)
	if err != nil {
		t.Fatalf("DecodeCursor failed: %v", err)
	}

	if decoded.Page != pageInfo.Page {
		t.Errorf("Expected page %d, got %d", pageInfo.Page, decoded.Page)
	}

	if decoded.PerPage != pageInfo.PerPage {
		t.Errorf("Expected perPage %d, got %d", pageInfo.PerPage, decoded.PerPage)
	}

	if decoded.TotalItems != pageInfo.TotalItems {
		t.Errorf("Expected totalItems %d, got %d", pageInfo.TotalItems, decoded.TotalItems)
	}
}

func TestDecodeEmptyCursor(t *testing.T) {
	decoded, err := DecodeCursor("")
	if err != nil {
		t.Fatalf("DecodeCursor failed: %v", err)
	}

	if decoded.Page != 1 {
		t.Errorf("Expected default page 1, got %d", decoded.Page)
	}

	if decoded.PerPage != 10 {
		t.Errorf("Expected default perPage 10, got %d", decoded.PerPage)
	}
}

func TestDecodeInvalidCursor(t *testing.T) {
	_, err := DecodeCursor("invalid-cursor-string")
	if err == nil {
		t.Error("Expected error for invalid cursor")
	}
}

func TestPaginateSlice(t *testing.T) {
	items := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// First page (3 items per page)
	response, err := PaginateSlice(items, "", 3)
	if err != nil {
		t.Fatalf("PaginateSlice failed: %v", err)
	}

	pageItems, ok := response.Items.([]int)
	if !ok {
		t.Fatal("Expected items to be []int")
	}

	if len(pageItems) != 3 {
		t.Errorf("Expected 3 items on page 1, got %d", len(pageItems))
	}

	if pageItems[0] != 1 || pageItems[1] != 2 || pageItems[2] != 3 {
		t.Errorf("Expected [1, 2, 3], got %v", pageItems)
	}

	if response.NextCursor == nil {
		t.Error("Expected nextCursor for more pages")
	}

	// Second page
	response2, err := PaginateSlice(items, Cursor(*response.NextCursor), 3)
	if err != nil {
		t.Fatalf("PaginateSlice page 2 failed: %v", err)
	}

	pageItems2, ok := response2.Items.([]int)
	if !ok {
		t.Fatal("Expected items to be []int")
	}

	if len(pageItems2) != 3 {
		t.Errorf("Expected 3 items on page 2, got %d", len(pageItems2))
	}

	if pageItems2[0] != 4 || pageItems2[1] != 5 || pageItems2[2] != 6 {
		t.Errorf("Expected [4, 5, 6], got %v", pageItems2)
	}
}

func TestPaginateSliceLastPage(t *testing.T) {
	items := []int{1, 2, 3, 4, 5}

	// Page with 3 items
	response, err := PaginateSlice(items, "", 3)
	if err != nil {
		t.Fatalf("PaginateSlice failed: %v", err)
	}

	// Second page (last page with 2 items)
	response2, err := PaginateSlice(items, Cursor(*response.NextCursor), 3)
	if err != nil {
		t.Fatalf("PaginateSlice page 2 failed: %v", err)
	}

	pageItems2, ok := response2.Items.([]int)
	if !ok {
		t.Fatal("Expected items to be []int")
	}

	if len(pageItems2) != 2 {
		t.Errorf("Expected 2 items on last page, got %d", len(pageItems2))
	}

	if response2.NextCursor != nil {
		t.Error("Expected no nextCursor on last page")
	}
}

func TestPaginateSliceEmpty(t *testing.T) {
	items := []int{}

	response, err := PaginateSlice(items, "", 10)
	if err != nil {
		t.Fatalf("PaginateSlice failed: %v", err)
	}

	pageItems, ok := response.Items.([]int)
	if !ok {
		t.Fatal("Expected items to be []int")
	}

	if len(pageItems) != 0 {
		t.Errorf("Expected 0 items, got %d", len(pageItems))
	}

	if response.NextCursor != nil {
		t.Error("Expected no nextCursor for empty results")
	}
}

func TestPaginateSliceBeyondEnd(t *testing.T) {
	items := []int{1, 2, 3}

	// Create cursor for page 10 (beyond the data)
	pageInfo := PageInfo{
		Page:    10,
		PerPage: 10,
	}
	cursor, _ := EncodeCursor(pageInfo)

	response, err := PaginateSlice(items, cursor, 10)
	if err != nil {
		t.Fatalf("PaginateSlice failed: %v", err)
	}

	pageItems, ok := response.Items.([]int)
	if !ok {
		t.Fatal("Expected items to be []int")
	}

	if len(pageItems) != 0 {
		t.Errorf("Expected 0 items beyond end, got %d", len(pageItems))
	}

	if response.NextCursor != nil {
		t.Error("Expected no nextCursor beyond end")
	}
}

func TestNewPaginatedResponse(t *testing.T) {
	items := []string{"a", "b", "c"}

	// With more pages
	response, err := NewPaginatedResponse(items, true, 2, 10)
	if err != nil {
		t.Fatalf("NewPaginatedResponse failed: %v", err)
	}

	if response.NextCursor == nil {
		t.Error("Expected nextCursor when hasMore is true")
	}

	// Last page
	response2, err := NewPaginatedResponse(items, false, 3, 10)
	if err != nil {
		t.Fatalf("NewPaginatedResponse failed: %v", err)
	}

	if response2.NextCursor != nil {
		t.Error("Expected no nextCursor when hasMore is false")
	}
}

func TestValidateCursor(t *testing.T) {
	// Valid cursor
	pageInfo := PageInfo{Page: 1, PerPage: 10}
	cursor, _ := EncodeCursor(pageInfo)

	err := ValidateCursor(cursor)
	if err != nil {
		t.Errorf("Expected valid cursor, got error: %v", err)
	}

	// Invalid cursor
	err = ValidateCursor("invalid")
	if err == nil {
		t.Error("Expected error for invalid cursor")
	}

	// Empty cursor (should be valid - defaults to first page)
	err = ValidateCursor("")
	if err != nil {
		t.Errorf("Expected empty cursor to be valid, got error: %v", err)
	}
}

func TestGetPageInfo(t *testing.T) {
	pageInfo := PageInfo{
		Page:       3,
		PerPage:    20,
		TotalItems: 100,
	}

	cursor, _ := EncodeCursor(pageInfo)
	retrieved, err := GetPageInfo(cursor)

	if err != nil {
		t.Fatalf("GetPageInfo failed: %v", err)
	}

	if retrieved.Page != pageInfo.Page {
		t.Errorf("Expected page %d, got %d", pageInfo.Page, retrieved.Page)
	}

	if retrieved.PerPage != pageInfo.PerPage {
		t.Errorf("Expected perPage %d, got %d", pageInfo.PerPage, retrieved.PerPage)
	}
}

func TestClampPageSize(t *testing.T) {
	tests := []struct {
		input    int
		expected int
	}{
		{0, DefaultPageSize},  // Zero becomes default
		{-1, DefaultPageSize}, // Negative becomes default
		{5, 5},                // Valid size stays same
		{50, 50},              // Valid size stays same
		{100, 100},            // Max size stays same
		{150, MaxPageSize},    // Over max becomes max
		{1000, MaxPageSize},   // Way over max becomes max
	}

	for _, test := range tests {
		result := ClampPageSize(test.input)
		if result != test.expected {
			t.Errorf("ClampPageSize(%d) = %d, expected %d", test.input, result, test.expected)
		}
	}
}

func TestDefaultPageSize(t *testing.T) {
	if DefaultPageSize != 10 {
		t.Errorf("Expected DefaultPageSize to be 10, got %d", DefaultPageSize)
	}
}

func TestMaxPageSize(t *testing.T) {
	if MaxPageSize != 100 {
		t.Errorf("Expected MaxPageSize to be 100, got %d", MaxPageSize)
	}
}

func TestPaginateSliceWithCustomPageSize(t *testing.T) {
	items := make([]int, 50)
	for i := range items {
		items[i] = i + 1
	}

	// Test with page size 15
	response, err := PaginateSlice(items, "", 15)
	if err != nil {
		t.Fatalf("PaginateSlice failed: %v", err)
	}

	pageItems, ok := response.Items.([]int)
	if !ok {
		t.Fatal("Expected items to be []int")
	}

	if len(pageItems) != 15 {
		t.Errorf("Expected 15 items, got %d", len(pageItems))
	}

	// Verify pagination continues correctly
	response2, err := PaginateSlice(items, Cursor(*response.NextCursor), 15)
	if err != nil {
		t.Fatalf("PaginateSlice page 2 failed: %v", err)
	}

	pageItems2, ok := response2.Items.([]int)
	if !ok {
		t.Fatal("Expected items to be []int")
	}

	if pageItems2[0] != 16 {
		t.Errorf("Expected first item on page 2 to be 16, got %d", pageItems2[0])
	}
}

func TestPaginateSliceStructs(t *testing.T) {
	type Item struct {
		ID   int
		Name string
	}

	items := []Item{
		{1, "One"},
		{2, "Two"},
		{3, "Three"},
		{4, "Four"},
		{5, "Five"},
	}

	response, err := PaginateSlice(items, "", 2)
	if err != nil {
		t.Fatalf("PaginateSlice failed: %v", err)
	}

	pageItems, ok := response.Items.([]Item)
	if !ok {
		t.Fatal("Expected items to be []Item")
	}

	if len(pageItems) != 2 {
		t.Errorf("Expected 2 items, got %d", len(pageItems))
	}

	if pageItems[0].ID != 1 || pageItems[0].Name != "One" {
		t.Errorf("Expected {1, One}, got {%d, %s}", pageItems[0].ID, pageItems[0].Name)
	}
}
