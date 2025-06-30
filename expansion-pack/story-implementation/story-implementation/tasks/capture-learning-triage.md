# Capture Learning Triage

## Task Overview

**Agent:** architect
**Action Type:** learning-triage
**Duration:** 5-8 minutes
**Purpose:** Extract actionable technical debt and improvement items

## Purpose

Analyze implementation experience, identify technical debt and improvement opportunities, output actionable items only.

## Inputs

- `implementation_summary` (object): Implementation details and decisions
- `quality_gates_status` (object): Quality validation results

## Outputs

- `technical_debt` (array): Technical debt items requiring fixes
- `architecture_improvements` (array): Architecture patterns to standardize
- `future_work` (array): Technical work for future planning

## Instructions

### Step 1: Identify Technical Debt (2-3 minutes)

From implementation experience, identify:

**Architecture improvements needed:**

- Performance bottlenecks discovered
- Code quality issues found
- Technical constraints that hindered development
- Integration problems encountered

**Technical decisions requiring follow-up:**

- Temporary solutions that need permanent fixes
- Libraries/approaches that should be standardized
- Missing tooling or automation gaps

### Step 2: Filter for Actionability (1-2 minutes)

**INCLUDE (actionable):**

- Specific technical debt items requiring fixes
- Performance issues with measurable impact
- Missing tooling that would improve development
- Architecture patterns that should be standardized

**EXCLUDE (not actionable):**

- General best practices already known
- Theoretical improvements without specific problems
- Process observations
- Congratulations or approval confirmations

### Step 3: Structure Learning Output (2-3 minutes)

Return structured data:

```yaml
technical_debt:
  - item: "Mobile responsive testing gaps"
    location: "CSS media queries"
    impact: "User experience on mobile devices"
    action: "Add mobile testing to CI pipeline"
    priority: "HIGH"
  - item: "Performance bottleneck in navigation"
    location: "BreadcrumbTrail component"
    impact: "Rendering delay >200ms"
    action: "Implement virtual scrolling"
    priority: "MEDIUM"

architecture_improvements:
  - pattern: "CSS standardization"
    current_issue: "Mixed CSS modules and inline styles"
    benefit: "Consistent styling approach"
    action: "Establish CSS standards document"
  - pattern: "State management consistency"
    current_issue: "Mix of stores and props"
    benefit: "Predictable data flow"
    action: "Define state management guidelines"

future_work:
  - work: "Automated accessibility testing"
    problem: "Manual a11y validation is slow"
    solution: "Integrate axe-core in CI"
    timeline: "Next sprint"
  - work: "Component documentation"
    problem: "No component usage examples"
    solution: "Add Storybook documentation"
    timeline: "Next epic"
```

**Max 10 actionable items total. Focus on technical improvements with clear next steps.**

## Success Criteria

- [ ] Technical debt items identified with specific actions
- [ ] Architecture improvements documented with clear benefits
- [ ] No theoretical or process-related items included
- [ ] All items have concrete next steps
- [ ] Ready for technical planning and prioritization

## Integration Points

- **Input from:** implement-story-development task
- **Output to:** party-mode-learning-review task
- **Data format:** Structured YAML objects
