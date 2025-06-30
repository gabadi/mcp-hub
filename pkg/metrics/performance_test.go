package metrics

import (
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	timer := NewTimer()
	
	// Sleep for a small amount to ensure measurable time
	time.Sleep(10 * time.Millisecond)
	
	duration := timer.Stop()
	if duration < 10 {
		t.Errorf("Expected duration >= 10ms, got %dms", duration)
	}
}

func TestRecordLoadOperation(t *testing.T) {
	// Reset metrics before test
	ResetMetrics()
	
	RecordLoadOperation(100, 5)
	
	metrics := GetSummary()
	
	if metrics["load_operations_total"] != 1 {
		t.Errorf("Expected load_operations_total = 1, got %d", metrics["load_operations_total"])
	}
	
	if metrics["load_time_total_ms"] != 100 {
		t.Errorf("Expected load_time_total_ms = 100, got %d", metrics["load_time_total_ms"])
	}
	
	if metrics["current_inventory_size"] != 5 {
		t.Errorf("Expected current_inventory_size = 5, got %d", metrics["current_inventory_size"])
	}
}

func TestRecordSaveOperation(t *testing.T) {
	// Reset metrics before test
	ResetMetrics()
	
	RecordSaveOperation(150, 10)
	
	metrics := GetSummary()
	
	if metrics["save_operations_total"] != 1 {
		t.Errorf("Expected save_operations_total = 1, got %d", metrics["save_operations_total"])
	}
	
	if metrics["save_time_total_ms"] != 150 {
		t.Errorf("Expected save_time_total_ms = 150, got %d", metrics["save_time_total_ms"])
	}
	
	if metrics["current_inventory_size"] != 10 {
		t.Errorf("Expected current_inventory_size = 10, got %d", metrics["current_inventory_size"])
	}
}

func TestRecordValidation(t *testing.T) {
	// Reset metrics before test
	ResetMetrics()
	
	RecordValidation(50)
	
	metrics := GetSummary()
	
	if metrics["validation_time_total_ms"] != 50 {
		t.Errorf("Expected validation_time_total_ms = 50, got %d", metrics["validation_time_total_ms"])
	}
}

func TestRecordRecoveryOperation(t *testing.T) {
	// Reset metrics before test
	ResetMetrics()
	
	RecordRecoveryOperation()
	
	metrics := GetSummary()
	
	if metrics["recovery_operations_total"] != 1 {
		t.Errorf("Expected recovery_operations_total = 1, got %d", metrics["recovery_operations_total"])
	}
}

func TestRecordBackupOperation(t *testing.T) {
	// Reset metrics before test
	ResetMetrics()
	
	RecordBackupOperation()
	
	metrics := GetSummary()
	
	if metrics["backup_operations_total"] != 1 {
		t.Errorf("Expected backup_operations_total = 1, got %d", metrics["backup_operations_total"])
	}
}

func TestMultipleOperations(t *testing.T) {
	// Reset metrics before test
	ResetMetrics()
	
	// Record multiple operations
	RecordLoadOperation(100, 5)
	RecordLoadOperation(200, 10)
	RecordSaveOperation(150, 10)
	RecordValidation(25)
	RecordValidation(35)
	
	metrics := GetSummary()
	
	// Check totals
	if metrics["load_operations_total"] != 2 {
		t.Errorf("Expected load_operations_total = 2, got %d", metrics["load_operations_total"])
	}
	
	if metrics["load_time_total_ms"] != 300 {
		t.Errorf("Expected load_time_total_ms = 300, got %d", metrics["load_time_total_ms"])
	}
	
	if metrics["save_operations_total"] != 1 {
		t.Errorf("Expected save_operations_total = 1, got %d", metrics["save_operations_total"])
	}
	
	if metrics["validation_time_total_ms"] != 60 {
		t.Errorf("Expected validation_time_total_ms = 60, got %d", metrics["validation_time_total_ms"])
	}
	
	// Inventory size should be the last recorded value
	if metrics["current_inventory_size"] != 10 {
		t.Errorf("Expected current_inventory_size = 10, got %d", metrics["current_inventory_size"])
	}
}

func TestResetMetrics(t *testing.T) {
	// Set some metrics
	RecordLoadOperation(100, 5)
	RecordSaveOperation(150, 10)
	
	// Reset them
	ResetMetrics()
	
	metrics := GetSummary()
	
	// All metrics should be zero
	for key, value := range metrics {
		if value != 0 {
			t.Errorf("Expected %s = 0 after reset, got %d", key, value)
		}
	}
}