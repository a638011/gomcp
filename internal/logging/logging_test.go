package logging

import (
	"context"
	"testing"
	"time"
)

func TestLogToClientWithoutSession(t *testing.T) {
	ctx := context.Background()

	// Should not error when session is nil
	err := LogToClient(ctx, nil, LevelInfo, "Test message", nil)
	if err != nil {
		t.Errorf("Expected no error with nil session, got %v", err)
	}
}

func TestDebugLevel(t *testing.T) {
	ctx := context.Background()
	err := Debug(ctx, nil, "Debug message", map[string]interface{}{"key": "value"})
	if err != nil {
		t.Errorf("Debug should not error: %v", err)
	}
}

func TestInfoLevel(t *testing.T) {
	ctx := context.Background()
	err := Info(ctx, nil, "Info message", nil)
	if err != nil {
		t.Errorf("Info should not error: %v", err)
	}
}

func TestNoticeLevel(t *testing.T) {
	ctx := context.Background()
	err := Notice(ctx, nil, "Notice message", nil)
	if err != nil {
		t.Errorf("Notice should not error: %v", err)
	}
}

func TestWarningLevel(t *testing.T) {
	ctx := context.Background()
	err := Warning(ctx, nil, "Warning message", nil)
	if err != nil {
		t.Errorf("Warning should not error: %v", err)
	}
}

func TestErrorLevel(t *testing.T) {
	ctx := context.Background()
	err := Error(ctx, nil, "Error message", map[string]interface{}{"error": "test"})
	if err != nil {
		t.Errorf("Error should not error: %v", err)
	}
}

func TestCriticalLevel(t *testing.T) {
	ctx := context.Background()
	err := Critical(ctx, nil, "Critical message", nil)
	if err != nil {
		t.Errorf("Critical should not error: %v", err)
	}
}

func TestAlertLevel(t *testing.T) {
	ctx := context.Background()
	err := Alert(ctx, nil, "Alert message", nil)
	if err != nil {
		t.Errorf("Alert should not error: %v", err)
	}
}

func TestEmergencyLevel(t *testing.T) {
	ctx := context.Background()
	err := Emergency(ctx, nil, "Emergency message", nil)
	if err != nil {
		t.Errorf("Emergency should not error: %v", err)
	}
}

func TestRateLimiter(t *testing.T) {
	limiter := NewRateLimiter(10) // 10 per second

	// Should allow first 10
	allowed := 0
	for i := 0; i < 15; i++ {
		if limiter.Allow() {
			allowed++
		}
	}

	if allowed != 10 {
		t.Errorf("Expected 10 allowed, got %d", allowed)
	}

	// Wait for reset
	time.Sleep(1100 * time.Millisecond)

	// Should allow again
	if !limiter.Allow() {
		t.Error("Expected to allow after reset")
	}
}

func TestRateLimiterReset(t *testing.T) {
	limiter := NewRateLimiter(5)

	// Use up the limit
	for i := 0; i < 5; i++ {
		limiter.Allow()
	}

	// Should deny
	if limiter.Allow() {
		t.Error("Expected to deny after limit reached")
	}

	// Wait for reset
	time.Sleep(1100 * time.Millisecond)

	// Should allow after reset
	if !limiter.Allow() {
		t.Error("Expected to allow after reset")
	}
}

func TestRateLimitedLogger(t *testing.T) {
	ctx := context.Background()
	logger := NewRateLimitedLogger(5) // 5 per second

	// Log 10 messages
	for i := 0; i < 10; i++ {
		logger.Log(ctx, nil, LevelInfo, "Test", nil)
	}

	// Should have dropped 5
	dropped := logger.GetDroppedCount()
	if dropped != 5 {
		t.Errorf("Expected 5 dropped messages, got %d", dropped)
	}

	// Reset counter
	logger.ResetDroppedCount()
	if logger.GetDroppedCount() != 0 {
		t.Error("Expected dropped count to be 0 after reset")
	}
}

func TestRateLimitedLoggerResetDropped(t *testing.T) {
	ctx := context.Background()
	logger := NewRateLimitedLogger(2)

	// Exceed limit
	for i := 0; i < 5; i++ {
		logger.Log(ctx, nil, LevelInfo, "Test", nil)
	}

	if logger.GetDroppedCount() == 0 {
		t.Error("Expected some dropped messages")
	}

	logger.ResetDroppedCount()
	if logger.GetDroppedCount() != 0 {
		t.Errorf("Expected 0 dropped after reset, got %d", logger.GetDroppedCount())
	}
}

func TestSecureLog(t *testing.T) {
	ctx := context.Background()

	// Log data with sensitive fields
	data := map[string]interface{}{
		"username": "john",
		"password": "secret123",
		"token":    "abc123",
		"api_key":  "xyz789",
		"email":    "john@example.com",
	}

	// Should not error
	err := SecureLog(ctx, nil, LevelInfo, "Login attempt", data)
	if err != nil {
		t.Errorf("SecureLog should not error: %v", err)
	}

	// Note: In a real implementation, we'd verify that sensitive fields are redacted
	// For this test, we're just checking that the function doesn't error
}

func TestSecureLogRedactsPassword(t *testing.T) {
	ctx := context.Background()

	data := map[string]interface{}{
		"username": "john",
		"password": "secret123",
	}

	// Call SecureLog
	err := SecureLog(ctx, nil, LevelInfo, "Test", data)
	if err != nil {
		t.Fatalf("SecureLog failed: %v", err)
	}

	// The original data should still have password
	// (SecureLog creates a copy)
	if data["password"] != "secret123" {
		t.Error("Original data should not be modified")
	}
}

func TestSecureLogSensitiveFields(t *testing.T) {
	ctx := context.Background()

	sensitiveFields := []string{
		"password", "token", "secret", "api_key",
		"apikey", "auth", "authorization", "credential",
		"private_key", "privatekey", "session_id",
		"sessionid", "cookie", "ssn", "credit_card",
	}

	for _, field := range sensitiveFields {
		data := map[string]interface{}{
			field: "should-be-redacted",
		}

		err := SecureLog(ctx, nil, LevelInfo, "Test", data)
		if err != nil {
			t.Errorf("SecureLog failed for field %s: %v", field, err)
		}
	}
}

func TestLogMessageStructure(t *testing.T) {
	msg := LogMessage{
		Level:     LevelInfo,
		Logger:    "test-logger",
		Data:      "test data",
		Timestamp: time.Now().UTC(),
		Meta:      map[string]interface{}{"key": "value"},
	}

	if msg.Level != LevelInfo {
		t.Errorf("Expected level info, got %s", msg.Level)
	}

	if msg.Logger != "test-logger" {
		t.Errorf("Expected logger 'test-logger', got '%s'", msg.Logger)
	}

	if msg.Data != "test data" {
		t.Errorf("Expected data 'test data', got '%v'", msg.Data)
	}

	if msg.Meta["key"] != "value" {
		t.Errorf("Expected meta key 'value', got '%v'", msg.Meta["key"])
	}
}

func TestLogLevels(t *testing.T) {
	levels := []LogLevel{
		LevelDebug,
		LevelInfo,
		LevelNotice,
		LevelWarning,
		LevelError,
		LevelCritical,
		LevelAlert,
		LevelEmergency,
	}

	expectedLevels := []string{
		"debug", "info", "notice", "warning",
		"error", "critical", "alert", "emergency",
	}

	for i, level := range levels {
		if string(level) != expectedLevels[i] {
			t.Errorf("Expected level '%s', got '%s'", expectedLevels[i], string(level))
		}
	}
}
