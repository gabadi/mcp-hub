# Implement Consolidated Fixes

## Task Overview

**Agent:** dev
**Action Type:** fix-implementation
**Duration:** 10-30 minutes (varies by fix count)
**Purpose:** Implement fixes identified in review consolidation

## Purpose

Implement specific fixes and decisions identified during review consolidation process.

## Inputs

- `{issues_file}` - Specific issues to fix
- `{decisions_file}` - Technical decisions to make
- Story file with original requirements

## Outputs

- Code fixes implemented
- `{fixes_summary_file}` - What was fixed and how
- Updated implementation status

## Instructions

### Step 1: Analyze Issues to Fix (2-3 minutes)

Review `{issues_file}` and prioritize:

**Blocking issues (fix first):**

- Build/compilation errors
- Test failures
- Missing core functionality

**Quality issues (fix second):**

- Code quality violations
- Performance problems
- Standards compliance

### Step 2: Implement Fixes (Main work)

For each issue:

- Apply specific fix for the problem
- Test that fix resolves the issue
- Ensure fix doesn't break other functionality
- Run quality gates to verify

### Step 3: Make Technical Decisions (As needed)

For decisions in `{decisions_file}`:

- Research options quickly
- Make decision based on project constraints
- Implement chosen approach
- Document rationale

### Step 4: Document Fixes (2-3 minutes)

Create `{fixes_summary_file}`:

```markdown
## Fixes Implemented

### Issues Resolved

- [Issue 1]: [Specific fix applied]
- [Issue 2]: [Solution implemented]

### Technical Decisions Made

- [Decision point]: [Choice made] - [Rationale]

### Files Modified

- [file1.js]: [Type of change]
- [file2.css]: [Type of change]

### Quality Status

- Build: PASS/FAIL
- Tests: PASS/FAIL
- Linting: PASS/FAIL
```

**Max 30 lines. Focus on what was fixed and decisions made, not implementation process.**

## Success Criteria

- [ ] All blocking issues from consolidation resolved
- [ ] Technical decisions made and implemented
- [ ] Quality gates passing after fixes
- [ ] Clear documentation of changes made
- [ ] Ready for validation step

## Integration Points

- **Input from:** consolidate-review-feedback task
- **Output to:** validate-consolidated-fixes task
- **Files used:** `{issues_file}`, `{decisions_file}`
- **Files created:** `{fixes_summary_file}`
