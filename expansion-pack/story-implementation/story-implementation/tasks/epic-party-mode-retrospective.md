# Epic Party Mode Retrospective

## Task Overview

**Agent:** sm
**Action Type:** epic-retrospective
**Duration:** 15-20 minutes
**Purpose:** Extract actionable patterns and anti-patterns for next epic planning

## Purpose

Analyze completed epic, identify success patterns to reuse and anti-patterns to avoid, output LLM-optimized guidance for future epics.

## Inputs

- Epic file with 100% completion status
- All completed story files from the epic
- `{learning_final_file}` files from all stories

## Outputs

- Epic retrospective file optimized for LLM consumption by future epic planning
- Actionable patterns and anti-patterns only

## Instructions

### Step 1: Identify Success Patterns (5-7 minutes)

From all epic stories, identify patterns that worked well:

**Implementation patterns:**

- Technical approaches that delivered quickly
- Architecture decisions that scaled well
- Tool/library choices that improved velocity
- Testing strategies that caught issues early

**Process patterns:**

- Workflow selections that matched complexity correctly
- Agent assignments that were most effective
- Quality gates that prevented rework

### Step 2: Identify Anti-Patterns (5-7 minutes)

From story retrospectives and problems encountered:

**What to avoid:**

- Technical approaches that caused delays
- Architecture decisions that created technical debt
- Tools that slowed development
- Process choices that added unnecessary overhead

**Problem patterns:**

- Workflow mismatches (complex process for simple changes)
- Quality gate failures that could have been prevented
- Technical decisions that required later rework

### Step 3: Generate Epic Retrospective (5-8 minutes)

Create epic retrospective file:

```markdown
# Epic {epic_number} Retrospective

## Success Patterns (Reuse in Future Epics)

### Technical Patterns

- [Specific tech decision] - [Why it worked] - [When to apply]
- [Architecture choice] - [Benefit gained] - [Conditions for reuse]

### Process Patterns

- [Workflow selection] - [Story types where effective] - [Selection criteria]
- [Quality approach] - [Problems prevented] - [When to apply]

## Anti-Patterns (Avoid in Future Epics)

### Technical Anti-Patterns

- [Technical choice] - [Problem caused] - [Better alternative]
- [Architecture decision] - [Technical debt created] - [Preferred approach]

### Process Anti-Patterns

- [Process choice] - [Overhead created] - [Streamlined alternative]
- [Workflow mismatch] - [Efficiency loss] - [Correct selection criteria]

## Action Items for Next Epic

- [Owner]: [Specific action] - [Timeline] - [Success criteria]
- [Owner]: [Technical improvement] - [Implementation approach] - [Validation]

## Business Outcomes Achieved

- [Outcome 1]: [Metric] - [Business impact]
- [Outcome 2]: [Measurement] - [User benefit]
```

**Max 100 lines total. Focus on actionable guidance for future epic planning, not achievement celebration.**

## Success Criteria

- [ ] Success patterns identified with specific reuse conditions
- [ ] Anti-patterns documented with better alternatives
- [ ] Action items assigned with clear ownership and timelines
- [ ] Business outcomes quantified for completion validation
- [ ] Optimized for LLM consumption in future epic planning
- [ ] No verbose explanations or process documentation

## Integration Points

- **Input from:** All story files and learning summaries from completed epic
- **Output to:** Future epic planning and story implementation guidance
- **Primary consumer:** LLMs planning next epics and making technical decisions
