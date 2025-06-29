package errors

import (
	"errors"
	"fmt"
	"testing"
)

func TestErrorType_String(t *testing.T) {
	testCases := []struct {
		errorType ErrorType
		expected  string
	}{
		{ErrorTypeConfiguration, "Configuration"},
		{ErrorTypeConnection, "Connection"},
		{ErrorTypeValidation, "Validation"},
		{ErrorTypePermission, "Permission"},
		{ErrorTypeFileSystem, "FileSystem"},
		{ErrorTypeUI, "UI"},
		{ErrorTypeMCP, "MCP"},
		{ErrorTypeUnknown, "Unknown"},
	}
	
	for _, tc := range testCases {
		result := tc.errorType.String()
		if result != tc.expected {
			t.Errorf("ErrorType %d: expected '%s', got '%s'", tc.errorType, tc.expected, result)
		}
	}
}

func TestAppError_Error(t *testing.T) {
	// Test error without cause
	err1 := &AppError{
		Type:    ErrorTypeConfiguration,
		Code:    "TEST_CODE",
		Message: "Test message",
	}
	
	expected1 := "[Configuration:TEST_CODE] Test message"
	if err1.Error() != expected1 {
		t.Errorf("Expected '%s', got '%s'", expected1, err1.Error())
	}
	
	// Test error with cause
	cause := errors.New("underlying error")
	err2 := &AppError{
		Type:    ErrorTypeConnection,
		Code:    "TEST_CODE",
		Message: "Test message",
		Cause:   cause,
	}
	
	expected2 := "[Connection:TEST_CODE] Test message: underlying error"
	if err2.Error() != expected2 {
		t.Errorf("Expected '%s', got '%s'", expected2, err2.Error())
	}
}

func TestAppError_Unwrap(t *testing.T) {
	originalErr := errors.New("original error")
	appErr := &AppError{
		Type:  ErrorTypeValidation,
		Code:  "TEST",
		Cause: originalErr,
	}
	
	unwrapped := appErr.Unwrap()
	if unwrapped != originalErr {
		t.Errorf("Expected unwrapped error to be original error")
	}
	
	// Test with no cause
	appErrNoCause := &AppError{
		Type: ErrorTypeValidation,
		Code: "TEST",
	}
	
	if appErrNoCause.Unwrap() != nil {
		t.Errorf("Expected nil when unwrapping error with no cause")
	}
}

func TestAppError_GetUserMessage(t *testing.T) {
	// Test with explicit user message
	err1 := &AppError{
		Type:        ErrorTypeConfiguration,
		Code:        "TEST",
		Message:     "Technical message",
		UserMessage: "User-friendly message",
	}
	
	if err1.GetUserMessage() != "User-friendly message" {
		t.Errorf("Expected user message, got '%s'", err1.GetUserMessage())
	}
	
	// Test without user message (should return technical message)
	err2 := &AppError{
		Type:    ErrorTypeConfiguration,
		Code:    "TEST",
		Message: "Technical message",
	}
	
	if err2.GetUserMessage() != "Technical message" {
		t.Errorf("Expected technical message, got '%s'", err2.GetUserMessage())
	}
}

func TestAppError_WithContext(t *testing.T) {
	err := &AppError{
		Type: ErrorTypeValidation,
		Code: "TEST",
	}
	
	err.WithContext("field", "username")
	err.WithContext("value", "invalid_user")
	
	if field, exists := err.GetContext("field"); !exists || field != "username" {
		t.Errorf("Expected context field 'username', got '%v' (exists: %t)", field, exists)
	}
	
	if value, exists := err.GetContext("value"); !exists || value != "invalid_user" {
		t.Errorf("Expected context value 'invalid_user', got '%v' (exists: %t)", value, exists)
	}
	
	if _, exists := err.GetContext("nonexistent"); exists {
		t.Error("Expected false for non-existent context key")
	}
}

func TestNew(t *testing.T) {
	err := New(ErrorTypeConfiguration, "TEST_CODE", "Test message")
	
	if err.Type != ErrorTypeConfiguration {
		t.Errorf("Expected type Configuration, got %v", err.Type)
	}
	
	if err.Code != "TEST_CODE" {
		t.Errorf("Expected code 'TEST_CODE', got '%s'", err.Code)
	}
	
	if err.Message != "Test message" {
		t.Errorf("Expected message 'Test message', got '%s'", err.Message)
	}
	
	if err.Cause != nil {
		t.Errorf("Expected no cause, got %v", err.Cause)
	}
}

func TestWrap(t *testing.T) {
	originalErr := errors.New("original error")
	err := Wrap(originalErr, ErrorTypeConnection, "WRAP_TEST", "Wrapped message")
	
	if err.Type != ErrorTypeConnection {
		t.Errorf("Expected type Connection, got %v", err.Type)
	}
	
	if err.Code != "WRAP_TEST" {
		t.Errorf("Expected code 'WRAP_TEST', got '%s'", err.Code)
	}
	
	if err.Message != "Wrapped message" {
		t.Errorf("Expected message 'Wrapped message', got '%s'", err.Message)
	}
	
	if err.Cause != originalErr {
		t.Errorf("Expected cause to be original error")
	}
}

func TestWithUserMessage(t *testing.T) {
	err := New(ErrorTypeValidation, "TEST", "Technical message")
	WithUserMessage(err, "User-friendly message")
	
	if err.UserMessage != "User-friendly message" {
		t.Errorf("Expected user message 'User-friendly message', got '%s'", err.UserMessage)
	}
}

func TestGetUserFriendlyMessage(t *testing.T) {
	// Test existing code
	msg := GetUserFriendlyMessage("CONFIG_NOT_FOUND")
	if msg != "No configuration file found. Using default settings." {
		t.Errorf("Unexpected message for CONFIG_NOT_FOUND: %s", msg)
	}
	
	// Test non-existent code
	msg = GetUserFriendlyMessage("NON_EXISTENT_CODE")
	if msg != "An unexpected error occurred. Please try again." {
		t.Errorf("Unexpected default message: %s", msg)
	}
}

func TestFormatErrorForUser(t *testing.T) {
	testCases := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name: "AppError with user message",
			err: &AppError{
				Type:        ErrorTypeConfiguration,
				Code:        "TEST",
				UserMessage: "Custom user message",
			},
			expected: "Custom user message",
		},
		{
			name: "AppError without user message",
			err: &AppError{
				Type: ErrorTypeConfiguration,
				Code: "CONFIG_NOT_FOUND",
			},
			expected: "No configuration file found. Using default settings.",
		},
		{
			name:     "Standard Go error - file not found",
			err:      errors.New("open config.json: no such file or directory"),
			expected: "The requested file could not be found.",
		},
		{
			name:     "Standard Go error - permission denied",
			err:      errors.New("open config.json: permission denied"),
			expected: "Permission denied. Please check file permissions.",
		},
		{
			name:     "Standard Go error - connection refused",
			err:      errors.New("dial tcp :8080: connection refused"),
			expected: "Unable to connect to the server. Please check if it's running.",
		},
		{
			name:     "Standard Go error - timeout",
			err:      errors.New("context deadline exceeded: timeout"),
			expected: "Operation timed out. Please try again.",
		},
		{
			name:     "Standard Go error - network",
			err:      errors.New("network is unreachable"),
			expected: "Network error occurred. Please check your connection.",
		},
		{
			name:     "Unknown error",
			err:      errors.New("some unknown error"),
			expected: "An unexpected error occurred. Please try again.",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := FormatErrorForUser(tc.err)
			if result != tc.expected {
				t.Errorf("Expected '%s', got '%s'", tc.expected, result)
			}
		})
	}
}

func TestGetRecoveryActions(t *testing.T) {
	testCases := []struct {
		name            string
		err             error
		expectedActions int
		firstAction     string
	}{
		{
			name: "CONFIG_NOT_FOUND",
			err: &AppError{
				Type: ErrorTypeConfiguration,
				Code: "CONFIG_NOT_FOUND",
			},
			expectedActions: 1,
			firstAction:     "create_default_config",
		},
		{
			name: "CONFIG_INVALID",
			err: &AppError{
				Type: ErrorTypeConfiguration,
				Code: "CONFIG_INVALID",
			},
			expectedActions: 2,
			firstAction:     "reset_config",
		},
		{
			name: "CONNECTION_TIMEOUT",
			err: &AppError{
				Type: ErrorTypeConnection,
				Code: "CONNECTION_TIMEOUT",
			},
			expectedActions: 2,
			firstAction:     "retry_connection",
		},
		{
			name: "TERMINAL_TOO_SMALL",
			err: &AppError{
				Type: ErrorTypeUI,
				Code: "TERMINAL_TOO_SMALL",
			},
			expectedActions: 1,
			firstAction:     "resize_terminal",
		},
		{
			name:            "Unknown error",
			err:             errors.New("unknown error"),
			expectedActions: 1,
			firstAction:     "retry",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actions := GetRecoveryActions(tc.err)
			
			if len(actions) != tc.expectedActions {
				t.Errorf("Expected %d actions, got %d", tc.expectedActions, len(actions))
			}
			
			if len(actions) > 0 && actions[0].Action != tc.firstAction {
				t.Errorf("Expected first action '%s', got '%s'", tc.firstAction, actions[0].Action)
			}
		})
	}
}

func TestIsRecoverable(t *testing.T) {
	testCases := []struct {
		name       string
		err        error
		recoverable bool
	}{
		{
			name: "CONFIG_NOT_FOUND - recoverable",
			err: &AppError{
				Type: ErrorTypeConfiguration,
				Code: "CONFIG_NOT_FOUND",
			},
			recoverable: true,
		},
		{
			name: "CONNECTION_TIMEOUT - recoverable",
			err: &AppError{
				Type: ErrorTypeConnection,
				Code: "CONNECTION_TIMEOUT",
			},
			recoverable: true,
		},
		{
			name: "CONFIG_INVALID - not recoverable",
			err: &AppError{
				Type: ErrorTypeConfiguration,
				Code: "CONFIG_INVALID",
			},
			recoverable: false,
		},
		{
			name:        "Standard Go error - not recoverable",
			err:         errors.New("some error"),
			recoverable: false,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := IsRecoverable(tc.err)
			if result != tc.recoverable {
				t.Errorf("Expected recoverable=%t, got %t", tc.recoverable, result)
			}
		})
	}
}

func TestIsCritical(t *testing.T) {
	testCases := []struct {
		name     string
		err      error
		critical bool
	}{
		{
			name: "Permission error - critical",
			err: &AppError{
				Type: ErrorTypePermission,
				Code: "ACCESS_DENIED",
			},
			critical: true,
		},
		{
			name: "Config validation error - critical",
			err: &AppError{
				Type: ErrorTypeConfiguration,
				Code: "CONFIG_VALIDATION",
			},
			critical: true,
		},
		{
			name: "Terminal not supported - critical",
			err: &AppError{
				Type: ErrorTypeUI,
				Code: "TERMINAL_NOT_SUPPORTED",
			},
			critical: true,
		},
		{
			name: "Connection timeout - not critical",
			err: &AppError{
				Type: ErrorTypeConnection,
				Code: "CONNECTION_TIMEOUT",
			},
			critical: false,
		},
		{
			name:     "Standard Go error - not critical",
			err:      errors.New("some error"),
			critical: false,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := IsCritical(tc.err)
			if result != tc.critical {
				t.Errorf("Expected critical=%t, got %t", tc.critical, result)
			}
		})
	}
}

func TestPredefinedErrors(t *testing.T) {
	// Test that predefined errors have correct types and codes
	testCases := []struct {
		name      string
		err       *AppError
		errorType ErrorType
		code      string
	}{
		{"ErrConfigNotFound", ErrConfigNotFound, ErrorTypeConfiguration, "CONFIG_NOT_FOUND"},
		{"ErrConfigInvalid", ErrConfigInvalid, ErrorTypeConfiguration, "CONFIG_INVALID"},
		{"ErrConnectionTimeout", ErrConnectionTimeout, ErrorTypeConnection, "CONNECTION_TIMEOUT"},
		{"ErrInvalidInput", ErrInvalidInput, ErrorTypeValidation, "INVALID_INPUT"},
		{"ErrAccessDenied", ErrAccessDenied, ErrorTypePermission, "ACCESS_DENIED"},
		{"ErrFileNotFound", ErrFileNotFound, ErrorTypeFileSystem, "FILE_NOT_FOUND"},
		{"ErrTerminalTooSmall", ErrTerminalTooSmall, ErrorTypeUI, "TERMINAL_TOO_SMALL"},
		{"ErrMCPServerNotFound", ErrMCPServerNotFound, ErrorTypeMCP, "MCP_SERVER_NOT_FOUND"},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.err.Type != tc.errorType {
				t.Errorf("Expected type %v, got %v", tc.errorType, tc.err.Type)
			}
			
			if tc.err.Code != tc.code {
				t.Errorf("Expected code '%s', got '%s'", tc.code, tc.err.Code)
			}
			
			if tc.err.Message == "" {
				t.Error("Expected non-empty message")
			}
		})
	}
}

func ExampleAppError() {
	// Create a new configuration error
	err := New(ErrorTypeConfiguration, "INVALID_VALUE", "The value is invalid")
	
	// Add context
	err.WithContext("field", "timeout")
	err.WithContext("value", -1)
	
	// Add user message
	WithUserMessage(err, "Please enter a positive number for timeout")
	
	fmt.Println("Error:", err.Error())
	fmt.Println("User message:", err.GetUserMessage())
	
	// Check if error is recoverable
	fmt.Println("Recoverable:", IsRecoverable(err))
	
	// Output:
	// Error: [Configuration:INVALID_VALUE] The value is invalid
	// User message: Please enter a positive number for timeout
	// Recoverable: false
}