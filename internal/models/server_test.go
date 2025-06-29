package models

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestServerStatus_String(t *testing.T) {
	testCases := []struct {
		status   ServerStatus
		expected string
	}{
		{StatusStopped, "stopped"},
		{StatusStarting, "starting"},
		{StatusRunning, "running"},
		{StatusError, "error"},
		{StatusStopping, "stopping"},
		{StatusUnknown, "unknown"},
	}
	
	for _, tc := range testCases {
		result := tc.status.String()
		if result != tc.expected {
			t.Errorf("Status %d: expected '%s', got '%s'", tc.status, tc.expected, result)
		}
	}
}

func TestServerStatus_Icon(t *testing.T) {
	testCases := []struct {
		status ServerStatus
		icon   string
	}{
		{StatusStopped, "⭕"},
		{StatusStarting, "🔄"},
		{StatusRunning, "✅"},
		{StatusError, "❌"},
		{StatusStopping, "⏹️"},
		{StatusUnknown, "❓"},
	}
	
	for _, tc := range testCases {
		result := tc.status.Icon()
		if result != tc.icon {
			t.Errorf("Status %d: expected icon '%s', got '%s'", tc.status, tc.icon, result)
		}
	}
}

func TestServerType_String(t *testing.T) {
	testCases := []struct {
		serverType ServerType
		expected   string
	}{
		{TypeGeneric, "generic"},
		{TypeFileSystem, "filesystem"},
		{TypeDatabase, "database"},
		{TypeAPI, "api"},
		{TypeTool, "tool"},
		{TypeCustom, "custom"},
	}
	
	for _, tc := range testCases {
		result := tc.serverType.String()
		if result != tc.expected {
			t.Errorf("Type %d: expected '%s', got '%s'", tc.serverType, tc.expected, result)
		}
	}
}

func TestServerType_Icon(t *testing.T) {
	testCases := []struct {
		serverType ServerType
		icon       string
	}{
		{TypeGeneric, "📦"},
		{TypeFileSystem, "📁"},
		{TypeDatabase, "🗄️"},
		{TypeAPI, "🌐"},
		{TypeTool, "🔧"},
		{TypeCustom, "⚙️"},
	}
	
	for _, tc := range testCases {
		result := tc.serverType.Icon()
		if result != tc.icon {
			t.Errorf("Type %d: expected icon '%s', got '%s'", tc.serverType, tc.icon, result)
		}
	}
}

func TestServerConfig_Validate(t *testing.T) {
	testCases := []struct {
		name          string
		config        ServerConfig
		expectErrors  bool
		errorFields   []string
	}{
		{
			name: "valid minimal config",
			config: ServerConfig{
				Name:    "Test Server",
				Command: "/usr/bin/test-server",
			},
			expectErrors: false,
		},
		{
			name: "valid address-based config",
			config: ServerConfig{
				Name:    "Network Server",
				Address: "localhost:8080",
			},
			expectErrors: false,
		},
		{
			name: "missing name",
			config: ServerConfig{
				Command: "/usr/bin/test-server",
			},
			expectErrors: true,
			errorFields:  []string{"name"},
		},
		{
			name: "missing command and address",
			config: ServerConfig{
				Name: "Test Server",
			},
			expectErrors: true,
			errorFields:  []string{"command/address"},
		},
		{
			name: "name too long",
			config: ServerConfig{
				Name:    strings.Repeat("a", 101),
				Command: "/usr/bin/test-server",
			},
			expectErrors: true,
			errorFields:  []string{"name"},
		},
		{
			name: "description too long",
			config: ServerConfig{
				Name:        "Test Server",
				Command:     "/usr/bin/test-server",
				Description: strings.Repeat("a", 501),
			},
			expectErrors: true,
			errorFields:  []string{"description"},
		},
		{
			name: "invalid address format",
			config: ServerConfig{
				Name:    "Test Server",
				Address: "invalid-address",
			},
			expectErrors: true,
			errorFields:  []string{"address"},
		},
		{
			name: "invalid timeout",
			config: ServerConfig{
				Name:    "Test Server",
				Command: "/usr/bin/test-server",
				Timeout: -1,
			},
			expectErrors: true,
			errorFields:  []string{"timeout"},
		},
		{
			name: "invalid max restarts",
			config: ServerConfig{
				Name:        "Test Server",
				Command:     "/usr/bin/test-server",
				MaxRestarts: -1,
			},
			expectErrors: true,
			errorFields:  []string{"max_restarts"},
		},
		{
			name: "invalid homepage URL",
			config: ServerConfig{
				Name:     "Test Server",
				Command:  "/usr/bin/test-server",
				Homepage: "ftp://not-http-url.com",
			},
			expectErrors: true,
			errorFields:  []string{"homepage"},
		},
		{
			name: "TLS enabled without cert file",
			config: ServerConfig{
				Name:       "Test Server",
				Command:    "/usr/bin/test-server",
				TLSEnabled: true,
				TLSKeyFile: "/path/to/key.pem",
			},
			expectErrors: true,
			errorFields:  []string{"tls_cert_file"},
		},
		{
			name: "TLS enabled without key file",
			config: ServerConfig{
				Name:        "Test Server",
				Command:     "/usr/bin/test-server",
				TLSEnabled:  true,
				TLSCertFile: "/path/to/cert.pem",
			},
			expectErrors: true,
			errorFields:  []string{"tls_key_file"},
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			errors := tc.config.Validate()
			
			if tc.expectErrors {
				if len(errors) == 0 {
					t.Error("Expected validation errors but got none")
					return
				}
				
				// Check that expected error fields are present
				errorFieldMap := make(map[string]bool)
				for _, err := range errors {
					errorFieldMap[err.Field] = true
				}
				
				for _, expectedField := range tc.errorFields {
					if !errorFieldMap[expectedField] {
						t.Errorf("Expected error for field '%s' but didn't find it", expectedField)
					}
				}
			} else {
				if len(errors) > 0 {
					t.Errorf("Expected no validation errors but got: %v", errors)
				}
			}
		})
	}
}

func TestServerConfig_SetDefaults(t *testing.T) {
	config := ServerConfig{
		Name:       "Test Server",
		Command:    "/usr/bin/test-server",
		RestartOn:  true,
		HealthCheck: true,
	}
	
	config.SetDefaults()
	
	// Check that defaults were set
	if config.Type != TypeGeneric {
		t.Errorf("Expected default type to be TypeGeneric, got %v", config.Type)
	}
	
	if config.Timeout != 30 {
		t.Errorf("Expected default timeout to be 30, got %d", config.Timeout)
	}
	
	if config.MaxRestarts != 3 {
		t.Errorf("Expected default max_restarts to be 3, got %d", config.MaxRestarts)
	}
	
	if config.HealthInterval != 30*time.Second {
		t.Errorf("Expected default health_interval to be 30s, got %v", config.HealthInterval)
	}
	
	if config.StartupTimeout != 30*time.Second {
		t.Errorf("Expected default startup_timeout to be 30s, got %v", config.StartupTimeout)
	}
	
	if config.ShutdownTimeout != 10*time.Second {
		t.Errorf("Expected default shutdown_timeout to be 10s, got %v", config.ShutdownTimeout)
	}
	
	if config.Environment == nil {
		t.Error("Expected environment map to be initialized")
	}
}

func TestServerConfig_Clone(t *testing.T) {
	original := ServerConfig{
		Name:        "Test Server",
		Command:     "/usr/bin/test-server",
		Args:        []string{"--arg1", "--arg2"},
		Tags:        []string{"tag1", "tag2"},
		Environment: map[string]string{"VAR1": "value1", "VAR2": "value2"},
	}
	
	clone := original.Clone()
	
	// Check that all fields are copied
	if clone.Name != original.Name {
		t.Errorf("Name not cloned correctly: expected '%s', got '%s'", original.Name, clone.Name)
	}
	
	if clone.Command != original.Command {
		t.Errorf("Command not cloned correctly: expected '%s', got '%s'", original.Command, clone.Command)
	}
	
	// Check that slices are deep copied
	if len(clone.Args) != len(original.Args) {
		t.Error("Args slice not cloned correctly")
	} else {
		clone.Args[0] = "modified"
		if original.Args[0] == "modified" {
			t.Error("Args slice not deep copied - original was modified")
		}
	}
	
	// Check that maps are deep copied
	if len(clone.Environment) != len(original.Environment) {
		t.Error("Environment map not cloned correctly")
	} else {
		clone.Environment["VAR1"] = "modified"
		if original.Environment["VAR1"] == "modified" {
			t.Error("Environment map not deep copied - original was modified")
		}
	}
}

func TestNewServer(t *testing.T) {
	config := ServerConfig{
		Name:    "Test Server",
		Command: "/usr/bin/test-server",
		Type:    TypeTool,
	}
	
	server, err := NewServer(config)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}
	
	// Check that server was created with correct values
	if server.Config.Name != config.Name {
		t.Errorf("Expected server name '%s', got '%s'", config.Name, server.Config.Name)
	}
	
	if server.Status != StatusStopped {
		t.Errorf("Expected initial status to be StatusStopped, got %v", server.Status)
	}
	
	if server.ID == "" {
		t.Error("Expected server ID to be generated")
	}
	
	if server.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}
	
	if server.UpdatedAt.IsZero() {
		t.Error("Expected UpdatedAt to be set")
	}
}

func TestNewServer_ValidationError(t *testing.T) {
	config := ServerConfig{
		// Missing required fields
	}
	
	_, err := NewServer(config)
	if err == nil {
		t.Error("Expected validation error for invalid config")
	}
	
	if _, ok := err.(ValidationError); !ok {
		t.Errorf("Expected ValidationError, got %T", err)
	}
}

func TestServer_UpdateStatus(t *testing.T) {
	server := &Server{
		Status:       StatusStopped,
		StatusReason: "initial",
	}
	
	initialUpdateTime := server.UpdatedAt
	time.Sleep(1 * time.Millisecond) // Ensure time difference
	
	server.UpdateStatus(StatusRunning, "started successfully")
	
	if server.Status != StatusRunning {
		t.Errorf("Expected status to be StatusRunning, got %v", server.Status)
	}
	
	if server.StatusReason != "started successfully" {
		t.Errorf("Expected status reason 'started successfully', got '%s'", server.StatusReason)
	}
	
	if !server.UpdatedAt.After(initialUpdateTime) {
		t.Error("Expected UpdatedAt to be updated")
	}
	
	if server.LastSeenAt.IsZero() {
		t.Error("Expected LastSeenAt to be set for running status")
	}
}

func TestServer_RecordError(t *testing.T) {
	server := &Server{
		Status: StatusRunning,
		Metrics: ServerMetrics{
			ErrorCount: 0,
		},
	}
	
	testError := ValidationError{Field: "test", Message: "test error"}
	server.RecordError(testError)
	
	if server.LastError != testError.Error() {
		t.Errorf("Expected last error '%s', got '%s'", testError.Error(), server.LastError)
	}
	
	if server.Metrics.ErrorCount != 1 {
		t.Errorf("Expected error count to be 1, got %d", server.Metrics.ErrorCount)
	}
	
	if server.Status != StatusError {
		t.Errorf("Expected status to change to StatusError, got %v", server.Status)
	}
	
	if server.LastErrorTime.IsZero() {
		t.Error("Expected LastErrorTime to be set")
	}
}

func TestServer_IsHealthy(t *testing.T) {
	testCases := []struct {
		name     string
		server   Server
		expected bool
	}{
		{
			name: "stopped server",
			server: Server{
				Status: StatusStopped,
			},
			expected: false,
		},
		{
			name: "running server without health check",
			server: Server{
				Status: StatusRunning,
				Config: ServerConfig{
					HealthCheck: false,
				},
			},
			expected: true,
		},
		{
			name: "running server with recent heartbeat",
			server: Server{
				Status: StatusRunning,
				Config: ServerConfig{
					HealthCheck:    true,
					HealthInterval: 30 * time.Second,
				},
				Metrics: ServerMetrics{
					LastHeartbeat: time.Now().Add(-10 * time.Second),
				},
			},
			expected: true,
		},
		{
			name: "running server with old heartbeat",
			server: Server{
				Status: StatusRunning,
				Config: ServerConfig{
					HealthCheck:    true,
					HealthInterval: 30 * time.Second,
				},
				Metrics: ServerMetrics{
					LastHeartbeat: time.Now().Add(-2 * time.Minute),
				},
			},
			expected: false,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.server.IsHealthy()
			if result != tc.expected {
				t.Errorf("Expected IsHealthy() to return %t, got %t", tc.expected, result)
			}
		})
	}
}

func TestServer_GetDisplayName(t *testing.T) {
	testCases := []struct {
		name     string
		server   Server
		expected string
	}{
		{
			name: "server with name",
			server: Server{
				Config: ServerConfig{
					Name: "My Test Server",
				},
			},
			expected: "My Test Server",
		},
		{
			name: "server with command but no name",
			server: Server{
				Config: ServerConfig{
					Command: "/usr/bin/test-server",
				},
			},
			expected: "Command: /usr/bin/test-server",
		},
		{
			name: "server with address but no name or command",
			server: Server{
				Config: ServerConfig{
					Address: "localhost:8080",
				},
			},
			expected: "Address: localhost:8080",
		},
		{
			name: "server with only ID",
			server: Server{
				ID: "test-server-123",
			},
			expected: "test-server-123",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.server.GetDisplayName()
			if result != tc.expected {
				t.Errorf("Expected display name '%s', got '%s'", tc.expected, result)
			}
		})
	}
}

func TestServer_JSONSerialization(t *testing.T) {
	server := Server{
		Config: ServerConfig{
			Name: "Test Server",
			Type: TypeTool,
		},
		Status: StatusRunning,
		ID:     "test-123",
	}
	
	// Marshal to JSON
	data, err := json.Marshal(server)
	if err != nil {
		t.Fatalf("Failed to marshal server to JSON: %v", err)
	}
	
	// Check that additional string fields are included
	var jsonData map[string]interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}
	
	if jsonData["status_string"] != "running" {
		t.Errorf("Expected status_string to be 'running', got %v", jsonData["status_string"])
	}
	
	if jsonData["type_string"] != "tool" {
		t.Errorf("Expected type_string to be 'tool', got %v", jsonData["type_string"])
	}
}

func TestValidationError_Error(t *testing.T) {
	err := ValidationError{
		Field:   "test_field",
		Value:   "invalid_value",
		Message: "test message",
	}
	
	expected := "validation error for field 'test_field': test message"
	if err.Error() != expected {
		t.Errorf("Expected error message '%s', got '%s'", expected, err.Error())
	}
}

func BenchmarkServerConfig_Validate(b *testing.B) {
	config := ServerConfig{
		Name:        "Test Server",
		Description: "A test server for benchmarking",
		Command:     "/usr/bin/test-server",
		Args:        []string{"--arg1", "--arg2"},
		Type:        TypeTool,
		AutoStart:   true,
		Tags:        []string{"test", "benchmark"},
		Environment: map[string]string{"VAR1": "value1"},
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		config.Validate()
	}
}

func BenchmarkNewServer(b *testing.B) {
	config := ServerConfig{
		Name:    "Test Server",
		Command: "/usr/bin/test-server",
		Type:    TypeTool,
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := NewServer(config)
		if err != nil {
			b.Fatalf("Failed to create server: %v", err)
		}
	}
}