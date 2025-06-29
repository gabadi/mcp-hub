// Package models provides data models and types for the MCP Manager CLI.
//
// This package defines the core data structures used throughout the application
// for representing MCP servers, their configurations, capabilities, and runtime state.
// All models are designed to be:
//   - JSON serializable for configuration storage
//   - Validatable with comprehensive error reporting
//   - Immutable for safe concurrent access
//   - Well-documented with clear field purposes
package models

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"
)

// ServerStatus represents the current operational status of an MCP server
type ServerStatus int

const (
	// StatusUnknown indicates the server status has not been determined
	StatusUnknown ServerStatus = iota
	
	// StatusStopped indicates the server is not running
	StatusStopped
	
	// StatusStarting indicates the server is in the process of starting up
	StatusStarting
	
	// StatusRunning indicates the server is running and responsive
	StatusRunning
	
	// StatusError indicates the server encountered an error
	StatusError
	
	// StatusStopping indicates the server is in the process of shutting down
	StatusStopping
)

// String returns the string representation of the server status
func (s ServerStatus) String() string {
	switch s {
	case StatusStopped:
		return "stopped"
	case StatusStarting:
		return "starting"
	case StatusRunning:
		return "running"
	case StatusError:
		return "error"
	case StatusStopping:
		return "stopping"
	default:
		return "unknown"
	}
}

// Icon returns a visual icon for the server status
func (s ServerStatus) Icon() string {
	switch s {
	case StatusStopped:
		return "⭕"
	case StatusStarting:
		return "🔄"
	case StatusRunning:
		return "✅"
	case StatusError:
		return "❌"
	case StatusStopping:
		return "⏹️"
	default:
		return "❓"
	}
}

// ServerType represents the type/category of MCP server
type ServerType int

const (
	// TypeGeneric represents a generic MCP server
	TypeGeneric ServerType = iota
	
	// TypeFileSystem represents a file system MCP server
	TypeFileSystem
	
	// TypeDatabase represents a database MCP server
	TypeDatabase
	
	// TypeAPI represents an API integration MCP server
	TypeAPI
	
	// TypeTool represents a tool/utility MCP server
	TypeTool
	
	// TypeCustom represents a custom/user-defined MCP server
	TypeCustom
)

// String returns the string representation of the server type
func (t ServerType) String() string {
	switch t {
	case TypeFileSystem:
		return "filesystem"
	case TypeDatabase:
		return "database"
	case TypeAPI:
		return "api"
	case TypeTool:
		return "tool"
	case TypeCustom:
		return "custom"
	default:
		return "generic"
	}
}

// Icon returns a visual icon for the server type
func (t ServerType) Icon() string {
	switch t {
	case TypeFileSystem:
		return "📁"
	case TypeDatabase:
		return "🗄️"
	case TypeAPI:
		return "🌐"
	case TypeTool:
		return "🔧"
	case TypeCustom:
		return "⚙️"
	default:
		return "📦"
	}
}

// ServerConfig represents the configuration for an MCP server
type ServerConfig struct {
	// Basic identification
	Name        string `json:"name"`        // Human-readable server name
	Description string `json:"description"` // Server description
	Type        ServerType `json:"type"`    // Server type/category
	
	// Connection settings
	Command     string            `json:"command"`              // Command to execute the server
	Args        []string          `json:"args,omitempty"`       // Command line arguments
	WorkingDir  string            `json:"working_dir,omitempty"` // Working directory for execution
	Environment map[string]string `json:"environment,omitempty"` // Environment variables
	
	// Network settings (for networked MCP servers)
	Address string `json:"address,omitempty"` // Server address (host:port or Unix socket)
	Timeout int    `json:"timeout,omitempty"` // Connection timeout in seconds
	
	// Startup settings
	AutoStart   bool `json:"auto_start"`   // Whether to start automatically
	RestartOn   bool `json:"restart_on"`   // Whether to restart on failure
	MaxRestarts int  `json:"max_restarts"` // Maximum restart attempts
	
	// Metadata
	Tags     []string `json:"tags,omitempty"`     // User-defined tags
	Version  string   `json:"version,omitempty"`  // Server version
	Author   string   `json:"author,omitempty"`   // Server author/maintainer
	Homepage string   `json:"homepage,omitempty"` // Server homepage URL
	
	// Advanced settings
	HealthCheck     bool          `json:"health_check"`      // Enable health checking
	HealthInterval  time.Duration `json:"health_interval"`   // Health check interval
	StartupTimeout  time.Duration `json:"startup_timeout"`   // Maximum time to wait for startup
	ShutdownTimeout time.Duration `json:"shutdown_timeout"`  // Maximum time to wait for shutdown
	
	// Security settings
	TLSEnabled   bool   `json:"tls_enabled,omitempty"`   // Enable TLS
	TLSCertFile  string `json:"tls_cert_file,omitempty"` // TLS certificate file
	TLSKeyFile   string `json:"tls_key_file,omitempty"`  // TLS private key file
	TLSCAFile    string `json:"tls_ca_file,omitempty"`   // TLS CA certificate file
	TLSInsecure  bool   `json:"tls_insecure,omitempty"`  // Skip TLS verification
}

// ServerCapabilities represents the capabilities exposed by an MCP server
type ServerCapabilities struct {
	// Core capabilities
	Tools     bool `json:"tools"`     // Server provides tools
	Resources bool `json:"resources"` // Server provides resources
	Prompts   bool `json:"prompts"`   // Server provides prompts
	
	// Advanced capabilities
	Sampling       bool `json:"sampling"`        // Server supports sampling
	LoggingEnabled bool `json:"logging_enabled"` // Server supports logging
	
	// Experimental capabilities
	Experimental map[string]interface{} `json:"experimental,omitempty"` // Experimental features
}

// ServerMetrics represents runtime metrics for an MCP server
type ServerMetrics struct {
	// Runtime information
	StartTime     time.Time     `json:"start_time"`     // When the server started
	Uptime        time.Duration `json:"uptime"`         // How long the server has been running
	LastHeartbeat time.Time     `json:"last_heartbeat"` // Last successful health check
	
	// Performance metrics
	RequestCount    int64         `json:"request_count"`    // Total requests processed
	ErrorCount      int64         `json:"error_count"`      // Total errors encountered
	AverageLatency  time.Duration `json:"average_latency"`  // Average request latency
	MemoryUsage     int64         `json:"memory_usage"`     // Memory usage in bytes
	CPUUsage        float64       `json:"cpu_usage"`        // CPU usage percentage
	
	// Connection metrics
	ConnectionCount    int `json:"connection_count"`     // Active connections
	MaxConnections     int `json:"max_connections"`      // Maximum concurrent connections
	TotalConnections   int `json:"total_connections"`    // Total connections since start
	FailedConnections  int `json:"failed_connections"`   // Failed connection attempts
	
	// Restart information
	RestartCount     int       `json:"restart_count"`      // Number of restarts
	LastRestartTime  time.Time `json:"last_restart_time"`  // When last restarted
	LastRestartError string    `json:"last_restart_error"` // Error that caused last restart
}

// Server represents a complete MCP server with configuration, state, and metrics
type Server struct {
	// Configuration (immutable after creation)
	Config ServerConfig `json:"config"`
	
	// Runtime state
	Status       ServerStatus        `json:"status"`       // Current operational status
	StatusReason string              `json:"status_reason"` // Reason for current status
	Capabilities ServerCapabilities  `json:"capabilities"` // Server capabilities
	Metrics      ServerMetrics       `json:"metrics"`      // Runtime metrics
	
	// Metadata
	ID           string    `json:"id"`            // Unique server identifier
	CreatedAt    time.Time `json:"created_at"`    // When the server was added
	UpdatedAt    time.Time `json:"updated_at"`    // When the server was last modified
	LastSeenAt   time.Time `json:"last_seen_at"`  // When the server was last contacted
	
	// Error information
	LastError     string    `json:"last_error,omitempty"`      // Last error message
	LastErrorTime time.Time `json:"last_error_time,omitempty"` // When the last error occurred
}

// ValidationError represents a server configuration validation error
type ValidationError struct {
	Field   string      `json:"field"`   // Field that failed validation
	Value   interface{} `json:"value"`   // Invalid value
	Message string      `json:"message"` // Error message
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error for field '%s': %s", e.Field, e.Message)
}

// Validate validates the server configuration and returns any errors
func (c *ServerConfig) Validate() []ValidationError {
	var errors []ValidationError
	
	// Validate required fields
	if strings.TrimSpace(c.Name) == "" {
		errors = append(errors, ValidationError{
			Field:   "name",
			Value:   c.Name,
			Message: "server name is required",
		})
	}
	
	if strings.TrimSpace(c.Command) == "" && strings.TrimSpace(c.Address) == "" {
		errors = append(errors, ValidationError{
			Field:   "command/address",
			Value:   fmt.Sprintf("command='%s', address='%s'", c.Command, c.Address),
			Message: "either command or address must be specified",
		})
	}
	
	// Validate name format
	if len(c.Name) > 100 {
		errors = append(errors, ValidationError{
			Field:   "name",
			Value:   c.Name,
			Message: "server name must be 100 characters or less",
		})
	}
	
	// Validate description length
	if len(c.Description) > 500 {
		errors = append(errors, ValidationError{
			Field:   "description",
			Value:   c.Description,
			Message: "description must be 500 characters or less",
		})
	}
	
	// Validate network address format
	if c.Address != "" {
		if !strings.Contains(c.Address, ":") && !strings.HasPrefix(c.Address, "/") {
			errors = append(errors, ValidationError{
				Field:   "address",
				Value:   c.Address,
				Message: "address must be in format 'host:port' or Unix socket path",
			})
		}
	}
	
	// Validate timeout values
	if c.Timeout < 0 || c.Timeout > 300 {
		errors = append(errors, ValidationError{
			Field:   "timeout",
			Value:   c.Timeout,
			Message: "timeout must be between 0 and 300 seconds",
		})
	}
	
	// Validate restart settings
	if c.MaxRestarts < 0 || c.MaxRestarts > 100 {
		errors = append(errors, ValidationError{
			Field:   "max_restarts",
			Value:   c.MaxRestarts,
			Message: "max_restarts must be between 0 and 100",
		})
	}
	
	// Validate homepage URL format
	if c.Homepage != "" {
		if parsedURL, err := url.Parse(c.Homepage); err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
			errors = append(errors, ValidationError{
				Field:   "homepage",
				Value:   c.Homepage,
				Message: "homepage must be a valid HTTP or HTTPS URL",
			})
		}
	}
	
	// Validate duration values
	if c.HealthInterval < 0 || c.HealthInterval > 24*time.Hour {
		errors = append(errors, ValidationError{
			Field:   "health_interval",
			Value:   c.HealthInterval,
			Message: "health_interval must be between 0 and 24 hours",
		})
	}
	
	if c.StartupTimeout < 0 || c.StartupTimeout > 10*time.Minute {
		errors = append(errors, ValidationError{
			Field:   "startup_timeout",
			Value:   c.StartupTimeout,
			Message: "startup_timeout must be between 0 and 10 minutes",
		})
	}
	
	if c.ShutdownTimeout < 0 || c.ShutdownTimeout > 5*time.Minute {
		errors = append(errors, ValidationError{
			Field:   "shutdown_timeout",
			Value:   c.ShutdownTimeout,
			Message: "shutdown_timeout must be between 0 and 5 minutes",
		})
	}
	
	// Validate TLS configuration
	if c.TLSEnabled {
		if c.TLSCertFile == "" {
			errors = append(errors, ValidationError{
				Field:   "tls_cert_file",
				Value:   c.TLSCertFile,
				Message: "TLS certificate file is required when TLS is enabled",
			})
		}
		if c.TLSKeyFile == "" {
			errors = append(errors, ValidationError{
				Field:   "tls_key_file",
				Value:   c.TLSKeyFile,
				Message: "TLS private key file is required when TLS is enabled",
			})
		}
	}
	
	return errors
}

// IsValid returns true if the server configuration is valid
func (c *ServerConfig) IsValid() bool {
	return len(c.Validate()) == 0
}

// SetDefaults sets default values for optional configuration fields
func (c *ServerConfig) SetDefaults() {
	if c.Type == 0 {
		c.Type = TypeGeneric
	}
	
	if c.Timeout == 0 {
		c.Timeout = 30
	}
	
	if c.MaxRestarts == 0 && c.RestartOn {
		c.MaxRestarts = 3
	}
	
	if c.HealthInterval == 0 && c.HealthCheck {
		c.HealthInterval = 30 * time.Second
	}
	
	if c.StartupTimeout == 0 {
		c.StartupTimeout = 30 * time.Second
	}
	
	if c.ShutdownTimeout == 0 {
		c.ShutdownTimeout = 10 * time.Second
	}
	
	if c.Environment == nil {
		c.Environment = make(map[string]string)
	}
}

// Clone creates a deep copy of the server configuration
func (c *ServerConfig) Clone() *ServerConfig {
	clone := *c
	
	// Deep copy slices and maps
	if c.Args != nil {
		clone.Args = make([]string, len(c.Args))
		copy(clone.Args, c.Args)
	}
	
	if c.Tags != nil {
		clone.Tags = make([]string, len(c.Tags))
		copy(clone.Tags, c.Tags)
	}
	
	if c.Environment != nil {
		clone.Environment = make(map[string]string)
		for k, v := range c.Environment {
			clone.Environment[k] = v
		}
	}
	
	return &clone
}

// NewServer creates a new server with the given configuration
func NewServer(config ServerConfig) (*Server, error) {
	// Validate configuration
	config.SetDefaults()
	if validationErrors := config.Validate(); len(validationErrors) > 0 {
		return nil, validationErrors[0] // Return first validation error
	}
	
	now := time.Now()
	
	server := &Server{
		Config:       config,
		Status:       StatusStopped,
		StatusReason: "newly created",
		ID:           generateServerID(config.Name),
		CreatedAt:    now,
		UpdatedAt:    now,
		Capabilities: ServerCapabilities{}, // Will be determined by server communication
		Metrics:      ServerMetrics{},      // Will be populated during runtime
	}
	
	return server, nil
}

// generateServerID generates a unique identifier for a server
func generateServerID(name string) string {
	// For now, use a simple approach based on name
	// In a real implementation, this might use UUID or hash
	clean := strings.ToLower(strings.ReplaceAll(name, " ", "-"))
	timestamp := time.Now().Unix()
	return fmt.Sprintf("%s-%d", clean, timestamp)
}

// UpdateStatus updates the server status and status reason
func (s *Server) UpdateStatus(status ServerStatus, reason string) {
	s.Status = status
	s.StatusReason = reason
	s.UpdatedAt = time.Now()
	
	if status == StatusRunning {
		s.LastSeenAt = s.UpdatedAt
		s.Metrics.LastHeartbeat = s.UpdatedAt
	}
}

// RecordError records an error for the server
func (s *Server) RecordError(err error) {
	s.LastError = err.Error()
	s.LastErrorTime = time.Now()
	s.Metrics.ErrorCount++
	s.UpdatedAt = time.Now()
	
	if s.Status == StatusRunning {
		s.UpdateStatus(StatusError, "error occurred during operation")
	}
}

// IsHealthy returns true if the server is considered healthy
func (s *Server) IsHealthy() bool {
	if s.Status != StatusRunning {
		return false
	}
	
	// Check if health checking is enabled
	if !s.Config.HealthCheck {
		return true // Assume healthy if health checking is disabled
	}
	
	// Check if last heartbeat is within acceptable range
	maxAge := s.Config.HealthInterval * 2 // Allow 2x the interval
	if maxAge == 0 {
		maxAge = 60 * time.Second // Default to 60 seconds
	}
	
	return time.Since(s.Metrics.LastHeartbeat) < maxAge
}

// GetDisplayName returns a human-readable display name for the server
func (s *Server) GetDisplayName() string {
	if s.Config.Name != "" {
		return s.Config.Name
	}
	if s.Config.Command != "" {
		return fmt.Sprintf("Command: %s", s.Config.Command)
	}
	if s.Config.Address != "" {
		return fmt.Sprintf("Address: %s", s.Config.Address)
	}
	return s.ID
}

// String returns a string representation of the server
func (s *Server) String() string {
	return fmt.Sprintf("Server{name=%s, status=%s, type=%s}", 
		s.GetDisplayName(), s.Status, s.Config.Type)
}

// MarshalJSON customizes JSON marshaling for the server
func (s Server) MarshalJSON() ([]byte, error) {
	type Alias Server
	return json.Marshal(&struct {
		Alias
		StatusString string `json:"status_string"`
		TypeString   string `json:"type_string"`
	}{
		Alias:        Alias(s),
		StatusString: s.Status.String(),
		TypeString:   s.Config.Type.String(),
	})
}