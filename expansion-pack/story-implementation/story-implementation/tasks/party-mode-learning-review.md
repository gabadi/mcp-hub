# Party Mode Learning Review

## Task Overview

**Agent:** architect
**Action Type:** collaborative-learning-review
**Duration:** 3-5 minutes
**Purpose:** Validate and prioritize technical learning items

## Purpose

Quick team validation of learning items from triage, prioritize for immediate action, assign ownership.

## Inputs

- `technical_debt` (array): Technical debt items from triage
- `architecture_improvements` (array): Architecture patterns to standardize
- `future_work` (array): Technical work for future planning

## Outputs

- `final_learning_items` (array): Prioritized actionable items with owners and timelines

## Instructions

### Step 1: Review Learning Items (1-2 minutes)

Review input learning items and validate:

**Keep items that:**

- Solve specific problems encountered
- Have measurable impact on development speed
- Address real technical debt discovered
- Improve future story implementation

**Remove items that:**

- Are theoretical improvements
- Duplicate existing technical debt tracking
- Lack specific next actions
- Don't address real problems found

### Step 2: Prioritize and Assign (2-3 minutes)

For each validated learning item:

**Prioritize:**

- HIGH: Blocks or slows future development
- MEDIUM: Improves development experience
- LOW: Nice-to-have optimization

**Assign ownership:**

- Technical debt → architect
- Tooling gaps → dev team
- Performance issues → architect + dev
- Standards → architect

### Step 3: Structure Final Learning Output

Return structured learning data:

```yaml
final_learning_items:
  - item: "Mobile responsive testing gaps"
    priority: "HIGH"
    owner: "architect"
    timeline: "Next sprint"
    action: "Add mobile testing to CI pipeline"
    impact: "Prevents mobile UX issues"

  - item: "CSS standardization needed"
    priority: "MEDIUM"
    owner: "architect"
    timeline: "Next epic"
    action: "Create CSS standards document"
    impact: "Consistent styling approach"

  - item: "Component documentation missing"
    priority: "LOW"
    owner: "dev team"
    timeline: "Backlog"
    action: "Add Storybook documentation"
    impact: "Improved developer experience"
```

**Max 5 actionable items with clear ownership and timelines.**

## Success Criteria

- [ ] All learning items validated for actionability
- [ ] Priority assigned based on development impact
- [ ] Clear ownership for each item
- [ ] Specific timelines for action items
- [ ] Ready for technical planning integration

## Integration Points

- **Input from:** capture-learning-triage task
- **Output to:** commit-and-prepare-pr task
- **Data format:** Structured YAML objects
