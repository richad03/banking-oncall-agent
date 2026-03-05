#!/bin/bash

# Test script for Banking On-Call Agent MCP Server
# This script tests basic MCP functionality

set -e

echo "🧪 Testing Banking On-Call Agent MCP Server"
echo "=========================================="

# Build if needed
if [ ! -f "bin/banking-mcp-server" ]; then
    echo "🔨 Building MCP server..."
    go build -o bin/banking-mcp-server ./cmd/mcp-server
fi

echo ""
echo "🚀 Testing MCP Server Communications"
echo ""

# Test 1: Initialize
echo "1️⃣ Testing initialization..."
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05"}}' | ./bin/banking-mcp-server | head -1

echo ""

# Test 2: List tools
echo "2️⃣ Testing tools list..."
echo '{"jsonrpc":"2.0","id":2,"method":"tools/list"}' | ./bin/banking-mcp-server | head -1

echo ""

# Test 3: Test incident analysis
echo "3️⃣ Testing incident analysis..."
cat << 'EOF' | ./bin/banking-mcp-server | head -1
{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"analyze_banking_incident","arguments":{"description":"payroll error PT-401","error_code":"PT-401"}}}
EOF

echo ""

# Test 4: Test runbook search
echo "4️⃣ Testing runbook search..."
cat << 'EOF' | ./bin/banking-mcp-server | head -1
{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"search_banking_runbooks","arguments":{"query":"tax processing"}}}
EOF

echo ""

# Test 5: Test service list
echo "5️⃣ Testing service overview..."
cat << 'EOF' | ./bin/banking-mcp-server | head -1
{"jsonrpc":"2.0","id":5,"method":"tools/call","params":{"name":"list_banking_services","arguments":{}}}
EOF

echo ""
echo "✅ MCP Server testing completed!"
echo ""
echo "🎯 Next steps:"
echo "   1. Add to Claude Code configuration"
echo "   2. Restart Claude Code"
echo "   3. Start using banking tools in Claude conversations"
echo ""
echo "📖 See MCP-README.md for detailed setup instructions"