# Commit and Prepare PR

## Task Overview

**Agent:** dev
**Action Type:** git-commit-and-pr-preparation
**Duration:** 3-5 minutes
**Purpose:** Commit changes and extract technical decisions from temp folder

## Purpose

Commit implementation changes with simple message, extract technical decisions from temp folder to story file.

## Inputs

- All implementation code changes
- `{implementation_file}` - Technical decisions made
- `{learning_final_file}` - Technical debt items

## Outputs

- Git commit with implementation changes
- Story file updated with technical decisions
- Ready for PR creation

## Instructions

### Step 1: Pre-Commit Validation (1 minute)

Verify commit readiness:

- [ ] All quality gates passing
- [ ] Implementation code complete
- [ ] No uncommitted changes remaining

### Step 2: Extract Technical Decisions (2 minutes)

From `{temp_folder}` temp folder, extract actionable information to story file:

```markdown
## Implementation Completed

**Status:** Complete
**Quality Gates:** PASS

### Technical Decisions Made

- [Decision 1]: [Rationale]
- [Decision 2]: [Rationale]

### Technical Debt Identified

- [High priority item]: [Owner] - [Timeline]
- [Medium priority item]: [Owner] - [Timeline]
```

### Step 3: Generate Simple Commit (1-2 minutes)

Create concise commit message:

```
[Epic{epic_number}.Story{story_number}] {what_was_implemented}

Technical decisions:
- {decision_1_brief}
- {decision_2_brief}

Quality gates: PASS
```

**Max 5 lines total. Focus on what was implemented and key technical choices.**

### Step 4: Commit Changes (30 seconds)

```bash
git add .
git commit -m "{simple_commit_message}"
```

## Success Criteria

- [ ] Code changes committed with clear, concise message
- [ ] Technical decisions extracted from temp folder to story file
- [ ] Story status updated to Complete
- [ ] Ready for PR creation step
- [ ] No verbose business context or process documentation

## Integration Points

- **Input from:** implement-story-development, learning review tasks
- **Output to:** create-comprehensive-pr task
- **Files used:** `{implementation_file}`, `{learning_final_file}`
- **Files updated:** Story file with technical decisions section
