#!/bin/bash

# Install All MCP Servers for Banking Operations
# PagerDuty, Jira, Datadog, Sentry, Presidio + Banking On-Call

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

log_info() { echo -e "${BLUE}ℹ️  $1${NC}"; }
log_success() { echo -e "${GREEN}✅ $1${NC}"; }
log_warning() { echo -e "${YELLOW}⚠️  $1${NC}"; }
log_error() { echo -e "${RED}❌ $1${NC}"; }

echo -e "${BLUE}"
echo "🚀 Installing Complete MCP Stack for Banking Operations"
echo "========================================================"
echo -e "${NC}"

# Step 1: Check prerequisites
log_info "Checking prerequisites..."

# Check Node.js
if ! command -v node &> /dev/null; then
    log_error "Node.js is required but not installed"
    log_info "Install Node.js: https://nodejs.org/"
    exit 1
fi

NODE_VERSION=$(node --version | grep -oE '[0-9]+' | head -1)
if [ "$NODE_VERSION" -lt 18 ]; then
    log_warning "Node.js version $NODE_VERSION detected. Recommend 18+ for best compatibility"
fi

log_success "Node.js $(node --version) found"

# Check npm/npx
if ! command -v npx &> /dev/null; then
    log_error "npx is required but not found"
    exit 1
fi

log_success "npm/npx found"

# Step 2: Test MCP server installations
log_info "Testing MCP server availability..."

echo ""
log_info "Testing PagerDuty MCP server..."
if npx -y @pagerduty/mcp-server --version &>/dev/null; then
    log_success "PagerDuty MCP server available"
else
    log_warning "PagerDuty MCP server installation may require manual setup"
fi

echo ""
log_info "Testing Atlassian/Jira MCP server..."
if npx -y @sooperset/mcp-atlassian --help &>/dev/null; then
    log_success "Atlassian MCP server available"
else
    log_warning "Atlassian MCP server installation may require manual setup"
fi

echo ""
log_info "Testing Datadog MCP server..."
if npx -y mcp-server-datadog --help &>/dev/null; then
    log_success "Datadog MCP server available"
else
    log_warning "Datadog MCP server installation may require manual setup"
fi

echo ""
log_info "Testing Sentry MCP server..."
if npx -y @sentry/mcp-server --help &>/dev/null; then
    log_success "Sentry MCP server available"
else
    log_warning "Sentry MCP server installation may require manual setup"
fi

echo ""
log_info "Testing Presidio MCP server..."
if npx -y @presidio-oss/specifai-mcp-server --help &>/dev/null; then
    log_success "Presidio MCP server available"
else
    log_warning "Presidio MCP server installation may require manual setup"
fi

# Step 3: Create API token setup guide
log_info "Creating API token setup guide..."

cat > "MCP_TOKENS_SETUP.md" << 'EOF'
# MCP Servers API Token Setup Guide

## 🔑 **Required API Tokens**

### **1. PagerDuty**
```bash
# Get API token from PagerDuty:
# 1. Go to https://your-company.pagerduty.com/api_keys
# 2. Create new API key with "Read/Write" access
# 3. Copy the API token

PAGERDUTY_API_TOKEN="your-pagerduty-api-token"
PAGERDUTY_USER_EMAIL="your-email@company.com"
```

### **2. Atlassian/Jira**
```bash
# Get API token from Atlassian:
# 1. Go to https://id.atlassian.com/manage-profile/security/api-tokens
# 2. Create new API token
# 3. Copy the token

ATLASSIAN_URL="https://your-company.atlassian.net"
ATLASSIAN_EMAIL="your-email@company.com"
ATLASSIAN_API_TOKEN="your-atlassian-api-token"
```

### **3. Datadog**
```bash
# Get API keys from Datadog:
# 1. Go to https://app.datadoghq.com/organization-settings/api-keys
# 2. Create new API key and App key
# 3. Copy both keys

DATADOG_API_KEY="your-datadog-api-key"
DATADOG_APP_KEY="your-datadog-app-key"
DATADOG_SITE="datadoghq.com"  # or datadoghq.eu, us3.datadoghq.com, etc.
```

### **4. Sentry**
```bash
# Get auth token from Sentry:
# 1. Go to https://sentry.io/settings/account/api/auth-tokens/
# 2. Create new auth token with appropriate permissions
# 3. Copy the token

SENTRY_AUTH_TOKEN="your-sentry-auth-token"
SENTRY_ORG="your-sentry-org-slug"
```

### **5. Presidio (Optional)**
```bash
# Presidio may not require API key for basic usage
PRESIDIO_API_KEY="optional-presidio-api-key"
```

## 🔧 **Environment Setup**

Create a `.env` file with all your tokens:

```bash
# PagerDuty
PAGERDUTY_API_TOKEN=your-token-here
PAGERDUTY_USER_EMAIL=your-email@company.com

# Atlassian/Jira
ATLASSIAN_URL=https://your-company.atlassian.net
ATLASSIAN_EMAIL=your-email@company.com
ATLASSIAN_API_TOKEN=your-token-here

# Datadog
DATADOG_API_KEY=your-api-key-here
DATADOG_APP_KEY=your-app-key-here
DATADOG_SITE=datadoghq.com

# Sentry
SENTRY_AUTH_TOKEN=your-token-here
SENTRY_ORG=your-org-slug

# Presidio (optional)
PRESIDIO_API_KEY=optional-key

# Banking On-Call
KNOWLEDGE_BASE_PATH=/path/to/banking-oncall
```

## 🚀 **Next Steps**

1. Get all required API tokens using the guides above
2. Update the Claude Code configuration with your tokens
3. Restart Claude Code
4. Test each MCP server with Claude!

EOF

log_success "Created MCP_TOKENS_SETUP.md"

# Step 4: Create Claude Code configuration
CLAUDE_CONFIG_DIR="$HOME/Library/Application Support/Claude"
CLAUDE_CONFIG_FILE="$CLAUDE_CONFIG_DIR/claude_desktop_config.json"

log_info "Updating Claude Code configuration..."

# Backup existing config if it exists
if [ -f "$CLAUDE_CONFIG_FILE" ]; then
    cp "$CLAUDE_CONFIG_FILE" "$CLAUDE_CONFIG_FILE.backup.$(date +%s)"
    log_success "Backed up existing Claude Code configuration"
fi

# Install the new configuration
mkdir -p "$CLAUDE_CONFIG_DIR"
cp "claude_mcp_config_full.json" "$CLAUDE_CONFIG_FILE"
log_success "Updated Claude Code configuration with all MCP servers"

echo ""
log_success "🎉 MCP Installation Complete!"
echo ""
log_warning "⚠️  IMPORTANT NEXT STEPS:"
echo "  1. Edit MCP_TOKENS_SETUP.md to get your API tokens"
echo "  2. Update $CLAUDE_CONFIG_FILE with your actual tokens"
echo "  3. Replace 'your-*-token-here' placeholders with real values"
echo "  4. Restart Claude Code completely"
echo "  5. Test each service!"
echo ""
log_info "🧪 Test Commands in Claude:"
echo "  • 'Show my PagerDuty incidents'"
echo "  • 'List my Jira issues'"
echo "  • 'Show Datadog dashboard'"
echo "  • 'Check Sentry errors'"
echo "  • 'Analyze this text for PII' (Presidio)"
echo "  • 'I'm seeing PT-401 errors in payroll' (Banking)"
echo ""
log_info "📝 Configuration files created:"
echo "  • $(pwd)/MCP_TOKENS_SETUP.md"
echo "  • $CLAUDE_CONFIG_FILE"
echo ""
log_success "Ready to supercharge your Claude with banking operations tools! 🚀"