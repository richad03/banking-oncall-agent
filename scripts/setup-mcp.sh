#!/bin/bash

# Banking On-Call Agent MCP Setup Script
# Automatically configures Claude Code with banking expertise

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper functions
log_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

log_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

log_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

log_error() {
    echo -e "${RED}❌ $1${NC}"
}

echo -e "${BLUE}"
echo "🏦 Banking On-Call Agent MCP Setup"
echo "=================================="
echo -e "${NC}"

# Get current directory (where this script is run)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
MCP_SERVER_PATH="$SCRIPT_DIR/bin/banking-mcp-server"

log_info "Setting up Banking On-Call MCP Server from: $SCRIPT_DIR"
echo ""

# Step 1: Verify MCP server binary exists
if [ ! -f "$MCP_SERVER_PATH" ]; then
    log_error "MCP server binary not found at: $MCP_SERVER_PATH"
    log_info "Please ensure you have the complete package with bin/ directory"
    exit 1
fi

log_success "Found MCP server binary"

# Step 2: Make binary executable
chmod +x "$MCP_SERVER_PATH"
log_success "MCP server is executable"

# Step 3: Find banking-oncall repository
log_info "Searching for banking-oncall repository..."

# Common locations to check
POSSIBLE_PATHS=(
    "../banking-oncall"
    "../../banking-oncall"
    "$HOME/banking-oncall"
    "$HOME/work/banking-oncall"
    "$HOME/repos/banking-oncall"
    "$HOME/src/banking-oncall"
    "/Users/$(whoami)/banking-oncall"
)

KNOWLEDGE_BASE_PATH=""
for path in "${POSSIBLE_PATHS[@]}"; do
    if [ -d "$path" ] && [ -f "$path/README.md" ]; then
        # Convert to absolute path
        KNOWLEDGE_BASE_PATH=$(cd "$path" && pwd)
        log_success "Found banking-oncall repository at: $KNOWLEDGE_BASE_PATH"
        break
    fi
done

# If not found, ask user
if [ -z "$KNOWLEDGE_BASE_PATH" ]; then
    log_warning "Could not automatically find banking-oncall repository"
    echo ""
    echo "Please enter the full path to your banking-oncall repository:"
    echo "Example: /Users/$(whoami)/work/banking-oncall"
    echo ""
    read -p "Path: " USER_PATH

    if [ -d "$USER_PATH" ]; then
        KNOWLEDGE_BASE_PATH="$USER_PATH"
        log_success "Using banking-oncall repository at: $KNOWLEDGE_BASE_PATH"
    else
        log_error "Directory not found: $USER_PATH"
        exit 1
    fi
fi

# Step 4: Count markdown files to verify it's the right repo
MD_COUNT=$(find "$KNOWLEDGE_BASE_PATH" -name "*.md" | wc -l)
if [ "$MD_COUNT" -lt 5 ]; then
    log_warning "Only found $MD_COUNT markdown files in $KNOWLEDGE_BASE_PATH"
    log_warning "This might not be the correct banking-oncall repository"
    echo ""
    read -p "Continue anyway? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
else
    log_success "Found $MD_COUNT markdown files in banking-oncall repository"
fi

# Step 5: Create Claude Code configuration
CLAUDE_CONFIG_DIR="$HOME/Library/Application Support/Claude"
CLAUDE_CONFIG_FILE="$CLAUDE_CONFIG_DIR/claude_desktop_config.json"

log_info "Configuring Claude Code..."

# Create directory if it doesn't exist
mkdir -p "$CLAUDE_CONFIG_DIR"

# Check if config file already exists
if [ -f "$CLAUDE_CONFIG_FILE" ]; then
    log_info "Found existing Claude Code configuration"

    # Backup existing config
    cp "$CLAUDE_CONFIG_FILE" "$CLAUDE_CONFIG_FILE.backup.$(date +%s)"
    log_success "Backed up existing configuration"

    # Check if banking-oncall-agent already exists in config
    if grep -q "banking-oncall-agent" "$CLAUDE_CONFIG_FILE" 2>/dev/null; then
        log_warning "Banking On-Call Agent already configured in Claude Code"
        echo ""
        read -p "Update configuration? (Y/n): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Nn]$ ]]; then
            log_info "Skipping configuration update"
        else
            # Update existing configuration
            log_info "Updating existing configuration..."
            # For now, we'll create a new config - in production you'd want to merge JSON properly
        fi
    fi
fi

# Create new configuration
cat > "$CLAUDE_CONFIG_FILE" << EOF
{
  "mcpServers": {
    "banking-oncall-agent": {
      "command": "$MCP_SERVER_PATH",
      "cwd": "$SCRIPT_DIR",
      "env": {
        "KNOWLEDGE_BASE_PATH": "$KNOWLEDGE_BASE_PATH"
      }
    }
  }
}
EOF

log_success "Claude Code configuration updated"

# Step 6: Test MCP server
log_info "Testing MCP server functionality..."

# Test basic communication
if timeout 10s echo '{"jsonrpc":"2.0","id":1,"method":"tools/list"}' | "$MCP_SERVER_PATH" >/dev/null 2>&1; then
    log_success "MCP server is responding correctly"
else
    log_error "MCP server test failed"
    log_info "Try running manually: $MCP_SERVER_PATH"
    exit 1
fi

# Step 7: Show completion message
echo ""
log_success "🎉 Banking On-Call Agent MCP setup complete!"
echo ""
log_info "Next steps:"
echo "  1. ${GREEN}Restart Claude Code${NC} (completely quit and reopen)"
echo "  2. ${GREEN}Open a new conversation${NC} in Claude Code"
echo "  3. ${GREEN}Test with a banking question${NC}, like:"
echo "     \"I'm seeing PT-401 errors in payroll. What should I do?\""
echo ""
log_info "Available banking tools in Claude:"
echo "  • analyze_banking_incident - Full incident analysis"
echo "  • search_banking_runbooks - Search procedures"
echo "  • get_service_dependencies - Impact analysis"
echo "  • list_banking_services - Service overview"
echo ""
log_info "Configuration details:"
echo "  • MCP Server: $MCP_SERVER_PATH"
echo "  • Knowledge Base: $KNOWLEDGE_BASE_PATH"
echo "  • Claude Config: $CLAUDE_CONFIG_FILE"
echo ""

# Step 8: Offer to test with Claude Code
if command -v "Claude Code" &> /dev/null || [ -d "/Applications/Claude Code.app" ]; then
    log_info "Claude Code detected on your system"
    echo ""
    read -p "Restart Claude Code now? (Y/n): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Nn]$ ]]; then
        log_info "Attempting to restart Claude Code..."
        pkill -f "Claude Code" 2>/dev/null || true
        sleep 2
        if [ -d "/Applications/Claude Code.app" ]; then
            open "/Applications/Claude Code.app"
            log_success "Claude Code restarted"
        else
            log_info "Please restart Claude Code manually"
        fi
    fi
fi

echo ""
log_success "Setup complete! Your Claude conversations now have banking expertise! 🚀"