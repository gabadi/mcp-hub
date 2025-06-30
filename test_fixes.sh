#!/bin/bash

echo "Building cc-mcp-manager..."
go build -o cc-mcp-manager ./main.go

if [ $? -eq 0 ]; then
    echo "Build successful!"
    echo ""
    echo "Starting cc-mcp-manager to test the fixes..."
    echo ""
    echo "Test the following fixes:"
    echo "1. Search Input Visibility:"
    echo "   - Press '/' or Tab to enter search mode"
    echo "   - Type some text - it should be visible in a gray box with white text"
    echo "   - The cursor should be a solid block (█) when in input mode"
    echo ""
    echo "2. Grid Selection with Search:"
    echo "   - Enter search mode and type 'git' or 'docker'"
    echo "   - Use arrow keys to navigate filtered results"
    echo "   - The selection should properly highlight filtered items"
    echo ""
    echo "3. Rendering:"
    echo "   - Press Ctrl+L to clear screen if you see artifacts"
    echo "   - The display should maintain consistent height"
    echo ""
    echo "4. Modal Windows:"
    echo "   - Press 'a' to see Add modal"
    echo "   - Press 'e' to see Edit modal"
    echo "   - Press 'd' to see Delete modal"
    echo "   - Press ESC to close modals"
    echo ""
    echo "Press Ctrl+C to exit the application"
    echo ""
    
    ./cc-mcp-manager
else
    echo "Build failed!"
    exit 1
fi