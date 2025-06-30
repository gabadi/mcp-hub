# Implement Story Development

**Agent:** dev
**Type:** BatchTask

## Purpose

Implement story requirements and document technical decisions made during implementation.

## Inputs

- Story file with acceptance criteria
- `blocking_issues` (array): Issues to fix from consolidation
- `technical_decisions` (array): Technical decisions to make

## Outputs

- `implementation_summary` (object): What was implemented and technical decisions made
- `implementation_status` (string): "Complete" | "Blocked"
- `quality_gates_status` (object): Status of project quality validation

## Instructions

### Workspace Enforcement

**CRITICAL**: All temporary files MUST be created in the designated workspace:

- Use workflow variables: `{temp_folder}`, `{issues_file}`, `{decisions_file}`, etc.
- NEVER create files outside of `{temp_folder}` directory
- File pattern: `{temp_folder}filename.md` where filename describes purpose
- Example: `{temp_folder}config-validation-epic{epic_number}-story{story_number}.md`

### Step 1: Implement Acceptance Criteria (Main work)

Read story file and implement all acceptance criteria:

- Follow project coding standards
- Write tests as required by project
- Use project's build/validation tools
- Fix any blocking issues from input

### Step 2: Make Technical Decisions (As needed)

For decisions in input data:

- Research options and trade-offs
- Make implementation choice based on project constraints
- Document rationale for future reference

### Step 3: Validate Quality Gates

Run project quality validation:

- Build process
- Test suite
- Linting/formatting
- Any project-specific checks
- Configuration AC validation (if story involves config changes)

#### Configuration AC Validation Process

When story involves configuration changes (ESLint, package.json, environment variables, etc.):

1. **Extract config-related acceptance criteria** from story file
2. **Verify each config AC** against actual implementation:
   - ESLint rules: Check configuration matches AC requirements
   - Package dependencies: Verify additions/changes match AC specifications
   - Environment variables: Confirm settings match AC requirements
   - Build configuration: Validate settings match AC requirements
3. **Document any AC-config mismatches** as blocking issues
4. **Include config validation status** in quality_gates_status output

**Example AC patterns to validate:**

- "ESLint must enforce X rule"
- "Add dependency Y for feature Z"
- "Environment variable X must be set to Y"
- "Build process must include step Z"

### Step 4: Structure Implementation Output

Return structured data:

```yaml
implementation_summary:
  components_built:
    - component: "PanicButton"
      purpose: "Emergency navigation component"
      files: ["PanicButton.svelte", "PanicButton.test.js"]
    - component: "BreadcrumbTrail"
      purpose: "Navigation history display"
      files: ["BreadcrumbTrail.svelte", "BreadcrumbTrail.test.js"]

  technical_decisions:
    - decision: "State management"
      choice: "Svelte store"
      rationale: "Reactive updates, simpler than props drilling"
    - decision: "Styling approach"
      choice: "CSS modules"
      rationale: "Component isolation, avoid class conflicts"

  issues_resolved:
    - issue: "CSS syntax error fixed"
      solution: "Corrected semicolon in line 23"
    - issue: "Test coverage added"
      solution: "Added unit tests for error handling"

implementation_status: "Complete"

quality_gates_status:
  build: "PASS"
  tests: "PASS"
  linting: "PASS"
  coverage: "95%"
  configuration_ac_validation: "PASS" # Added when config changes involved
  config_ac_details: # Added when config changes involved
    - ac: "ESLint no-console rule in production"
      status: "VERIFIED"
      config_location: ".eslintrc.js"
```

## Success Criteria

- [ ] All acceptance criteria implemented
- [ ] All blocking issues resolved
- [ ] Technical decisions documented with rationale
- [ ] Quality gates passing
- [ ] Structured output data provided

## Integration Points

- **Input from:** consolidate-review-feedback task
- **Output to:** validate-consolidated-fixes task
- **Data format:** Structured YAML objects
