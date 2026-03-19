# Installation Guide - Banking On-Call MCP

## 🚀 **Super Quick Setup (Recommended)**

```bash
# 1. Download and extract this package
# 2. Run the setup script
./setup-mcp.sh

# 3. Restart Claude Code
# 4. Done! ✨
```

## 📋 **What You Get**

After setup, Claude will have **4 banking tools**:

### **analyze_banking_incident**
Ask Claude: *"I'm seeing PT-401 errors in payroll processing"*

Claude Response:
```
Let me analyze this banking incident for you.

🚨 **Payroll Incident - High Severity**
**Error Code:** PT-401 - Tax Calculation Failure

**Immediate Actions:**
1. Check reckoner worker status in go/payrolladmin
2. Restart worker if queue is backed up
3. Verify tax calculation endpoints
4. Check for recent tax rate updates

**Runbook:** Section 23.4 - Tax Processing Errors
**Dependencies:** Affects employee payments and tax filing
```

### **search_banking_runbooks**
Ask Claude: *"Show me billpay timeout procedures"*

### **get_service_dependencies**
Ask Claude: *"What happens if deposits service goes down?"*

### **list_banking_services**
Ask Claude: *"What banking services do we have?"*

## ⚡ **Manual Installation**

If you prefer manual setup:

### **1. Find Your Paths**
```bash
# MCP Server location
/path/to/banking-oncall-mcp-release/bin/banking-mcp-server

# Your banking-oncall repository
/path/to/your/banking-oncall
```

### **2. Configure Claude Code**

Edit: `~/Library/Application Support/Claude/claude_desktop_config.json`

```json
{
  "mcpServers": {
    "banking-oncall-agent": {
      "command": "/absolute/path/to/banking-mcp-server",
      "cwd": "/absolute/path/to/banking-oncall-mcp-release",
      "env": {
        "KNOWLEDGE_BASE_PATH": "/absolute/path/to/banking-oncall"
      }
    }
  }
}
```

**Important:** Use absolute paths (starting with `/`)

### **3. Restart Claude Code**
Completely quit and reopen Claude Code.

### **4. Test**
Ask Claude: *"What banking services are available?"*

## 🐛 **Troubleshooting**

### **Claude doesn't see banking tools**
1. Check configuration syntax with JSON validator
2. Verify all paths are absolute and correct
3. Restart Claude Code completely
4. Check Console.app for Claude Code error logs

### **Tools return empty results**
1. Verify `KNOWLEDGE_BASE_PATH` points to correct repository
2. Check that banking-oncall contains .md files
3. Run `./test-mcp.sh` to debug

### **Permission denied errors**
```bash
chmod +x bin/banking-mcp-server
```

### **Path not found errors**
Use `pwd` to get absolute paths:
```bash
cd /path/to/banking-oncall-mcp-release
echo "MCP Server: $(pwd)/bin/banking-mcp-server"
echo "Knowledge Base: /path/to/banking-oncall"
```

## 📞 **Getting Help**

1. **Run diagnostics:** `./test-mcp.sh`
2. **Check logs:** Console.app → search "Claude"
3. **Verify setup:** Re-run `./setup-mcp.sh`

## 🔄 **Updates**

To update the banking knowledge:
- No restart needed! The MCP server reads files fresh each time
- Update your banking-oncall repository
- Banking expertise automatically updates in Claude

---

**Need help?** Check that your `banking-oncall` repository path is correct and contains markdown files.