# Testing Standards - Definition of Done

**Document Version:** 1.0  
**Created:** 2025-06-30  
**Last Updated:** 2025-06-30  
**Owner:** Development Team  
**Approved By:** Technical Lead & QA Team

## Overview

This document defines the comprehensive Testing Definition of Done (DoD) standards for the CC MCP Manager project. These standards ensure production-ready quality assurance and full process compliance for all code changes.

## Testing Definition of Done Criteria

### Mandatory Requirements

Every code change MUST satisfy ALL of the following criteria before being considered complete:

#### 1. Unit Test Coverage
- **Minimum Coverage:** 85% line coverage for all new code
- **Quality Standard:** Tests must verify both happy path and error conditions
- **Isolation:** All unit tests must be independent and not rely on external dependencies
- **Naming Convention:** Test functions must follow `Test[MethodName]_[Scenario]` pattern
- **Assertions:** Each test must have clear, specific assertions with descriptive error messages

#### 2. Integration Test Coverage
- **End-to-End Workflows:** All user workflows must have integration tests
- **Component Integration:** Tests must verify interaction between major components
- **Error Handling:** Integration tests must cover error scenarios and recovery
- **Data Persistence:** Tests must verify data integrity across storage operations

#### 3. Test Organization and Structure
- **File Organization:** Test files must be co-located with source files (`*_test.go`)
- **Test Categories:** Tests must be categorized (unit, integration, benchmark)
- **Setup/Teardown:** Tests must properly clean up resources and test data
- **Test Data:** Use builders and fixtures for consistent test data

#### 4. Performance and Benchmark Tests
- **Critical Paths:** Performance tests for user-facing operations
- **Regression Prevention:** Benchmark tests to prevent performance degradation
- **Resource Usage:** Memory and CPU usage verification for long-running operations

#### 5. Error Handling and Edge Cases
- **Error Path Coverage:** All error conditions must be tested
- **Boundary Conditions:** Tests for edge cases and boundary values
- **Input Validation:** Tests for invalid input handling
- **Race Conditions:** Concurrent access testing where applicable

### Code Quality Standards

#### Test Code Quality
- **Readability:** Tests must be self-documenting with clear intent
- **Maintainability:** Test code follows same quality standards as production code
- **DRY Principle:** Common test utilities in `internal/testutil` package
- **Test Helpers:** Reusable builders and helpers for complex test scenarios

#### Test Execution Standards
- **Speed:** Unit tests must complete in < 100ms per test
- **Reliability:** Tests must be deterministic and not flaky
- **Isolation:** Tests must not depend on external services or files
- **Parallelization:** Tests must be safe to run in parallel

### Implementation Verification

#### Automated Checks
- **CI Integration:** All tests must pass in CI pipeline
- **Coverage Reports:** Coverage reports generated and reviewed
- **Performance Baselines:** Benchmark results compared against baselines
- **Lint Compliance:** Test code must pass all linting rules

#### Manual Verification
- **Code Review:** Test coverage and quality reviewed during code review
- **Functional Testing:** Manual verification of critical user workflows
- **Cross-Platform:** Testing on target platforms (macOS, Linux, Windows)

## Testing Framework Standards

### Go Testing Best Practices
- Use standard `testing` package for all tests
- Leverage `testify` assertions for better error messages
- Use table-driven tests for multiple scenarios
- Implement proper test fixtures and cleanup

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

### Minimum Coverage Thresholds
- **Overall Project:** 85% line coverage
- **New Features:** 90% line coverage
- **Critical Components:** 95% line coverage (storage, UI handlers)
- **Error Handling:** 100% error path coverage

### Coverage Verification
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

### Development Phase Gates
1. **Code Complete:** All tests written and passing locally
2. **Review Ready:** Coverage meets thresholds, tests reviewed
3. **Merge Ready:** CI passes, all quality gates satisfied

### Release Phase Gates
1. **Integration Testing:** All workflows tested end-to-end
2. **Performance Validation:** Benchmarks meet performance criteria
3. **Production Readiness:** Full test suite passes with confidence

## Continuous Improvement

### Metrics Tracking
- Test execution time trends
- Coverage percentage over time
- Flaky test identification and resolution
- Performance benchmark trends

### Regular Reviews
- **Monthly:** Review test quality and coverage metrics
- **Per Story:** Assess testing approach effectiveness
- **Per Release:** Comprehensive testing retrospective

## Enforcement and Accountability

### Development Team Responsibilities
- Write tests following these standards
- Maintain existing test suites
- Review test quality in code reviews
- Report and address testing gaps

### Quality Assurance Team Responsibilities
- Validate testing standards compliance
- Provide guidance on testing approaches
- Review complex testing scenarios
- Maintain testing documentation

### Technical Lead Responsibilities
- Enforce testing standards in code reviews
- Approve exceptions to testing requirements
- Ensure team training on testing practices
- Monitor and report testing metrics

---

**Compliance Status:** This document defines the complete Testing Definition of Done that addresses the Process Review feedback. Implementation of these standards will move the Process Review from CONDITIONAL (85%) to APPROVED (95%+).

**Next Steps:**
1. Team training on new testing standards
2. Implementation of missing test categories
3. CI pipeline updates for automated verification
4. Regular compliance audits and metrics reporting