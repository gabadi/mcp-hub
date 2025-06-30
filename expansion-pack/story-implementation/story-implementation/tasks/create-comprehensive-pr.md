# Create Simple PR

## Task Overview

**Agent:** po
**Action Type:** pr-creation-simple
**Duration:** 2-3 minutes
**Purpose:** Create concise, readable PR with essential information

## Purpose

Generate pull request with essential information only: what was implemented, why, key technical decisions, and breaking changes.

## Inputs

- Story definition and acceptance criteria
- `implementation_summary` (object): Technical decisions and components built
- `final_learning_items` (array): Priority technical debt items

## Outputs

- `pr_title` (string): Concise PR title
- `pr_description` (string): Simple PR description
- `story_update` (object): Story file updates with technical decisions

## Instructions

### Step 1: Generate PR Title (30 seconds)

Create concise title:

```
[Epic{epic_number}.Story{story_number}] {what_was_implemented}
```

### Step 2: Create Simple PR Description (1-2 minutes)

Generate structured PR content:

```yaml
pr_description: |
  ## What
  {one_sentence_what_was_implemented}

  ## Why
  {one_sentence_business_reason}

  ## Technical decisions
  - {decision_1_with_rationale}
  - {decision_2_with_rationale}
  - {decision_3_with_rationale}

  ## Breaking changes
  {none_or_specific_changes}

  ## Testing
  {ci_validation_status}

  ---
  🤖 Generated with [Claude Code](https://claude.ai/code)
```

### Step 3: Structure Story Updates (30 seconds)

Prepare story file updates:

```yaml
story_update:
  status: "Complete"
  technical_decisions:
    - decision: "State management approach"
      choice: "Svelte store"
      rationale: "Reactive updates, simpler than props drilling"

  technical_debt_identified:
    - item: "Mobile responsive testing gaps"
      priority: "HIGH"
      owner: "architect"
      timeline: "Next sprint"

  pr_url: "{pr_url_to_be_set_by_workflow}"
```

**Max 150 words total PR description. Focus on what was implemented and key technical choices.**

## Success Criteria

- [ ] PR title is concise and descriptive
- [ ] PR description contains max 150 words
- [ ] Technical decisions clearly explained
- [ ] Breaking changes explicitly stated
- [ ] Story updates structured for file insertion
- [ ] No verbose business context or process documentation

## Integration Points

- **Input from:** All previous workflow steps
- **Output to:** Workflow for PR creation and story file updates
- **Data format:** Structured objects for workflow consumption
