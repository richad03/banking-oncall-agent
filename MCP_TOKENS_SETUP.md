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

