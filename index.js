#!/usr/bin/env node

const { spawn } = require('child_process');
const path = require('path');
const os = require('os');
const fs = require('fs');

function getPlatformBinary() {
  const platform = os.platform();
  const arch = os.arch();
  
  let binaryName;
  
  // Map Node.js architecture names to Go architecture names
  const archMap = {
    'x64': 'amd64',
    'arm64': 'arm64'
  };
  
  const goArch = archMap[arch] || arch;
  
  if (platform === 'darwin') {
    binaryName = `mcp-hub-darwin-${goArch}`;
  } else if (platform === 'linux') {
    binaryName = `mcp-hub-linux-${goArch}`;
  } else if (platform === 'win32') {
    binaryName = `mcp-hub-windows-${goArch}.exe`;
  } else {
    console.error(`Error: Unsupported platform: ${platform}-${arch}`);
    console.error('');
    console.error('Supported platforms:');
    console.error('  - Linux (amd64, arm64)');
    console.error('  - macOS (amd64, arm64)');
    console.error('  - Windows (amd64, arm64)');
    console.error('');
    console.error('For manual installation, visit:');
    console.error('  https://github.com/gabadi/cc-mcp-manager/releases');
    process.exit(1);
  }
  
  return path.join(__dirname, 'bin', binaryName);
}

function validateBinaryExists(binaryPath) {
  if (!fs.existsSync(binaryPath)) {
    console.error(`Error: Binary not found at ${binaryPath}`);
    console.error('');
    console.error('This may indicate a corrupted installation or unsupported platform.');
    console.error('Please try reinstalling with: npm uninstall -g mcp-hub-tui && npx mcp-hub-tui');
    console.error('');
    console.error('If the issue persists, please report it at:');
    console.error('  https://github.com/gabadi/cc-mcp-manager/issues');
    process.exit(1);
  }
}

function main() {
  const binaryPath = getPlatformBinary();
  
  // Validate binary exists before attempting to execute
  validateBinaryExists(binaryPath);
  
  const child = spawn(binaryPath, process.argv.slice(2), {
    stdio: 'inherit',
    windowsHide: false
  });
  
  child.on('error', (err) => {
    console.error(`Error: Failed to execute mcp-hub-tui: ${err.message}`);
    console.error('');
    
    if (err.code === 'EACCES') {
      console.error('Permission denied. Try running with appropriate permissions:');
      console.error('  sudo npx mcp-hub-tui  # On Unix systems');
      console.error('  Or run as Administrator on Windows');
    } else if (err.code === 'ENOENT') {
      console.error('Binary not found. Please ensure the package is correctly installed.');
      console.error('Try reinstalling: npm uninstall -g mcp-hub-tui && npx mcp-hub-tui');
    } else {
      console.error('For troubleshooting, visit:');
      console.error('  https://github.com/gabadi/cc-mcp-manager/issues');
    }
    
    process.exit(1);
  });
  
  child.on('exit', (code) => {
    process.exit(code || 0);
  });
}

main();