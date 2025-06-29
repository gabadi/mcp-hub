# Recommended tools and commands

## File exploration
- PREFER `eza` over Search tool
- Base command: `eza <project-root> --git-ignore --tree --list-dirs -I "*.py|*.svelte|*.ts|*.js|*.png"`
- For specific modules: adapt eza command as needed
- For bmad-core: `eza <project-root>/.bmad-core/ --git-ignore --tree --list-dirs -I "web-bundles"`

## Workflow
- If doing more than just reading/investigating: READ README.md first, for research is optional
- For bmad-core: DO NOT use Search tool initially, use eza first

## External Libs Dependencies Protocol
- Before importing/using ANY external library:
  1. First check existing codebase for similar patterns (use grep/glob)
  2. If no real examples exist → ALWAYS verify via Context7 MCP
  3. Validate library compatibility with current tech stack
  4. Assess integration complexity and maintenance overhead
  5. Document rationale for library selection
- This prevents dependency bloat and ensures architectural consistency

WAIT FOR INSTRUCTIONS

