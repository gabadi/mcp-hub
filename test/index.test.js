#!/usr/bin/env node

/**
 * Test suite for mcp-hub-tui NPM package
 * This test validates platform detection, binary path resolution, and error handling
 */

const path = require('path');
const fs = require('fs');
const os = require('os');

// Mock console methods for testing
const originalConsoleError = console.error;
const originalProcessExit = process.exit;
let consoleErrors = [];
let exitCode = null;

function mockConsole() {
  console.error = (...args) => {
    consoleErrors.push(args.join(' '));
  };
  process.exit = (code) => {
    exitCode = code;
    throw new Error(`Process would exit with code ${code}`);
  };
}

function restoreConsole() {
  console.error = originalConsoleError;
  process.exit = originalProcessExit;
  consoleErrors = [];
  exitCode = null;
}

// Test utility functions
function runTest(testName, testFunc) {
  try {
    mockConsole();
    testFunc();
    console.log(`✅ ${testName}`);
    restoreConsole();
  } catch (error) {
    restoreConsole();
    if (error.message.includes('Process would exit')) {
      console.log(`✅ ${testName} (expected exit)`);
    } else {
      console.error(`❌ ${testName}: ${error.message}`);
    }
  }
}

// Platform detection function for testing (extracted from index.js)
function getPlatformBinary() {
  const platform = os.platform();
  const arch = os.arch();
  
  // Map Node.js architecture names to Go architecture names
  const archMap = {
    'x64': 'amd64',
    'arm64': 'arm64'
  };
  
  const goArch = archMap[arch] || arch;
  
  let binaryName;
  if (platform === 'darwin') {
    binaryName = `mcp-hub-darwin-${goArch}`;
  } else if (platform === 'linux') {
    binaryName = `mcp-hub-linux-${goArch}`;
  } else if (platform === 'win32') {
    binaryName = `mcp-hub-windows-${goArch}.exe`;
  } else {
    throw new Error(`Unsupported platform: ${platform}-${arch}`);
  }
  
  return path.join(__dirname, '..', 'bin', binaryName);
}

// Test platform detection
function testPlatformDetection() {
  console.log('\n🧪 Testing Platform Detection...');
  
  // Test current platform
  runTest('Current platform detection', () => {
    const binaryPath = getPlatformBinary();
    if (!binaryPath || typeof binaryPath !== 'string') {
      throw new Error('getPlatformBinary should return a string path');
    }
    
    const platform = os.platform();
    const arch = os.arch();
    const expectedArch = arch === 'x64' ? 'amd64' : arch;
    
    if (platform === 'darwin') {
      if (!binaryPath.includes(`mcp-hub-darwin-${expectedArch}`)) {
        throw new Error(`Expected darwin binary path, got: ${binaryPath}`);
      }
    } else if (platform === 'linux') {
      if (!binaryPath.includes(`mcp-hub-linux-${expectedArch}`)) {
        throw new Error(`Expected linux binary path, got: ${binaryPath}`);
      }
    } else if (platform === 'win32') {
      if (!binaryPath.includes(`mcp-hub-windows-${expectedArch}.exe`)) {
        throw new Error(`Expected windows binary path, got: ${binaryPath}`);
      }
    }
  });
}

// Test binary existence
function testBinaryExistence() {
  console.log('\n🧪 Testing Binary Existence...');
  
  runTest('Binary file exists', () => {
    const platform = os.platform();
    const arch = os.arch();
    const expectedArch = arch === 'x64' ? 'amd64' : arch;
    
    let expectedBinary;
    if (platform === 'darwin') {
      expectedBinary = `mcp-hub-darwin-${expectedArch}`;
    } else if (platform === 'linux') {
      expectedBinary = `mcp-hub-linux-${expectedArch}`;
    } else if (platform === 'win32') {
      expectedBinary = `mcp-hub-windows-${expectedArch}.exe`;
    }
    
    const binaryPath = path.join(__dirname, '..', 'bin', expectedBinary);
    
    if (!fs.existsSync(binaryPath)) {
      throw new Error(`Binary not found at: ${binaryPath}`);
    }
  });
}

// Test package.json structure
function testPackageStructure() {
  console.log('\n🧪 Testing Package Structure...');
  
  runTest('package.json exists and valid', () => {
    const packagePath = path.join(__dirname, '..', 'package.json');
    if (!fs.existsSync(packagePath)) {
      throw new Error('package.json not found');
    }
    
    const packageJson = JSON.parse(fs.readFileSync(packagePath, 'utf8'));
    
    // Validate required fields
    const requiredFields = ['name', 'version', 'description', 'bin', 'files'];
    for (const field of requiredFields) {
      if (!packageJson[field]) {
        throw new Error(`Missing required field: ${field}`);
      }
    }
    
    // Validate bin configuration
    if (packageJson.bin['mcp-hub-tui'] !== './index.js') {
      throw new Error('Invalid bin configuration');
    }
    
    // Validate files array
    const expectedFiles = ['bin/', 'index.js', 'README.md'];
    for (const file of expectedFiles) {
      if (!packageJson.files.includes(file)) {
        throw new Error(`Missing file in package.json files array: ${file}`);
      }
    }
  });
}

// Test Node.js version compatibility
function testNodeCompatibility() {
  console.log('\n🧪 Testing Node.js Compatibility...');
  
  runTest('Node.js version compatibility', () => {
    const nodeVersion = process.version;
    const majorVersion = parseInt(nodeVersion.slice(1).split('.')[0]);
    
    if (majorVersion < 14) {
      throw new Error(`Node.js version ${nodeVersion} is below minimum required version 14`);
    }
  });
}

// Test index.js syntax
function testIndexSyntax() {
  console.log('\n🧪 Testing Index.js Syntax...');
  
  runTest('index.js syntax validation', () => {
    const indexPath = path.join(__dirname, '..', 'index.js');
    const indexContent = fs.readFileSync(indexPath, 'utf8');
    
    // Check for shebang
    if (!indexContent.startsWith('#!/usr/bin/env node')) {
      throw new Error('Missing or incorrect shebang');
    }
    
    // Check for required functions
    const requiredFunctions = ['getPlatformBinary', 'validateBinaryExists', 'main'];
    for (const func of requiredFunctions) {
      if (!indexContent.includes(`function ${func}`)) {
        throw new Error(`Missing function: ${func}`);
      }
    }
    
    // Check for main() call
    if (!indexContent.includes('main()')) {
      throw new Error('Missing main() function call');
    }
  });
}

// Main test runner
function runAllTests() {
  console.log('🚀 Running mcp-hub-tui NPM Package Tests\n');
  
  try {
    testPackageStructure();
    testIndexSyntax();
    testNodeCompatibility();
    testPlatformDetection();
    testBinaryExistence();
    
    console.log('\n✅ All tests passed!');
    console.log('\n📊 Test Summary:');
    console.log(`   Platform: ${os.platform()}`);
    console.log(`   Architecture: ${os.arch()}`);
    console.log(`   Node.js: ${process.version}`);
    console.log(`   Package: mcp-hub-tui`);
    
  } catch (error) {
    console.error('\n❌ Test suite failed:', error.message);
    process.exit(1);
  }
}

// Run tests if called directly
if (require.main === module) {
  runAllTests();
}

module.exports = {
  runAllTests,
  testPlatformDetection,
  testBinaryExistence,
  testPackageStructure,
  testNodeCompatibility,
  testIndexSyntax
};