# Banking On-Call MCP Distribution Guide

## 🚀 **Quick Team Setup**

### **For Team Members:**

1. **Download** the MCP package (banking-oncall-mcp-v1.0.0.tar.gz)
2. **Extract** anywhere on their machine
3. **Run setup:** `./setup-mcp.sh`
4. **Restart Claude Code**
5. **Test:** Ask Claude "What banking services are available?"

**That's it!** ✨ They now have banking expertise in Claude.

---

## 📦 **Distribution Methods**

### **Method 1: Internal File Share**
```bash
# Share the archive directly
banking-oncall-mcp-v1.0.0.tar.gz (1.9MB)

# Recipients extract and run:
tar -xzf banking-oncall-mcp-v1.0.0.tar.gz
cd banking-oncall-mcp-release
./setup-mcp.sh
```

### **Method 2: GitHub Release**
Create a release on your banking-oncall-agent repository:

```bash
# 1. Push to GitHub
git add .
git commit -m "Add MCP server release package"
git push

# 2. Create GitHub release with banking-oncall-mcp-v1.0.0.tar.gz
# 3. Share release URL with team
```

### **Method 3: Internal Package Repository**
```bash
# Add to your internal package manager
# Team installs with: brew install banking-oncall-mcp
# or: npm install @square/banking-oncall-mcp
```

### **Method 4: Shared Network Drive**
```bash
# Place on shared drive:
/shared/tools/banking-oncall-mcp-release/

# Team members run:
/shared/tools/banking-oncall-mcp-release/setup-mcp.sh
```

---

## 🎯 **What Team Members Get**

After setup, every team member can ask Claude:

**Incident Response:**
- *"I'm seeing PT-401 errors in payroll processing"*
- *"DEP-500 timeout in deposits - what do I do?"*
- *"Billpay service seems down, need runbook"*

**Service Information:**
- *"What are the dependencies for savings service?"*
- *"Show me all error codes for outflows"*
- *"What banking services do we have?"*

**Procedure Lookup:**
- *"How do I restart the reckoner worker?"*
- *"What's the procedure for checkbook timeouts?"*
- *"Show me tax processing steps"*

Claude automatically uses banking tools to provide expert guidance!

---

## 🔧 **Administrative Notes**

### **Requirements for Each User:**
- Claude Code (Desktop app)
- Access to banking-oncall repository
- macOS/Linux (Windows users need WSL)

### **Setup Time Per User:**
- **Automated:** 2 minutes with setup script
- **Manual:** 5 minutes following INSTALL.md

### **Maintenance:**
- **Zero maintenance** - reads banking-oncall docs in real-time
- **Updates automatically** when banking-oncall is updated
- **No server infrastructure** needed

### **Security:**
- **Local only** - no external API calls
- **Uses existing banking-oncall permissions**
- **Private Claude conversations** - not shared

---

## 📊 **Rollout Strategy**

### **Phase 1: Early Adopters (Week 1)**
- Share with 2-3 senior engineers
- Gather feedback and iterate
- Document common use cases

### **Phase 2: On-Call Team (Week 2-3)**
- Roll out to all on-call engineers
- Include in on-call training materials
- Create usage examples and best practices

### **Phase 3: Full Banking Team (Month 1)**
- Share with entire banking engineering team
- Add to onboarding documentation
- Measure impact on incident response times

### **Success Metrics:**
- Reduced mean time to resolution (MTTR)
- Fewer escalations due to missing context
- Improved consistency in incident response
- Higher confidence in on-call procedures

---

## 🚀 **Next Steps**

1. **Choose distribution method** (file share recommended for quick start)
2. **Share with pilot group** (2-3 engineers)
3. **Gather feedback** and iterate
4. **Scale to full team**
5. **Integrate into on-call training**

**Ready to transform your team's incident response!** 🏦🤖