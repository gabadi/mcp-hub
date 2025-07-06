# Story 3.1: NPX Distribution Implementation for mcp-hub-tui

## Status: Ready for Development

## Story Approved for Development

**Status:** Approved  
**Approved by:** SM (Scrum Master Agent)  
**Ready for:** Development  
**Created by:** SM (Scrum Master Agent)  
**MVP Focus:** Applied (Core NPX distribution functionality)

## Story

As a developer using Claude Code with MCPs,
I want to install and run mcp-hub using `npx mcp-hub-tui` command,
so that I can easily access the MCP management tool without manual installation or complex setup.

## Business Context

This story implements NPX distribution for mcp-hub, enabling developers to run the tool directly using `npx mcp-hub-tui` without requiring manual binary downloads or installation procedures. This addresses the key user experience gap in MCP tooling distribution and makes mcp-hub easily accessible to the broader Claude Code developer community.

**Epic Context:** Epic 3 - Distribution & Package Management  
**Prerequisites:** Core mcp-hub functionality (Epic 1 & 2) - COMPLETED  
**Timeline:** 8-12 hours (NPX setup, build automation, testing)  
**Strategic Value:** Eliminates installation friction and provides seamless access to mcp-hub across platforms

## Acceptance Criteria

### AC1: NPX Package Structure & Configuration
- **Given** the mcp-hub project needs NPX distribution
- **When** the npm package is configured
- **Then** a complete npm package structure is established:
  - `package.json` with proper npm configuration for "mcp-hub-tui"
  - `bin` directory with platform-specific executables
  - Platform detection and binary selection logic
  - Proper package metadata and dependencies
- **And** the package supports all 6 target platforms:
  - `@mcp-hub-tui/linux-amd64`, `@mcp-hub-tui/linux-arm64`
  - `@mcp-hub-tui/darwin-amd64`, `@mcp-hub-tui/darwin-arm64`
  - `@mcp-hub-tui/windows-amd64`, `@mcp-hub-tui/windows-arm64`
- **And** package version is synchronized with Go module version

### AC2: NPX Command Execution
- **Given** a developer has Node.js installed
- **When** they run `npx mcp-hub-tui`
- **Then** the command executes successfully:
  - Downloads the package on first run
  - Detects the correct platform automatically
  - Executes the appropriate binary for their system
  - Displays the mcp-hub TUI interface
- **And** subsequent runs use cached package for fast execution
- **And** the command works identically to direct binary execution

### AC3: Cross-Platform Binary Support
- **Given** the npm package contains multiple platform binaries
- **When** `npx mcp-hub-tui` is executed on any supported platform
- **Then** the correct binary is selected and executed:
  - Platform detection works reliably
  - Binary permissions are set correctly (Unix systems)
  - Windows executable runs without additional setup
  - Unsupported platforms show clear error messages
- **And** binary selection is transparent to the user
- **And** no platform-specific configuration is required

### AC4: Automated Build & Release Process
- **Given** the project uses GoReleaser for releases
- **When** a new version is tagged and released
- **Then** the automated process:
  - Builds binaries for all 6 target platforms
  - Creates npm package with all binaries
  - Publishes to npm registry automatically
  - Updates GitHub release with npm package info
- **And** the process is fully automated via GitHub Actions
- **And** manual intervention is only required for version tagging

### AC5: Error Handling & User Experience
- **Given** various error scenarios can occur
- **When** issues are encountered during NPX execution
- **Then** clear error messages are provided:
  - Unsupported platform notifications
  - Binary execution failures
  - Permission errors with resolution guidance
  - Network/download failures with retry suggestions
- **And** error messages include helpful troubleshooting information
- **And** fallback instructions for manual installation are provided

## Technical Implementation Requirements

### NPM Package Structure
```
mcp-hub-tui/
├── package.json
├── bin/
│   ├── mcp-hub-linux-amd64
│   ├── mcp-hub-linux-arm64
│   ├── mcp-hub-darwin-amd64
│   ├── mcp-hub-darwin-arm64
│   ├── mcp-hub-windows-amd64.exe
│   └── mcp-hub-windows-arm64.exe
├── index.js (platform detection & execution)
└── README.md
```

### Package.json Configuration
```json
{
  "name": "mcp-hub-tui",
  "version": "1.0.0",
  "description": "TUI for managing Claude MCP configurations",
  "bin": {
    "mcp-hub-tui": "./index.js"
  },
  "files": [
    "bin/",
    "index.js",
    "README.md"
  ],
  "engines": {
    "node": ">=14.0.0"
  },
  "keywords": ["claude", "mcp", "terminal", "tui", "cli"],
  "repository": {
    "type": "git",
    "url": "git+https://github.com/gabadi/cc-mcp-manager.git"
  },
  "license": "MIT",
  "homepage": "https://github.com/gabadi/cc-mcp-manager#readme"
}
```

### Platform Detection Logic (index.js)
```javascript
#!/usr/bin/env node

const { spawn } = require('child_process');
const path = require('path');
const os = require('os');

function getPlatformBinary() {
  const platform = os.platform();
  const arch = os.arch();
  
  let binaryName;
  
  if (platform === 'darwin') {
    binaryName = arch === 'arm64' ? 'mcp-hub-darwin-arm64' : 'mcp-hub-darwin-amd64';
  } else if (platform === 'linux') {
    binaryName = arch === 'arm64' ? 'mcp-hub-linux-arm64' : 'mcp-hub-linux-amd64';
  } else if (platform === 'win32') {
    binaryName = arch === 'arm64' ? 'mcp-hub-windows-arm64.exe' : 'mcp-hub-windows-amd64.exe';
  } else {
    console.error(`Unsupported platform: ${platform}-${arch}`);
    process.exit(1);
  }
  
  return path.join(__dirname, 'bin', binaryName);
}

function main() {
  const binaryPath = getPlatformBinary();
  
  const child = spawn(binaryPath, process.argv.slice(2), {
    stdio: 'inherit',
    windowsHide: false
  });
  
  child.on('error', (err) => {
    console.error(`Failed to execute mcp-hub-tui: ${err.message}`);
    process.exit(1);
  });
  
  child.on('exit', (code) => {
    process.exit(code);
  });
}

main();
```

### GoReleaser Configuration Updates
```yaml
# .goreleaser.yml additions
builds:
  - id: mcp-hub
    main: ./main.go
    binary: mcp-hub
    goos:
      - linux
      - darwin
      - windows
    goarch:
      - amd64
      - arm64
    env:
      - CGO_ENABLED=0
    ldflags:
      - -s -w -X main.version={{.Version}}

archives:
  - id: binaries
    builds:
      - mcp-hub
    format: tar.gz
    format_overrides:
      - goos: windows
        format: zip
    files:
      - LICENSE
      - README.md

# New NPM package configuration
npm:
  - ids:
      - mcp-hub
    package_name: mcp-hub-tui
    registry: https://registry.npmjs.org/
    folder: npm-package
    package_template: |
      {
        "name": "mcp-hub-tui",
        "version": "{{ .Version }}",
        "description": "TUI for managing Claude MCP configurations",
        "bin": {
          "mcp-hub-tui": "./index.js"
        },
        "files": [
          "bin/",
          "index.js",
          "README.md"
        ],
        "engines": {
          "node": ">=14.0.0"
        },
        "keywords": ["claude", "mcp", "terminal", "tui", "cli"],
        "repository": {
          "type": "git",
          "url": "git+https://github.com/gabadi/cc-mcp-manager.git"
        },
        "license": "MIT",
        "homepage": "https://github.com/gabadi/cc-mcp-manager#readme"
      }
```

### GitHub Actions CI/CD Updates
```yaml
# .github/workflows/release.yml
name: Release

on:
  push:
    tags:
      - 'v*'

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
      
      - uses: actions/setup-go@v5
        with:
          go-version: '1.24'
      
      - uses: actions/setup-node@v4
        with:
          node-version: '18'
          registry-url: 'https://registry.npmjs.org'
      
      - name: Run GoReleaser
        uses: goreleaser/goreleaser-action@v5
        with:
          version: latest
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          NPM_TOKEN: ${{ secrets.NPM_TOKEN }}
```

## Manual Setup Requirements

### NPM Account & Authentication
**Developer must complete these steps manually:**

1. **Create NPM Account:**
   - Visit https://www.npmjs.com/signup
   - Create account with email verification
   - Choose appropriate account type (free tier sufficient)

2. **Generate NPM Access Token:**
   - Login to NPM account
   - Go to "Access Tokens" in account settings
   - Create "Automation" token with "Publish" permissions
   - Copy token value securely

3. **Configure GitHub Repository Secrets:**
   - Navigate to repository Settings > Secrets and variables > Actions
   - Add new repository secret: `NPM_TOKEN`
   - Paste NPM access token as the value
   - Ensure secret is available to GitHub Actions

4. **Verify Package Name Availability:**
   - Check https://www.npmjs.com/package/mcp-hub-tui
   - Confirm package name is available or owned by correct account
   - Reserve package name if necessary

### Repository Configuration
**Developer must configure these settings:**

1. **GitHub Actions Permissions:**
   - Repository Settings > Actions > General
   - Enable "Allow GitHub Actions to create and approve pull requests"
   - Set "Workflow permissions" to "Read and write permissions"

2. **Release Settings:**
   - Enable "Automatically delete head branches" (optional)
   - Configure branch protection rules for main branch

3. **GoReleaser Configuration:**
   - Ensure `.goreleaser.yml` is in repository root
   - Validate configuration with `goreleaser check`

## Step-by-Step Implementation Guide

### Phase 1: NPM Package Structure Setup (2-3 hours)
1. **Create package.json:**
   - Configure package metadata and dependencies for "mcp-hub-tui"
   - Set up bin configuration for npx execution with "mcp-hub-tui" command
   - Define files to include in package

2. **Create mcp-hub-tui wrapper script (index.js):**
   - Implement platform and architecture detection
   - Add binary selection logic for mcp-hub binaries
   - Include error handling for unsupported platforms

3. **Set up bin directory structure:**
   - Create bin/ directory for platform binaries
   - Document binary naming convention
   - Add placeholder files for testing

### Phase 2: GoReleaser Integration (3-4 hours)
1. **Update .goreleaser.yml:**
   - Add npm package configuration
   - Configure binary builds for all platforms
   - Set up package template and metadata

2. **Create NPM package template:**
   - Design package structure for GoReleaser
   - Configure binary placement and permissions
   - Set up package documentation

3. **Test local package generation:**
   - Run `goreleaser build --snapshot --clean`
   - Verify binary generation for all platforms
   - Test mcp-hub-tui package structure and contents

### Phase 3: GitHub Actions Automation (2-3 hours)
1. **Create/update release workflow:**
   - Configure GitHub Actions for automated releases
   - Set up Node.js and Go environments
   - Add NPM authentication and publishing

2. **Configure repository secrets:**
   - Add NPM_TOKEN to GitHub secrets
   - Verify GitHub Actions permissions
   - Test workflow with dry-run

3. **Test end-to-end release process:**
   - Create test tag and release
   - Verify binary building and NPM publishing
   - Validate package installation with `npx mcp-hub-tui`

### Phase 4: Testing & Validation (2-3 hours)
1. **Platform-specific testing:**
   - Test `npx mcp-hub-tui` execution on Linux (amd64, arm64)
   - Test `npx mcp-hub-tui` execution on macOS (Intel, Apple Silicon)
   - Test `npx mcp-hub-tui` execution on Windows (amd64, arm64)

2. **Error scenario testing:**
   - Test unsupported platform handling
   - Test binary execution failures
   - Test network/download error scenarios

3. **Performance & UX testing:**
   - Verify first-run download experience
   - Test cached execution performance
   - Validate error message clarity

### Phase 5: Documentation & Release (1-2 hours)
1. **Update project documentation:**
   - Add NPX installation instructions to README (npx mcp-hub-tui)
   - Document troubleshooting for common issues
   - Create release notes for NPX feature

2. **Create first official release:**
   - Tag version for NPX distribution
   - Verify automated release process
   - Test final package installation

## Testing Requirements

### Functional Testing
- [ ] NPX package installs correctly on all supported platforms
- [ ] Binary selection logic works for all platform/architecture combinations
- [ ] Error handling provides clear messages for unsupported platforms
- [ ] Package caching works properly for subsequent runs
- [ ] All command-line arguments pass through correctly

### Integration Testing
- [ ] GoReleaser builds npm package correctly
- [ ] GitHub Actions workflow publishes to npm registry
- [ ] Package metadata is accurate and complete
- [ ] Binary permissions are set correctly on Unix systems
- [ ] Windows executable runs without additional setup

### Manual Testing Checklist
- [ ] Test `npx mcp-hub-tui` on Linux amd64
- [ ] Test `npx mcp-hub-tui` on Linux arm64
- [ ] Test `npx mcp-hub-tui` on macOS Intel
- [ ] Test `npx mcp-hub-tui` on macOS Apple Silicon
- [ ] Test `npx mcp-hub-tui` on Windows amd64
- [ ] Test `npx mcp-hub-tui` on Windows arm64
- [ ] Test error handling for unsupported platforms
- [ ] Test network failure scenarios
- [ ] Test permission error scenarios
- [ ] Verify package caching behavior

## Dependencies

### Prerequisites
- ✅ Core mcp-hub functionality (Epic 1 & 2) - COMPLETED
- ✅ Go project with working binary compilation
- ✅ GitHub repository with Actions enabled
- ✅ GoReleaser configuration (basic)

### Required External Accounts
- [ ] NPM account with publish permissions
- [ ] NPM access token with automation permissions
- [ ] GitHub repository with Actions permissions

### Technical Dependencies
- GoReleaser v1.21+ for NPM package support
- Node.js 14+ for NPX execution
- GitHub Actions for automated publishing
- NPM registry for package distribution

## Risk Assessment

### Technical Risks
- **Medium Risk:** GoReleaser NPM integration complexity
- **Low Risk:** Platform detection accuracy
- **Low Risk:** Binary permissions on Unix systems
- **Low Risk:** NPM package publishing automation

### Business Risks
- **Low Risk:** NPM package name "mcp-hub-tui" availability
- **Low Risk:** NPM registry service reliability
- **Low Risk:** GitHub Actions service availability

### Mitigation Strategies
- Test GoReleaser configuration thoroughly before release
- Implement comprehensive platform detection with fallbacks
- Provide clear error messages and manual installation instructions
- Set up monitoring for NPM package download metrics
- Create documentation for manual binary installation as backup
- Ensure wrapper script name (mcp-hub-tui) aligns with package name

## Success Metrics

### Primary Success Criteria
- [ ] `npx mcp-hub-tui` command works on all 6 supported platforms
- [ ] Package installs and executes within 30 seconds on first run
- [ ] Cached execution completes within 5 seconds
- [ ] Error handling provides actionable guidance
- [ ] Automated release process requires zero manual intervention

### Quality Metrics
- [ ] NPM package size under 50MB (all binaries included)
- [ ] Platform detection accuracy: 100%
- [ ] Error message clarity: User-friendly and actionable
- [ ] Documentation completeness: Installation and troubleshooting

## Definition of Done

### Functional Requirements
- [ ] All 5 acceptance criteria pass validation testing
- [ ] `npx mcp-hub-tui` command executes successfully on all supported platforms
- [ ] Automated build and release process works end-to-end
- [ ] Error handling provides clear guidance for all failure scenarios
- [ ] Package installation and execution performance meets targets

### Quality Requirements
- [ ] Comprehensive testing across all platform/architecture combinations
- [ ] Error scenarios tested and validated
- [ ] Documentation is complete and accurate
- [ ] NPM package metadata is properly configured
- [ ] Release automation requires no manual intervention

### Technical Requirements
- [ ] GoReleaser configuration includes NPM package generation
- [ ] GitHub Actions workflow publishes to NPM registry
- [ ] Binary permissions are correctly set for Unix systems
- [ ] Package structure follows NPM best practices
- [ ] Version synchronization between Go module and NPM package

## Implementation Constraints

### Technical Constraints
- Must support Node.js 14+ (NPX compatibility)
- Binary size limitations for NPM package distribution
- Platform detection must be reliable across all supported systems
- Error messages must be clear and actionable
- Package installation must be fast and efficient

### Business Constraints
- NPM account and token setup requires manual configuration
- GitHub Actions permissions must be properly configured
- Package name availability dependent on NPM registry
- Release process must be fully automated after initial setup

## Technical Debt Considerations

### Potential Technical Debt
- **NPM Package Maintenance:** Ongoing maintenance of NPM package structure
- **Binary Distribution:** Managing multiple platform binaries in single package
- **Version Synchronization:** Keeping Go module and NPM package versions aligned
- **Platform Support:** Adding new platforms requires NPM package updates

### Mitigation Strategies
- Implement automated version synchronization
- Create comprehensive testing for new platform additions
- Document NPM package maintenance procedures
- Set up monitoring for package download and usage metrics

## Future Enhancements

### Planned Improvements
- Selective binary download (download only required platform)
- Binary signature verification for security
- Automatic updates notification system
- Enhanced error reporting and diagnostics

### Considerations for Future Stories
- Package size optimization strategies
- Alternative distribution methods (Homebrew, Chocolatey)
- Enhanced platform detection and fallback mechanisms
- Integration with package managers beyond NPM

## Notes & Considerations

### Design Decisions
- **All-in-one package:** Include all platform binaries in single NPM package (mcp-hub-tui) for simplicity
- **Platform detection:** Use Node.js built-in modules for reliable platform identification
- **Error handling:** Provide clear error messages with actionable resolution steps
- **Automation:** Fully automated release process after initial manual setup
- **Naming Strategy:** NPM package uses "mcp-hub-tui" while internal binary remains "mcp-hub"

### Development Notes
- NPM package approach chosen for broad compatibility and ease of use
- Package name "mcp-hub-tui" chosen to distinguish from potential binary-only "mcp-hub" package
- GoReleaser NPM integration provides seamless automation
- Platform detection logic must be robust and maintainable
- Error scenarios require comprehensive testing and validation

### User Experience Considerations
- First-run experience should be smooth despite package download
- Command name "npx mcp-hub-tui" should be intuitive and memorable
- Cached execution should be nearly instantaneous
- Error messages should guide users to successful resolution
- Documentation should cover common installation and usage scenarios

---

**Created by:** SM (Scrum Master Agent)  
**Epic:** Epic 3 - Distribution & Package Management  
**MVP Focus:** Applied (Core NPX distribution functionality)  
**Target Timeline:** 8-12 hours  
**Dependencies:** Epic 1 & 2 (COMPLETED), Manual NPM setup required  
**Status:** Ready for Development

---

## Manual Setup Checklist for Developer

**Complete these steps before starting development:**

### NPM Account Setup
- [ ] Create NPM account at https://www.npmjs.com/signup
- [ ] Verify email address and complete account setup
- [ ] Generate NPM access token with "Automation" and "Publish" permissions
- [ ] Add NPM_TOKEN to GitHub repository secrets

### Repository Configuration
- [ ] Enable GitHub Actions with read/write permissions
- [ ] Verify GoReleaser configuration is present
- [ ] Check "mcp-hub-tui" package name availability on NPM registry
- [ ] Configure branch protection rules (optional)

### Development Environment
- [ ] Install Node.js 14+ for local testing
- [ ] Install GoReleaser for local package generation
- [ ] Test basic NPX functionality with existing packages
- [ ] Verify all development dependencies are available
- [ ] Test npx command execution with mcp-hub-tui naming

**Only proceed with development after completing all manual setup steps.**