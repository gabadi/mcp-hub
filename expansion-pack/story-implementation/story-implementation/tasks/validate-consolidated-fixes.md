# Validate Consolidated Fixes

**Agent:** sm
**Type:** BatchTask

## Purpose

Validate that consolidated fixes have been implemented correctly, approve or identify remaining issues.

## Inputs

- Story file with implementation details
- `{issues_file}` - Original issues to validate
- Implementation code changes

## Outputs

- Validation decision: APPROVED | NEEDS_FIXES
- `{validation_results_file}` - Issues remaining if any
- Story file updated with validation status

## Instructions

### Step 1: Validate Issues Resolved (3-5 minutes)

Check each issue from `{issues_file}`:

**For each blocking issue:**

- Verify specific error/problem is fixed
- Test functionality works as expected
- Confirm quality gates now pass

**For each technical decision:**

- Verify decision was made and implemented
- Check technical choice is documented with rationale

### Step 2: Quality Gate Validation (2-3 minutes)

Run project quality validation:

- Build: PASS/FAIL
- Tests: PASS/FAIL
- Linting: PASS/FAIL
- Any project-specific quality checks

### Step 3: Document Validation Results (1-2 minutes)

If issues remain, create `{validation_results_file}`:

```markdown
## Validation Results

### Remaining Issues

- [Specific issue] - [Evidence/test that shows it's not fixed]
- [Quality gate failure] - [Error message/output]

### Additional Fixes Needed

- [Fix description] - [Specific action required]
```

If all issues resolved:

```markdown
## Validation Results

**Status:** APPROVED
**All issues resolved:** YES
**Quality gates:** PASS
```

**Max 20 lines. Focus on specific issues remaining, not validation process details.**

### Step 4: Update Story Status (30 seconds)

Update story file:

```markdown
## Validation Complete

**Status:** [APPROVED | NEEDS_FIXES]
**Validated by:** SM
**Issues remaining:** [count or NONE]
```

## Success Criteria

- [ ] All original issues validated as resolved or remaining issues identified
- [ ] Quality gates verified as passing
- [ ] Clear approval decision made
- [ ] Specific guidance provided for any remaining work
- [ ] No verbose testing protocols or evidence documentation

## Integration Points

- **Input from:** implement-consolidated-fixes task
- **Output to:** capture-learning-triage task (if approved) or back to implementation (if fixes needed)
- **Files used:** `{issues_file}`
- **Files created:** `{validation_results_file}` (if issues remain)
