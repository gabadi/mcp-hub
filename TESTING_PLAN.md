# Cmd+V Paste Functionality Testing Plan

## Pre-Testing Setup

### 1. Grant macOS Permissions
- [ ] System Settings → Privacy & Security → Accessibility → Add Warp
- [ ] System Settings → Privacy & Security → Automation → Enable Warp → System Events
- [ ] Restart Warp Terminal after granting permissions

### 2. Prepare Test Content
Copy this text to clipboard before testing:
```
test clipboard content with spaces
```

## Testing Sequence

### Phase 1: Enhanced Debug Tool
```bash
./debug-tool
```

**Expected Results:**
- Shows system info (macOS, WarpTerminal, versions)
- Detects Cmd+V key press
- Shows clipboard content when Cmd+V is pressed
- No clipboard errors

**Test Cases:**
1. Press Cmd+V → Should show clipboard content
2. Press Ctrl+V → Should show clipboard content  
3. Press other keys → Should register normally

### Phase 2: Main Application Testing
```bash
./mcp-manager-test
```

**Search Mode Tests:**
1. Enter search mode (press `/`)
2. Press Cmd+V → Should paste clipboard content to search
3. Verify success message appears
4. Verify search query contains pasted content
5. Test with existing text + paste
6. Test with multi-line clipboard content (should be cleaned)

**Modal Form Tests:**
1. Add MCP → Enter form
2. Press Cmd+V in each field → Should paste content
3. Verify existing paste functionality still works

### Phase 3: Error Condition Testing

**Permission Denied Test:**
1. Revoke Warp accessibility permissions
2. Try Cmd+V in search mode
3. Should show error message: "Failed to paste from clipboard: ..."

**Empty Clipboard Test:**
1. Clear clipboard: `pbcopy < /dev/null`
2. Try Cmd+V in search mode
3. Should paste empty content (no error)

**Large Content Test:**
1. Copy large text (>1000 characters) with newlines
2. Paste in search mode
3. Should clean and paste content properly

## Success Criteria

### Must Pass:
- [ ] Cmd+V detected in debug tool
- [ ] Clipboard content shown in debug tool
- [ ] Cmd+V pastes in search mode
- [ ] Success message appears after paste
- [ ] Content properly cleaned (no newlines)
- [ ] Existing modal paste functionality unchanged

### Should Pass:
- [ ] Error messages for permission issues
- [ ] Ctrl+V works as alternative
- [ ] Multi-line content cleaned properly
- [ ] Large content handles gracefully

### Known Limitations:
- [ ] WarpTerminal Copy-on-Select may not work in TUI mode
- [ ] Some clipboard managers may interfere
- [ ] macOS permission dialogs may appear on first use

## Debugging Failed Tests

### If Cmd+V not detected:
1. Check `key_debug.log` for actual key strings
2. Verify terminal environment variables
3. Test in different terminal (Terminal.app)

### If Clipboard access fails:
1. Check macOS Console.app for permission errors
2. Verify osascript works: `osascript -e 'get the clipboard'`
3. Reset clipboard daemon: `sudo killall pboard`

### If Paste doesn't work despite detection:
1. Check search.go implementation
2. Verify imports are correct
3. Look for compilation errors

## Regression Testing

Ensure these still work after changes:
- [ ] Modal form paste functionality
- [ ] Search mode typing
- [ ] Search mode backspace
- [ ] All other keyboard shortcuts
- [ ] Application exit (Ctrl+C, Esc)

## Performance Testing

- [ ] Paste response time < 100ms
- [ ] No memory leaks with repeated paste
- [ ] Large clipboard content doesn't freeze app