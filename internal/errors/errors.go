// Package errors provides centralized error handling and user-friendly error messages
// for the MCP Manager CLI application.
package errors

import (
	"fmt"
	"strings"
)

// ErrorType represents different categories of errors
type ErrorType int

const (
	// ErrorTypeUnknown represents an unknown error type
	ErrorTypeUnknown ErrorType = iota
	
	// ErrorTypeConfiguration represents configuration-related errors
	ErrorTypeConfiguration
	
	// ErrorTypeConnection represents network/connection errors
	ErrorTypeConnection
	
	// ErrorTypeValidation represents input validation errors
	ErrorTypeValidation
	
	// ErrorTypePermission represents permission/access errors
	ErrorTypePermission
	
	// ErrorTypeFileSystem represents file system errors
	ErrorTypeFileSystem
	
	// ErrorTypeUI represents user interface errors
	ErrorTypeUI
	
	// ErrorTypeMCP represents MCP-specific errors
	ErrorTypeMCP
)

// String returns the string representation of the error type
func (e ErrorType) String() string {
	switch e {
	case ErrorTypeConfiguration:
		return "Configuration"
	case ErrorTypeConnection:
		return "Connection"
	case ErrorTypeValidation:
		return "Validation"
	case ErrorTypePermission:
		return "Permission"
	case ErrorTypeFileSystem:
		return "FileSystem"
	case ErrorTypeUI:
		return "UI"
	case ErrorTypeMCP:
		return "MCP"
	default:
		return "Unknown"
	}
}

// AppError represents an application error with additional context
type AppError struct {
	Type        ErrorType
	Code        string
	Message     string
	UserMessage string
	Cause       error
	Context     map[string]interface{}
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s:%s] %s: %v", e.Type, e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s:%s] %s", e.Type, e.Code, e.Message)
}

// Unwrap implements the error unwrapping interface
func (e *AppError) Unwrap() error {
	return e.Cause
}

// GetUserMessage returns a user-friendly error message
func (e *AppError) GetUserMessage() string {
	if e.UserMessage != "" {
		return e.UserMessage
	}
	return e.Message
}

// WithContext adds context information to the error
func (e *AppError) WithContext(key string, value interface{}) *AppError {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

// GetContext retrieves context information from the error
func (e *AppError) GetContext(key string) (interface{}, bool) {
	if e.Context == nil {
		return nil, false
	}
	value, exists := e.Context[key]
	return value, exists
}

// New creates a new application error
func New(errorType ErrorType, code, message string) *AppError {
	return &AppError{
		Type:    errorType,
		Code:    code,
		Message: message,
	}
}

// Wrap creates a new application error that wraps an existing error
func Wrap(err error, errorType ErrorType, code, message string) *AppError {
	return &AppError{
		Type:    errorType,
		Code:    code,
		Message: message,
		Cause:   err,
	}
}

// WithUserMessage adds a user-friendly message to an error
func WithUserMessage(err *AppError, userMessage string) *AppError {
	err.UserMessage = userMessage
	return err
}

// Configuration errors
var (
	ErrConfigNotFound = New(ErrorTypeConfiguration, "CONFIG_NOT_FOUND", 
		"Configuration file not found")
	ErrConfigInvalid = New(ErrorTypeConfiguration, "CONFIG_INVALID", 
		"Configuration file is invalid")
	ErrConfigValidation = New(ErrorTypeConfiguration, "CONFIG_VALIDATION", 
		"Configuration validation failed")
)

// Connection errors
var (
	ErrConnectionTimeout = New(ErrorTypeConnection, "CONNECTION_TIMEOUT", 
		"Connection timed out")
	ErrConnectionRefused = New(ErrorTypeConnection, "CONNECTION_REFUSED", 
		"Connection refused by server")
	ErrConnectionLost = New(ErrorTypeConnection, "CONNECTION_LOST", 
		"Connection lost")
)

// Validation errors
var (
	ErrInvalidInput = New(ErrorTypeValidation, "INVALID_INPUT", 
		"Invalid input provided")
	ErrMissingRequired = New(ErrorTypeValidation, "MISSING_REQUIRED", 
		"Required field is missing")
	ErrOutOfRange = New(ErrorTypeValidation, "OUT_OF_RANGE", 
		"Value is out of range")
)

// Permission errors
var (
	ErrAccessDenied = New(ErrorTypePermission, "ACCESS_DENIED", 
		"Access denied")
	ErrInsufficientPermissions = New(ErrorTypePermission, "INSUFFICIENT_PERMISSIONS", 
		"Insufficient permissions")
)

// File system errors
var (
	ErrFileNotFound = New(ErrorTypeFileSystem, "FILE_NOT_FOUND", 
		"File not found")
	ErrDirectoryNotFound = New(ErrorTypeFileSystem, "DIRECTORY_NOT_FOUND", 
		"Directory not found")
	ErrFilePermission = New(ErrorTypeFileSystem, "FILE_PERMISSION", 
		"File permission error")
)

// UI errors
var (
	ErrTerminalTooSmall = New(ErrorTypeUI, "TERMINAL_TOO_SMALL", 
		"Terminal window is too small")
	ErrTerminalNotSupported = New(ErrorTypeUI, "TERMINAL_NOT_SUPPORTED", 
		"Terminal not supported")
	ErrRenderFailed = New(ErrorTypeUI, "RENDER_FAILED", 
		"Failed to render UI")
)

// MCP errors
var (
	ErrMCPServerNotFound = New(ErrorTypeMCP, "MCP_SERVER_NOT_FOUND", 
		"MCP server not found")
	ErrMCPInvalidResponse = New(ErrorTypeMCP, "MCP_INVALID_RESPONSE", 
		"Invalid response from MCP server")
	ErrMCPMethodNotSupported = New(ErrorTypeMCP, "MCP_METHOD_NOT_SUPPORTED", 
		"MCP method not supported")
)

// UserFriendlyMessages maps error codes to user-friendly messages
var UserFriendlyMessages = map[string]string{
	"CONFIG_NOT_FOUND":         "No configuration file found. Using default settings.",
	"CONFIG_INVALID":           "Configuration file is corrupted. Please check the format.",
	"CONFIG_VALIDATION":        "Some configuration values are invalid. Please review your settings.",
	"CONNECTION_TIMEOUT":       "Unable to connect to the server. Please check your network connection.",
	"CONNECTION_REFUSED":       "Server is not responding. Please verify the server is running.",
	"CONNECTION_LOST":          "Connection to server was lost. Attempting to reconnect...",
	"INVALID_INPUT":            "Please check your input and try again.",
	"MISSING_REQUIRED":         "Some required information is missing. Please complete all fields.",
	"OUT_OF_RANGE":             "Value is outside the allowed range. Please adjust and try again.",
	"ACCESS_DENIED":            "You don't have permission to perform this action.",
	"INSUFFICIENT_PERMISSIONS": "Additional permissions are required to complete this operation.",
	"FILE_NOT_FOUND":           "The requested file could not be found.",
	"DIRECTORY_NOT_FOUND":      "The specified directory does not exist.",
	"FILE_PERMISSION":          "Unable to access the file due to permission restrictions.",
	"TERMINAL_TOO_SMALL":       "Terminal window is too small. Please resize to at least 40 columns.",
	"TERMINAL_NOT_SUPPORTED":   "Your terminal may not support all features of this application.",
	"RENDER_FAILED":            "Display error occurred. Please try resizing your terminal.",
	"MCP_SERVER_NOT_FOUND":     "The requested MCP server could not be found.",
	"MCP_INVALID_RESPONSE":     "Received unexpected response from MCP server.",
	"MCP_METHOD_NOT_SUPPORTED": "The requested operation is not supported by the MCP server.",
}

// GetUserFriendlyMessage returns a user-friendly message for an error code
func GetUserFriendlyMessage(code string) string {
	if msg, exists := UserFriendlyMessages[code]; exists {
		return msg
	}
	return "An unexpected error occurred. Please try again."
}

// FormatErrorForUser formats an error for display to the end user
func FormatErrorForUser(err error) string {
	if appErr, ok := err.(*AppError); ok {
		userMsg := appErr.GetUserMessage()
		if userMsg == "" {
			userMsg = GetUserFriendlyMessage(appErr.Code)
		}
		return userMsg
	}
	
	// Handle common Go errors
	errStr := err.Error()
	switch {
	case strings.Contains(errStr, "no such file or directory"):
		return "The requested file could not be found."
	case strings.Contains(errStr, "permission denied"):
		return "Permission denied. Please check file permissions."
	case strings.Contains(errStr, "connection refused"):
		return "Unable to connect to the server. Please check if it's running."
	case strings.Contains(errStr, "timeout"):
		return "Operation timed out. Please try again."
	case strings.Contains(errStr, "network"):
		return "Network error occurred. Please check your connection."
	default:
		return "An unexpected error occurred. Please try again."
	}
}

// Recovery provides error recovery suggestions
type Recovery struct {
	Action      string
	Description string
	AutoRetry   bool
}

// GetRecoveryActions returns possible recovery actions for an error
func GetRecoveryActions(err error) []Recovery {
	if appErr, ok := err.(*AppError); ok {
		switch appErr.Code {
		case "CONFIG_NOT_FOUND":
			return []Recovery{
				{
					Action:      "create_default_config",
					Description: "Create a default configuration file",
					AutoRetry:   true,
				},
			}
		case "CONFIG_INVALID":
			return []Recovery{
				{
					Action:      "reset_config",
					Description: "Reset to default configuration",
					AutoRetry:   false,
				},
				{
					Action:      "edit_config",
					Description: "Edit configuration file manually",
					AutoRetry:   false,
				},
			}
		case "CONNECTION_TIMEOUT", "CONNECTION_REFUSED":
			return []Recovery{
				{
					Action:      "retry_connection",
					Description: "Retry connection",
					AutoRetry:   true,
				},
				{
					Action:      "check_server_status",
					Description: "Check if server is running",
					AutoRetry:   false,
				},
			}
		case "TERMINAL_TOO_SMALL":
			return []Recovery{
				{
					Action:      "resize_terminal",
					Description: "Resize terminal window",
					AutoRetry:   true,
				},
			}
		}
	}
	
	return []Recovery{
		{
			Action:      "retry",
			Description: "Try again",
			AutoRetry:   false,
		},
	}
}

// IsRecoverable determines if an error can be automatically recovered from
func IsRecoverable(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		switch appErr.Code {
		case "CONFIG_NOT_FOUND", "CONNECTION_TIMEOUT", "CONNECTION_LOST", "TERMINAL_TOO_SMALL":
			return true
		}
	}
	return false
}

// IsCritical determines if an error is critical and should stop the application
func IsCritical(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		switch appErr.Type {
		case ErrorTypePermission:
			return true
		case ErrorTypeConfiguration:
			return appErr.Code == "CONFIG_VALIDATION"
		case ErrorTypeUI:
			return appErr.Code == "TERMINAL_NOT_SUPPORTED"
		}
	}
	return false
}