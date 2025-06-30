package models

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"time"
)

// MCPType represents the different types of MCP configurations
type MCPType string

const (
	MCPTypeCommand MCPType = "command"
	MCPTypeSSE     MCPType = "sse"
	MCPTypeJSON    MCPType = "json"
	MCPTypeHTTP    MCPType = "http"
)

// MCP represents a Model Context Protocol configuration
type MCP struct {
	ID          string    `json:"id" validate:"required,min=1,max=100"`
	Name        string    `json:"name" validate:"required,min=1,max=255"`
	Type        MCPType   `json:"type" validate:"required,oneof=command sse json http"`
	Config      MCPConfig `json:"config" validate:"required"`
	CreatedAt   time.Time `json:"created_at" validate:"required"`
	UpdatedAt   time.Time `json:"updated_at" validate:"required"`
	Description string    `json:"description,omitempty" validate:"max=1000"`
	Enabled     bool      `json:"enabled"`
}

// MCPConfig holds the configuration for different MCP types
type MCPConfig struct {
	// Command/Binary MCP fields
	Command string   `json:"command,omitempty" validate:"max=500"`
	Args    []string `json:"args,omitempty" validate:"max=50"`
	
	// SSE Server MCP fields
	ServerURL string `json:"server_url,omitempty" validate:"url"`
	
	// JSON Configuration MCP fields
	JSONConfig map[string]interface{} `json:"json_config,omitempty"`
	
	// HTTP MCP fields
	Endpoint string            `json:"endpoint,omitempty" validate:"url"`
	Headers  map[string]string `json:"headers,omitempty" validate:"max=20"`
	Method   string            `json:"method,omitempty" validate:"oneof=GET POST PUT DELETE PATCH HEAD OPTIONS"`
}

// MCPInventory represents the complete inventory of MCP configurations
type MCPInventory struct {
	Version   string           `json:"version" validate:"required,min=1,max=10"`
	MCPs      []MCP            `json:"mcps" validate:"max=1000"`
	Metadata  InventoryMetadata `json:"metadata" validate:"required"`
	UpdatedAt time.Time        `json:"updated_at" validate:"required"`
}

// InventoryMetadata holds metadata about the inventory
type InventoryMetadata struct {
	Created   time.Time `json:"created" validate:"required"`
	FileCount int       `json:"file_count" validate:"min=0,max=1000"`
	LastSync  time.Time `json:"last_sync,omitempty"`
}

// Validate checks if the MCP configuration is valid
func (m *MCP) Validate() error {
	// First apply struct tag validation
	if err := ValidateStruct(m); err != nil {
		return err
	}
	
	// Then apply business logic validation
	switch m.Type {
	case MCPTypeCommand:
		if m.Config.Command == "" {
			return fmt.Errorf("command MCP must have a command specified")
		}
	case MCPTypeSSE:
		if m.Config.ServerURL == "" {
			return fmt.Errorf("SSE MCP must have a server URL specified")
		}
	case MCPTypeJSON:
		if m.Config.JSONConfig == nil || len(m.Config.JSONConfig) == 0 {
			return fmt.Errorf("JSON MCP must have configuration specified")
		}
	case MCPTypeHTTP:
		if m.Config.Endpoint == "" {
			return fmt.Errorf("HTTP MCP must have an endpoint specified")
		}
		if m.Config.Method == "" {
			m.Config.Method = "GET" // default to GET
		}
	default:
		return fmt.Errorf("unsupported MCP type: %s", m.Type)
	}
	
	return nil
}

// NewMCP creates a new MCP with basic initialization
func NewMCP(id, name string, mcpType MCPType) *MCP {
	now := time.Now()
	return &MCP{
		ID:        id,
		Name:      name,
		Type:      mcpType,
		Config:    MCPConfig{},
		CreatedAt: now,
		UpdatedAt: now,
		Enabled:   true,
	}
}

// NewInventory creates a new empty inventory
func NewInventory() *MCPInventory {
	now := time.Now()
	return &MCPInventory{
		Version:   "1.0",
		MCPs:      []MCP{},
		Metadata: InventoryMetadata{
			Created:   now,
			FileCount: 0,
		},
		UpdatedAt: now,
	}
}

// AddMCP adds an MCP to the inventory
func (inv *MCPInventory) AddMCP(mcp MCP) error {
	if err := mcp.Validate(); err != nil {
		return fmt.Errorf("invalid MCP: %w", err)
	}
	
	// Check for duplicate IDs
	for _, existingMCP := range inv.MCPs {
		if existingMCP.ID == mcp.ID {
			return fmt.Errorf("MCP with ID %s already exists", mcp.ID)
		}
	}
	
	inv.MCPs = append(inv.MCPs, mcp)
	inv.Metadata.FileCount = len(inv.MCPs)
	inv.UpdatedAt = time.Now()
	
	return nil
}

// RemoveMCP removes an MCP from the inventory by ID
func (inv *MCPInventory) RemoveMCP(id string) error {
	for i, mcp := range inv.MCPs {
		if mcp.ID == id {
			inv.MCPs = append(inv.MCPs[:i], inv.MCPs[i+1:]...)
			inv.Metadata.FileCount = len(inv.MCPs)
			inv.UpdatedAt = time.Now()
			return nil
		}
	}
	return fmt.Errorf("MCP with ID %s not found", id)
}

// GetMCP retrieves an MCP by ID
func (inv *MCPInventory) GetMCP(id string) (*MCP, error) {
	for i, mcp := range inv.MCPs {
		if mcp.ID == id {
			return &inv.MCPs[i], nil
		}
	}
	return nil, fmt.Errorf("MCP with ID %s not found", id)
}

// UpdateMCP updates an existing MCP in the inventory
func (inv *MCPInventory) UpdateMCP(updatedMCP MCP) error {
	if err := updatedMCP.Validate(); err != nil {
		return fmt.Errorf("invalid MCP: %w", err)
	}
	
	for i, mcp := range inv.MCPs {
		if mcp.ID == updatedMCP.ID {
			updatedMCP.UpdatedAt = time.Now()
			inv.MCPs[i] = updatedMCP
			inv.UpdatedAt = time.Now()
			return nil
		}
	}
	
	return fmt.Errorf("MCP with ID %s not found", updatedMCP.ID)
}

// Validate checks if the inventory is valid
func (inv *MCPInventory) Validate() error {
	// First apply struct tag validation
	if err := ValidateStruct(inv); err != nil {
		return err
	}
	
	// Then apply business logic validation
	// Check for duplicate IDs
	ids := make(map[string]bool)
	for _, mcp := range inv.MCPs {
		if ids[mcp.ID] {
			return fmt.Errorf("duplicate MCP ID found: %s", mcp.ID)
		}
		ids[mcp.ID] = true
		
		if err := mcp.Validate(); err != nil {
			return fmt.Errorf("invalid MCP %s: %w", mcp.ID, err)
		}
	}
	
	return nil
}

// ToJSON converts the inventory to JSON format
func (inv *MCPInventory) ToJSON() ([]byte, error) {
	return json.MarshalIndent(inv, "", "  ")
}

// FromJSON creates an inventory from JSON data
func FromJSON(data []byte) (*MCPInventory, error) {
	var inventory MCPInventory
	if err := json.Unmarshal(data, &inventory); err != nil {
		return nil, fmt.Errorf("failed to unmarshal inventory: %w", err)
	}
	
	if err := inventory.Validate(); err != nil {
		return nil, fmt.Errorf("invalid inventory: %w", err)
	}
	
	return &inventory, nil
}

// ValidateStruct validates a struct based on its validation tags
func ValidateStruct(s interface{}) error {
	v := reflect.ValueOf(s)
	t := reflect.TypeOf(s)
	
	// Handle pointer types
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return fmt.Errorf("cannot validate nil pointer")
		}
		v = v.Elem()
		t = t.Elem()
	}
	
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("validation only supports struct types")
	}
	
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		
		// Skip unexported fields
		if !field.CanInterface() {
			continue
		}
		
		// Get validation tag
		tag := fieldType.Tag.Get("validate")
		if tag == "" {
			continue
		}
		
		// Skip validation for zero time values if not required
		if field.Type() == reflect.TypeOf(time.Time{}) && field.Interface().(time.Time).IsZero() {
			if !strings.Contains(tag, "required") {
				continue
			}
		}
		
		if err := validateField(field, fieldType.Name, tag); err != nil {
			return err
		}
	}
	
	return nil
}

// validateField validates a single field based on its validation tag
func validateField(field reflect.Value, fieldName, tag string) error {
	rules := strings.Split(tag, ",")
	
	for _, rule := range rules {
		rule = strings.TrimSpace(rule)
		if rule == "" {
			continue
		}
		
		if err := applyValidationRule(field, fieldName, rule); err != nil {
			return err
		}
	}
	
	return nil
}

// applyValidationRule applies a single validation rule to a field
func applyValidationRule(field reflect.Value, fieldName, rule string) error {
	parts := strings.SplitN(rule, "=", 2)
	ruleName := parts[0]
	var ruleValue string
	if len(parts) > 1 {
		ruleValue = parts[1]
	}
	
	switch ruleName {
	case "required":
		if field.Kind() == reflect.String && field.String() == "" {
			return fmt.Errorf("field %s is required", fieldName)
		}
		if field.Kind() == reflect.Slice && field.Len() == 0 {
			return fmt.Errorf("field %s is required", fieldName)
		}
		if field.Type() == reflect.TypeOf(time.Time{}) && field.Interface().(time.Time).IsZero() {
			return fmt.Errorf("field %s is required", fieldName)
		}
		
	case "min":
		if ruleValue == "" {
			return fmt.Errorf("min rule requires a value")
		}
		min := parseIntOrDefault(ruleValue, 0)
		if field.Kind() == reflect.String && len(field.String()) < min {
			return fmt.Errorf("field %s must be at least %d characters", fieldName, min)
		}
		if field.Kind() == reflect.Int && int(field.Int()) < min {
			return fmt.Errorf("field %s must be at least %d", fieldName, min)
		}
		
	case "max":
		if ruleValue == "" {
			return fmt.Errorf("max rule requires a value")
		}
		max := parseIntOrDefault(ruleValue, 0)
		if field.Kind() == reflect.String && len(field.String()) > max {
			return fmt.Errorf("field %s must be at most %d characters", fieldName, max)
		}
		if field.Kind() == reflect.Int && int(field.Int()) > max {
			return fmt.Errorf("field %s must be at most %d", fieldName, max)
		}
		if field.Kind() == reflect.Slice && field.Len() > max {
			return fmt.Errorf("field %s must have at most %d items", fieldName, max)
		}
		if field.Kind() == reflect.Map && field.Len() > max {
			return fmt.Errorf("field %s must have at most %d items", fieldName, max)
		}
		
	case "oneof":
		if ruleValue == "" {
			return fmt.Errorf("oneof rule requires values")
		}
		validValues := strings.Split(ruleValue, " ")
		fieldValue := field.String()
		for _, validValue := range validValues {
			if fieldValue == validValue {
				return nil
			}
		}
		return fmt.Errorf("field %s must be one of: %s", fieldName, strings.Join(validValues, ", "))
		
	case "url":
		if field.Kind() != reflect.String {
			return fmt.Errorf("url validation only applies to string fields")
		}
		urlStr := field.String()
		if urlStr == "" {
			return nil // Empty URLs are valid if not required
		}
		if _, err := url.ParseRequestURI(urlStr); err != nil {
			return fmt.Errorf("field %s must be a valid URL", fieldName)
		}
	}
	
	return nil
}

// parseIntOrDefault parses an integer string or returns default value
func parseIntOrDefault(s string, defaultValue int) int {
	if s == "" {
		return defaultValue
	}
	
	var result int
	if _, err := fmt.Sscanf(s, "%d", &result); err != nil {
		return defaultValue
	}
	return result
}