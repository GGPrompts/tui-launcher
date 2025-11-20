#!/bin/bash
# TUI Launcher Installation Script

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}╔════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║     TUI Launcher Installation         ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════╝${NC}"
echo ""

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# Step 1: Build
echo -e "${YELLOW}→${NC} Building tui-launcher..."
go build -o tui-launcher
if [ $? -ne 0 ]; then
    echo -e "${RED}✗${NC} Build failed"
    exit 1
fi
echo -e "${GREEN}✓${NC} Build successful"

# Step 2: Create bin directory
BIN_DIR="$HOME/.local/bin"
echo -e "${YELLOW}→${NC} Creating $BIN_DIR..."
mkdir -p "$BIN_DIR"
echo -e "${GREEN}✓${NC} Directory ready"

# Step 3: Install binary
echo -e "${YELLOW}→${NC} Installing binary to $BIN_DIR/tui-launcher..."
cp tui-launcher "$BIN_DIR/tui-launcher"
chmod +x "$BIN_DIR/tui-launcher"
echo -e "${GREEN}✓${NC} Binary installed"

# Step 4: Create wrapper script
echo -e "${YELLOW}→${NC} Creating 'tl' wrapper..."
cat > "$BIN_DIR/tl" << 'EOF'
#!/bin/bash
# TUI Launcher wrapper - run from any directory
~/.local/bin/tui-launcher "$@"
EOF
chmod +x "$BIN_DIR/tl"
echo -e "${GREEN}✓${NC} Wrapper created"

# Step 5: Create config directory
CONFIG_DIR="$HOME/.config/tui-launcher"
echo -e "${YELLOW}→${NC} Setting up config directory..."
mkdir -p "$CONFIG_DIR"

# Copy sample config if it doesn't exist
if [ ! -f "$CONFIG_DIR/config.yaml" ]; then
    if [ -f "$HOME/.config/tui-launcher/config.yaml" ]; then
        echo -e "${GREEN}✓${NC} Config already exists"
    else
        echo -e "${YELLOW}→${NC} Config will be created on first run"
    fi
else
    echo -e "${GREEN}✓${NC} Config already exists"
fi

echo ""
echo -e "${GREEN}╔════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║     Installation Complete! 🎉          ║${NC}"
echo -e "${GREEN}╚════════════════════════════════════════╝${NC}"
echo ""

# Check if ~/.local/bin is in PATH
if [[ ":$PATH:" != *":$HOME/.local/bin:"* ]]; then
    echo -e "${YELLOW}⚠${NC}  Add ~/.local/bin to your PATH:"
    echo ""
    echo -e "   ${BLUE}# Add to ~/.bashrc or ~/.zshrc:${NC}"
    echo -e "   export PATH=\"\$HOME/.local/bin:\$PATH\""
    echo ""
    echo -e "   ${BLUE}# Then reload:${NC}"
    echo -e "   source ~/.bashrc  ${BLUE}# or${NC} source ~/.zshrc"
    echo ""
else
    echo -e "${GREEN}✓${NC} ~/.local/bin is already in PATH"
    echo ""
fi

echo -e "${GREEN}Ready to use!${NC}"
echo ""
echo -e "${BLUE}Usage:${NC}"
echo -e "  tl                    ${BLUE}# Launch TUI${NC}"
echo -e "  tui-launcher          ${BLUE}# Full command${NC}"
echo ""
echo -e "${BLUE}Keybindings:${NC}"
echo -e "  ↑/↓ or j/k            ${BLUE}# Navigate${NC}"
echo -e "  Space                 ${BLUE}# Expand category OR select command${NC}"
echo -e "  Enter                 ${BLUE}# Launch selected${NC}"
echo -e "  c                     ${BLUE}# Clear selections${NC}"
echo -e "  t                     ${BLUE}# Toggle tmux mode${NC}"
echo -e "  q                     ${BLUE}# Quit${NC}"
echo ""
echo -e "${BLUE}Config:${NC} ~/.config/tui-launcher/config.yaml"
echo ""
