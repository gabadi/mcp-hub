# Approve Story for Development

## Purpose

Product Owner validation and approval of story for development readiness.

## Inputs

- `story_file`: Path to the story file requiring approval
- `epic_number`: Epic number for alignment validation
- `approval_threshold`: Minimum approval percentage required (e.g., 90)

## Outputs

- Story approval decision: APPROVED | NEEDS_REVISION
- Story file updated with approval status
- `{approval_notes_file}` (if issues found)

## Instructions

### Workspace Enforcement

**CRITICAL**: All temporary files MUST be created using workflow variables:

- Use `{approval_notes_file}` for approval issues
- NEVER create files outside designated workspace
- Pattern: Workflow variables resolve to `{temp_folder}` directory

### Step 1: Load Story and Epic Context (1-2 minutes)

- Read the complete story file
- Read the parent epic file for business context
- Extract user story, acceptance criteria, and business context

### Step 2: Validate Story Readiness (2-3 minutes)

Check essential approval criteria:

**Business alignment:**

- Story supports epic business objectives
- User value is clear and measurable
- Acceptance criteria are specific and testable

**Development readiness:**

- Requirements are clear and unambiguous
- Scope is appropriate for single story
- Dependencies are identified

### Step 3: Make Approval Decision (1 minute)

**CRITICAL: Story must achieve minimum {approval_threshold}% approval threshold to proceed to development.**

Evaluate overall story readiness score based on validation criteria:

- Business alignment: Pass/Fail
- Development readiness: Pass/Fail
- All acceptance criteria clear: Pass/Fail
- Scope appropriateness: Pass/Fail
- Dependencies identified: Pass/Fail

**Minimum {approval_threshold}% ({approval_threshold \* 5 / 100}/5) criteria must pass for approval.**

**If story meets {approval_threshold}% threshold:**
Update story file:

```markdown
## Story Approved for Development

**Status:** Approved ({approval_threshold}%+ threshold met)
**Approved by:** PO
**Ready for:** Development
**Approval Score:** [X/5 criteria passed]
```

**If story falls below {approval_threshold}% threshold:**
Create `{approval_notes_file}` and mark as NEEDS_REVISION:

```markdown
## Story Approval Issues

**Approval Score:** [X/5 criteria passed] - **Below {approval_threshold}% threshold**

### Failed Criteria

- [List specific criteria that failed]
- [Include actionable guidance for each]

### Required Actions

- [Specific steps to reach {approval_threshold}% threshold]
```

**Max 15 lines. Focus on specific issues to fix, not approval process details.**

## Success Criteria

- [ ] Story business alignment validated
- [ ] Acceptance criteria are clear and testable
- [ ] Development scope is appropriate
- [ ] Clear approval decision made
- [ ] Issues documented with specific guidance if needed

## Integration Points

- **Input from:** Story creation process
- **Output to:** implement-story-development task (if approved)
- **Files created:** `{approval_notes_file}` (if issues found)
