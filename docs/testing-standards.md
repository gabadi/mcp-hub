# Testing Standards - Definition of Done

**Document Version:** 1.0 | **Created:** 2025-06-30 | **Last Updated:** 2025-06-30 | **Owner:** Development Team | **Approved By:** Technical Lead & QA Team

## Overview

This document defines the Testing Definition of Done (DoD) standards for the CC MCP Manager project, ensuring production-ready quality assurance and full process compliance for all code changes.

## Testing Definition of Done Criteria

### Mandatory Requirements

Every code change MUST satisfy ALL criteria before completion:

**Unit Test Coverage:** 85% line coverage for new code, verify happy path and error conditions, independent tests without external dependencies, follow `Test[MethodName]_[Scenario]` naming, clear assertions with descriptive error messages

**Integration Test Coverage:** End-to-end workflows, component integration verification, error scenarios and recovery, data persistence integrity validation

**Test Organization:** Co-located test files (`*_test.go`), categorized tests (unit, integration, benchmark), proper setup/teardown with resource cleanup, builders and fixtures for consistent test data

**Performance and Benchmarks:** Performance tests for user-facing operations, benchmark tests for regression prevention, memory and CPU usage verification for long-running operations

**Error Handling and Edge Cases:** All error conditions tested, boundary conditions and edge cases, input validation tests, concurrent access testing where applicable

### Code Quality Standards

**Test Code Quality:** Self-documenting tests with clear intent, same quality standards as production code, common utilities in `internal/testutil` package, reusable builders and helpers for complex scenarios

**Test Execution:** Unit tests complete in < 100ms, deterministic and non-flaky tests, no external service dependencies, parallel-safe execution

### Implementation Verification

**Automated Checks:** All tests pass in CI pipeline, coverage reports generated and reviewed, benchmark results compared against baselines, test code passes all linting rules

**Manual Verification:** Test coverage and quality reviewed during code review, manual verification of critical user workflows, cross-platform testing (macOS, Linux, Windows)

## Testing Framework Standards

**Go Testing Best Practices:** Use standard `testing` package, leverage `testify` assertions for better error messages, table-driven tests for multiple scenarios, proper test fixtures and cleanup

### Test Categories

#### 1. Unit Tests
```go
// Example: Unit test with proper structure
func TestServiceMethod_ValidInput_ReturnsExpectedResult(t *testing.T) {
    // Arrange
    svc := NewService()
    input := "valid-input"
    
    // Act
    result, err := svc.Method(input)
    
    // Assert
    require.NoError(t, err)
    assert.Equal(t, "expected-result", result)
}
```

#### 2. Integration Tests
```go
// Example: Integration test with proper setup/teardown
func TestUserWorkflow_AddMCP_Integration(t *testing.T) {
    // Setup
    tempDir := t.TempDir()
    model := testutil.NewTestModel().WithTempStorage(tempDir).Build()
    
    // Test complete workflow
    // ... test steps ...
    
    // Verify persistence
    savedData := loadFromStorage(tempDir)
    assert.Len(t, savedData, 1)
}
```

#### 3. Benchmark Tests
```go
// Example: Performance benchmark
func BenchmarkMCPService_FilterMCPs(b *testing.B) {
    mcps := generateLargeMCPDataset(1000)
    service := NewMCPService(mcps)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        service.FilterMCPs("github")
    }
}
```

## Test Coverage Requirements

**Minimum Thresholds:** Overall Project 85%, New Features 90%, Critical Components 95% (storage, UI handlers), Error Handling 100%

**Coverage Verification:**
```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
# Verify minimum thresholds
go test -cover ./... | grep -E "coverage: [0-9]+\.[0-9]+%" 
```

## Testing Tools and Dependencies

### Required Testing Packages
```go
// Standard testing
import "testing"

// Enhanced assertions
import "github.com/stretchr/testify/assert"
import "github.com/stretchr/testify/require"

// Test utilities
import "cc-mcp-manager/internal/testutil"
```

### Development Tools
- **Coverage Analysis:** Built-in Go coverage tools
- **Test Runners:** Standard `go test` with CI integration
- **Benchmarking:** Built-in Go benchmark tools
- **Mocking:** Interface-based dependency injection (no external mocks needed)

## Compliance Verification

### Pre-Commit Checks
- [ ] All tests pass locally
- [ ] Coverage meets minimum thresholds
- [ ] No test code lint violations
- [ ] Benchmark tests don't show regression

### Code Review Checklist
- [ ] New code has corresponding tests
- [ ] Tests cover both success and error scenarios
- [ ] Test names clearly describe what is being tested
- [ ] Complex scenarios use test builders from testutil
- [ ] Integration tests verify end-to-end workflows

### CI/CD Pipeline Requirements
- [ ] All test suites pass in CI
- [ ] Coverage reports generated and archived
- [ ] Performance benchmarks executed and compared
- [ ] Test results published for review

## Quality Gates

**Development Phase:** Code Complete (all tests written and passing locally), Review Ready (coverage meets thresholds, tests reviewed), Merge Ready (CI passes, all quality gates satisfied)

**Release Phase:** Integration Testing (all workflows tested end-to-end), Performance Validation (benchmarks meet performance criteria), Production Readiness (full test suite passes with confidence)

## Continuous Improvement

**Metrics Tracking:** Test execution time trends, coverage percentage over time, flaky test identification and resolution, performance benchmark trends

**Regular Reviews:** Monthly (test quality and coverage metrics), Per Story (testing approach effectiveness), Per Release (comprehensive testing retrospective)

## Enforcement and Accountability

**Development Team:** Write tests following standards, maintain existing test suites, review test quality in code reviews, report and address testing gaps

**Quality Assurance Team:** Validate standards compliance, provide guidance on testing approaches, review complex scenarios, maintain testing documentation

**Technical Lead:** Enforce standards in code reviews, approve exceptions to requirements, ensure team training on practices, monitor and report testing metrics

---

**Compliance Status:** This document defines the complete Testing Definition of Done that addresses the Process Review feedback. Implementation of these standards will move the Process Review from CONDITIONAL (85%) to APPROVED (95%+).

**Next Steps:**
1. Team training on new testing standards
2. Implementation of missing test categories
3. CI pipeline updates for automated verification
4. Regular compliance audits and metrics reporting