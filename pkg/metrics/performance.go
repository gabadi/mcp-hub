package metrics

import (
	"expvar"
	"time"
)

// PerformanceMetrics holds performance tracking data
type PerformanceMetrics struct {
	LoadOperations   *expvar.Int
	SaveOperations   *expvar.Int
	LoadTime         *expvar.Int // in milliseconds
	SaveTime         *expvar.Int // in milliseconds
	InventorySize    *expvar.Int
	ValidationTime   *expvar.Int // in milliseconds
	RecoveryOperations *expvar.Int
	BackupOperations   *expvar.Int
}

var (
	// Global metrics instance
	globalMetrics *PerformanceMetrics
)

// init initializes the global metrics
func init() {
	globalMetrics = &PerformanceMetrics{
		LoadOperations:     expvar.NewInt("mcp_load_operations_total"),
		SaveOperations:     expvar.NewInt("mcp_save_operations_total"),
		LoadTime:          expvar.NewInt("mcp_load_time_ms"),
		SaveTime:          expvar.NewInt("mcp_save_time_ms"),
		InventorySize:     expvar.NewInt("mcp_inventory_size"),
		ValidationTime:    expvar.NewInt("mcp_validation_time_ms"),
		RecoveryOperations: expvar.NewInt("mcp_recovery_operations_total"),
		BackupOperations:   expvar.NewInt("mcp_backup_operations_total"),
	}

	// Register additional computed metrics
	expvar.Publish("mcp_avg_load_time_ms", expvar.Func(func() interface{} {
		loads := globalMetrics.LoadOperations.Value()
		if loads == 0 {
			return 0
		}
		return globalMetrics.LoadTime.Value() / loads
	}))

	expvar.Publish("mcp_avg_save_time_ms", expvar.Func(func() interface{} {
		saves := globalMetrics.SaveOperations.Value()
		if saves == 0 {
			return 0
		}
		return globalMetrics.SaveTime.Value() / saves
	}))

	expvar.Publish("mcp_avg_validation_time_ms", expvar.Func(func() interface{} {
		ops := globalMetrics.LoadOperations.Value() + globalMetrics.SaveOperations.Value()
		if ops == 0 {
			return 0
		}
		return globalMetrics.ValidationTime.Value() / ops
	}))
}

// GetMetrics returns the global metrics instance
func GetMetrics() *PerformanceMetrics {
	return globalMetrics
}

// Timer provides a simple way to measure operation duration
type Timer struct {
	start time.Time
}

// NewTimer creates a new timer
func NewTimer() *Timer {
	return &Timer{start: time.Now()}
}

// Stop returns the elapsed time in milliseconds
func (t *Timer) Stop() int64 {
	return time.Since(t.start).Milliseconds()
}

// RecordLoadOperation records a load operation with its duration
func RecordLoadOperation(durationMs int64, inventorySize int) {
	globalMetrics.LoadOperations.Add(1)
	globalMetrics.LoadTime.Add(durationMs)
	globalMetrics.InventorySize.Set(int64(inventorySize))
}

// RecordSaveOperation records a save operation with its duration
func RecordSaveOperation(durationMs int64, inventorySize int) {
	globalMetrics.SaveOperations.Add(1)
	globalMetrics.SaveTime.Add(durationMs)
	globalMetrics.InventorySize.Set(int64(inventorySize))
}

// RecordValidation records validation time
func RecordValidation(durationMs int64) {
	globalMetrics.ValidationTime.Add(durationMs)
}

// RecordRecoveryOperation records a recovery operation
func RecordRecoveryOperation() {
	globalMetrics.RecoveryOperations.Add(1)
}

// RecordBackupOperation records a backup operation
func RecordBackupOperation() {
	globalMetrics.BackupOperations.Add(1)
}

// GetSummary returns a summary of current metrics
func GetSummary() map[string]int64 {
	return map[string]int64{
		"load_operations_total":    globalMetrics.LoadOperations.Value(),
		"save_operations_total":    globalMetrics.SaveOperations.Value(),
		"load_time_total_ms":       globalMetrics.LoadTime.Value(),
		"save_time_total_ms":       globalMetrics.SaveTime.Value(),
		"current_inventory_size":   globalMetrics.InventorySize.Value(),
		"validation_time_total_ms": globalMetrics.ValidationTime.Value(),
		"recovery_operations_total": globalMetrics.RecoveryOperations.Value(),
		"backup_operations_total":   globalMetrics.BackupOperations.Value(),
	}
}

// ResetMetrics resets all metrics (mainly for testing)
func ResetMetrics() {
	globalMetrics.LoadOperations.Set(0)
	globalMetrics.SaveOperations.Set(0)
	globalMetrics.LoadTime.Set(0)
	globalMetrics.SaveTime.Set(0)
	globalMetrics.InventorySize.Set(0)
	globalMetrics.ValidationTime.Set(0)
	globalMetrics.RecoveryOperations.Set(0)
	globalMetrics.BackupOperations.Set(0)
}