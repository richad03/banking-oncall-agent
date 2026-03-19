# Banking On-Call Agent - MCP Integration

Use your Banking On-Call Agent directly in Claude conversations via **Model Context Protocol (MCP)**!

## 🚀 **What is MCP Integration?**

Instead of using this as a Slack bot, you can connect it directly to **Claude Code** or other MCP clients to get banking incident expertise right in your Claude conversations.

### **Benefits:**
- 🤖 **Direct Claude Access** - Ask Claude about banking incidents directly
- 📚 **Instant Runbook Search** - No need to leave your conversation
- 🔍 **Error Code Analysis** - Get immediate guidance on PT-401, DEP-500, etc.
- 🏗️ **Service Dependencies** - Understand impact scope instantly
- ⚡ **Real-Time Knowledge** - Always up-to-date with your banking-oncall docs

## 🛠️ **Quick Setup**

### **1. Build the MCP Server**
```bash
# Build MCP server
go build -o bin/banking-mcp-server ./cmd/mcp-server

# Or use make
make mcp-build
```

### **2. Configure Environment**
```bash
# Set knowledge base path
export KNOWLEDGE_BASE_PATH="../banking-oncall"

# Or create .env file
echo "KNOWLEDGE_BASE_PATH=../banking-oncall" > .env
```

### **3. Add to Claude Code**
Add this to your Claude Code MCP configuration (`~/Library/Application Support/Claude/claude_desktop_config.json`):

```json
{
  "mcpServers": {
    "banking-oncall-agent": {
      "command": "go",
      "args": ["run", "./cmd/mcp-server/main.go"],
      "cwd": "/path/to/banking-oncall-agent",
      "env": {
        "KNOWLEDGE_BASE_PATH": "../banking-oncall"
      }
    }
  }
}
```

Or use the pre-built binary:
```json
{
  "mcpServers": {
    "banking-oncall-agent": {
      "command": "/path/to/banking-oncall-agent/bin/banking-mcp-server",
      "cwd": "/path/to/banking-oncall-agent",
      "env": {
        "KNOWLEDGE_BASE_PATH": "../banking-oncall"
      }
    }
  }
}
```

### **4. Restart Claude Code**
Restart Claude Code to load the new MCP server.

## 🎯 **Available MCP Tools**

Once connected, Claude will have access to these banking tools:

### **1. `analyze_banking_incident`**
Analyze a banking service incident and get comprehensive guidance.

**Parameters:**
- `description` (required) - Description of the incident
- `error_code` (optional) - Specific error code (PT-401, DEP-500, etc.)
- `service` (optional) - Banking service name

**Example:**
> *"I'm seeing PT-401 errors in payroll processing"*

### **2. `search_banking_runbooks`**
Search runbooks and documentation for procedures.

**Parameters:**
- `query` (required) - Search terms
- `service` (optional) - Limit to specific service

**Example:**
> *"Search for tax calculation procedures"*

### **3. `get_service_dependencies`**
Get dependency information for impact analysis.

**Parameters:**
- `service` (required) - Service name

**Example:**
> *"What are the dependencies for payroll service?"*

### **4. `list_banking_services`**
Get overview of all banking services and coverage.

**No parameters required**

## 💬 **Usage Examples**

Once configured, you can ask Claude directly:

### **Incident Analysis**
```
Claude, I'm seeing timeout errors in the deposits service.
Can you analyze this incident and tell me what to do?
```

Claude will automatically:
1. Use `analyze_banking_incident` tool
2. Detect this is an outflows/deposits issue
3. Find relevant runbook sections
4. Provide immediate actions and admin links

### **Error Code Lookup**
```
What should I do about PT-401 error?
```

Claude will:
1. Search for PT-401 in runbooks
2. Provide specific resolution steps
3. Show related error codes and procedures

### **Service Impact Analysis**
```
If payroll service goes down, what other services are affected?
```

Claude will:
1. Get payroll service dependencies
2. Show upstream and downstream impact
3. Identify critical path services

### **Runbook Search**
```
How do I handle billpay timeouts?
```

Claude will:
1. Search runbooks for billpay timeout procedures
2. Return relevant sections with steps
3. Provide admin console links

## 🔧 **Advanced Configuration**

### **Custom Knowledge Base Path**
```json
{
  "mcpServers": {
    "banking-oncall-agent": {
      "command": "/path/to/bin/banking-mcp-server",
      "env": {
        "KNOWLEDGE_BASE_PATH": "/custom/path/to/docs"
      }
    }
  }
}
```

### **Multiple Environments**
```json
{
  "mcpServers": {
    "banking-oncall-prod": {
      "command": "/path/to/bin/banking-mcp-server",
      "env": {
        "KNOWLEDGE_BASE_PATH": "/prod/banking-oncall"
      }
    },
    "banking-oncall-staging": {
      "command": "/path/to/bin/banking-mcp-server",
      "env": {
        "KNOWLEDGE_BASE_PATH": "/staging/banking-oncall"
      }
    }
  }
}
```

## 🚀 **Development & Testing**

### **Test MCP Server**
```bash
# Run MCP server manually (for testing)
./bin/banking-mcp-server

# Test with MCP client tools
echo '{"jsonrpc":"2.0","id":1,"method":"tools/list"}' | ./bin/banking-mcp-server
```

### **Debug Mode**
Enable logging by setting environment variable:
```bash
export MCP_DEBUG=1
./bin/banking-mcp-server
```

### **Hot Reload During Development**
```bash
# Use go run for auto-rebuild
go run ./cmd/mcp-server/main.go
```

## 🎯 **Benefits vs Slack Bot**

| Feature | Slack Bot | MCP Integration |
|---------|-----------|-----------------|
| **Access** | Need Slack workspace | Direct in Claude conversations |
| **Context** | Separate from work | Integrated with your thinking |
| **Speed** | Wait for bot response | Instant Claude integration |
| **Privacy** | Shared in channels | Private conversations |
| **Workflow** | Context switching | Seamless integration |

## 🛠️ **Troubleshooting**

### **MCP Server Won't Start**
```bash
# Check if knowledge base exists
ls $KNOWLEDGE_BASE_PATH

# Verify Go build works
go build -o bin/test ./cmd/mcp-server && ./bin/test
```

### **Claude Code Can't Connect**
1. Check Claude Code configuration syntax
2. Verify file paths are absolute
3. Restart Claude Code after changes
4. Check Claude Code logs for errors

### **No Banking Tools Available**
1. Verify MCP server is listed in Claude Code settings
2. Check that knowledge base path contains markdown files
3. Ensure banking-oncall repository is accessible

### **Tools Return Empty Results**
1. Verify banking-oncall docs are properly formatted
2. Check knowledge base loading logs
3. Test with simple service queries first

## 📚 **Next Steps**

1. **Set up MCP server** with your banking-oncall repository
2. **Configure Claude Code** with MCP connection
3. **Test basic queries** to verify functionality
4. **Train your team** on MCP usage patterns
5. **Customize** for your specific runbook structure

---

**Ready to revolutionize your banking incident response directly in Claude!** 🏦🤖