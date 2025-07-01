# WarpTerminal Cmd+V Clipboard Fix Guide

## Issue
Cmd+V is detected but paste functionality doesn't work in cc-mcp-manager when using WarpTerminal on macOS.

## Root Cause
1. **Missing paste handler in search mode** - The search input handler didn't implement paste functionality
2. **macOS clipboard permissions** - Terminal apps need accessibility permissions for clipboard access
3. **WarpTerminal TUI limitations** - Known issue with clipboard operations in TUI applications

## Solutions

### 1. Code Fix (Implemented)
The search mode now handles paste operations:
```go
case "ctrl+v", "cmd+v", "⌘v", "command+v":
    // Paste clipboard content to search query
    model = pasteToSearchQuery(model)
```

### 2. macOS Permissions Setup

#### Grant Terminal Accessibility Permissions:
1. Open **System Settings** → **Privacy & Security**
2. Click **Accessibility** in the sidebar
3. Click the **+** button and add **Warp**
4. Ensure **Warp** is checked/enabled

#### Grant Automation Permissions:
1. Open **System Settings** → **Privacy & Security**
2. Click **Automation** in the sidebar
3. Find **Warp** and expand it
4. Enable **System Events** if present

#### Alternative: Grant Full Disk Access (if needed):
1. Open **System Settings** → **Privacy & Security**
2. Click **Full Disk Access** in the sidebar
3. Click **+** button and add **Warp**
4. Enable **Warp**

### 3. WarpTerminal Specific Fixes

#### Check Warp Settings:
1. Open Warp settings (Cmd+,)
2. Navigate to **Features** → **Terminal**
3. Ensure clipboard features are enabled

#### Alternative Paste Methods in Warp:
- Use **Cmd+Shift+V** for plain text paste
- Use **Edit** menu → **Paste** option
- Use right-click context menu paste

### 4. Testing & Debugging

#### Test with Enhanced Debug Tool:
```bash
cd /Users/gabadi/workspace/melech/cc-mcp-manager
go run cmd/debug/main.go
```

This will show:
- Key detection status
- Clipboard access errors
- System environment info
- Real-time paste testing

#### Manual Clipboard Test:
```bash
# Test if clipboard works from terminal
echo "test" | pbcopy
pbpaste
```

### 5. Alternative Solutions

#### If permissions don't work:
1. **Restart Warp** after granting permissions
2. **Restart macOS** to reset permission caches
3. **Reset permissions** with: `sudo tccutil reset Accessibility`
4. **Try different terminal** (Terminal.app, iTerm2) to isolate issue

#### Fallback keyboard shortcuts:
- **Ctrl+V** should always work (cross-platform compatibility)
- **Cmd+Shift+V** for plain text in some terminals

## Expected Results

After implementing these fixes:
1. **Cmd+V should work** in search mode with visual feedback
2. **Error messages** will appear if clipboard access fails
3. **Success messages** will confirm successful paste operations
4. **Debug tool** will show clipboard content when paste keys are pressed

## Troubleshooting

### If Cmd+V still doesn't work:
1. Check debug logs for clipboard errors
2. Verify permissions are actually granted (lock icon unlocked)
3. Try other terminal applications to isolate Warp-specific issues
4. Use `sudo killall pboard` to reset clipboard daemon

### Error Messages:
- "Failed to paste from clipboard: ..." = Permission or clipboard access issue
- No feedback at all = Key not being detected (check debug tool)
- Success message but no paste = Implementation bug (check code)