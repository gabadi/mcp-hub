# Consolidate Review Feedback

## Task Overview

**Agent:** sm
**Action Type:** feedback-consolidation
**Duration:** 5-8 minutes
**Purpose:** Extract actionable items from review feedback

## Purpose

Analyze 5 review outputs, identify blocking issues and required decisions, output actionable items only.

## Inputs

- Architecture review results
- Business review results
- Process review results
- QA review results
- UX review results

## Outputs

- `blocking_issues` (array): Issues that block story completion
- `technical_decisions` (array): Technical decisions required
- `constraints_discovered` (array): Technical constraints found

## Instructions

### Step 1: Extract Blocking Issues (2-3 minutes)

Review all 5 feedback sources and identify:

**Issues that block story completion:**

- Build/compilation errors
- Test failures
- Missing core functionality
- Business requirement gaps

**Technical decisions needed:**

- Architecture choices requiring decision
- Library/framework selections
- Implementation approach options

### Step 2: Filter Out Non-Actionables (1-2 minutes)

**IGNORE (not actionable):**

- "Review completed successfully"
- "Meets requirements" confirmations
- Process compliance statements
- General approval comments
- Scope expansion suggestions

**INCLUDE (actionable):**

- Specific errors to fix
- Missing functionality to implement
- Technical choices to make
- Performance issues to resolve

### Step 3: Structure Output Data (2-3 minutes)

Return structured data:

```yaml
blocking_issues:
  - issue: "CSS syntax error in PanicButton.svelte:23"
    location: "PanicButton.svelte:23"
    type: "build_error"
  - issue: "Missing test coverage for error handling"
    location: "utils/auth.js"
    type: "missing_functionality"

technical_decisions:
  - decision: "State management approach"
    options: ["Svelte store", "props drilling"]
    criteria: "Simplicity vs consistency"
  - decision: "Styling approach"
    options: ["CSS modules", "Tailwind"]
    criteria: "Component isolation"

constraints_discovered:
  - constraint: "Must work with legacy auth system"
    impact: "Cannot use modern auth patterns"
  - constraint: "Performance requirement <200ms"
    impact: "Must optimize rendering"
```

**Max 20 actionable items total. Focus on actions, not analysis.**

## Success Criteria

- [ ] All blocking issues identified with specific locations
- [ ] All required decisions extracted with clear options
- [ ] Output contains only actionable items
- [ ] No process documentation or approval confirmations
- [ ] Ready for implementation phase

## Integration Points

- **Input from:** Round 1 reviews (5 agents)
- **Output to:** implement-consolidated-fixes task
- **Data format:** Structured YAML objects
