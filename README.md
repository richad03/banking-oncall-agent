# Banking On-Call Agent

An intelligent Slack bot that assists Square banking engineers with incident response, built using Go and Claude AI.

## Features

🚨 **Incident Triage & Analysis**
- Automatically detects banking service incidents from Slack messages
- Classifies severity levels (Critical, High, Medium, Low)
- Identifies affected services (Payroll, Outflows, Inflows, Savings, etc.)

🤖 **AI-Powered Assistance**
- Uses Claude AI to provide intelligent incident analysis
- Searches knowledge base for relevant runbook sections
- Suggests immediate actions based on error codes and symptoms

📚 **Knowledge Base Integration**
- Indexes all banking-oncall documentation
- Maps error codes to specific resolution procedures
- Provides service dependency analysis for impact assessment

💬 **Slack Integration**
- Responds to direct mentions and trigger words
- Supports slash commands for help and status
- Formats responses with actionable information

## Architecture

```
banking-oncall-agent/
├── cmd/bot/                   # Main application entry point
├── internal/
│   ├── agent/                 # Core AI agent logic
│   ├── knowledge/             # Knowledge base management
│   ├── slack/                 # Slack bot integration
│   └── config/                # Configuration management
├── pkg/
│   ├── models/                # Data models and types
│   └── utils/                 # Utility functions
└── docs/                      # Documentation
```

## Quick Start

### Prerequisites

- Go 1.21+
- Slack workspace with bot permissions
- Claude API key (Anthropic)
- Access to banking-oncall repository

### Installation

1. **Clone and build**
```bash
git clone <repository-url>
cd banking-oncall-agent
go mod tidy
```

2. **Set up configuration**
```bash
cp .env.example .env
# Edit .env with your credentials
```

3. **Configure Slack App**
   - Create a Slack app at https://api.slack.com/apps
   - Enable Socket Mode and get App Token
   - Create Bot User and get Bot Token
   - Add required bot scopes:
     - `app_mentions:read`
     - `channels:history`
     - `chat:write`
     - `commands`
     - `im:history`

4. **Build and run**
```bash
go build -o banking-agent cmd/bot/main.go
./banking-agent
```

### Docker Deployment

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN go build -o banking-agent cmd/bot/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/banking-agent .
CMD ["./banking-agent"]
```

## Usage

### In Slack Channels

**Direct mentions:**
```
@banking-agent I'm seeing PT-401 errors in payroll
@banking-agent Deposits service appears down
@banking-agent Help with checkbook timeout issues
```

**Trigger words (automatic response):**
```
Seeing incident with billpay processing
Error DEP-500 in deposits
Outage affecting savings accounts
```

**Slash commands:**
```
/banking-help          # Show help and examples
/banking-status        # Show service overview
```

### Response Format

```
🚨 **Payroll Incident - High**
🔍 **Error Code:** PT-401

🤖 **Analysis:**
Tax calculation failure detected in payroll core service.
This typically indicates reckoner worker issues.

🔧 **Immediate Actions:**
• Check reckoner worker status in admin console
• Restart worker if needed
• Verify tax calculation endpoints

📖 **Runbook:** Section 23.4 - Tax Processing Errors
🔗 **Admin Consoles:**
• go/payrolladmin
• go/reckoneradmin

⚠️ **Escalation recommended** - Contact on-call manager
```

## Configuration

### Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `SLACK_BOT_TOKEN` | Bot User OAuth Token (xoxb-...) | Yes |
| `SLACK_APP_TOKEN` | App-Level Token (xapp-...) | Yes |
| `ANTHROPIC_API_KEY` | Claude API key | Yes |
| `KNOWLEDGE_BASE_PATH` | Path to banking-oncall repo | No |
| `PORT` | Application port | No |
| `ENVIRONMENT` | deployment environment | No |

### Slack App Configuration

**OAuth Scopes (Bot Token):**
- `app_mentions:read` - Detect @mentions
- `channels:history` - Read channel messages
- `chat:write` - Send messages
- `commands` - Handle slash commands
- `im:history` - Read DMs

**Event Subscriptions:**
- `app_mention` - Bot mentions
- `message.channels` - Channel messages (if needed)
- `message.im` - Direct messages

**Slash Commands:**
- `/banking-help` - Show help information
- `/banking-status` - Show service status

## Supported Services & Error Codes

| Service | Error Codes | Documentation |
|---------|-------------|---------------|
| **Payroll** | PT-*, PR-* | Payroll processing, tax calculations |
| **Outflows** | DEP-*, BP-*, CB-* | Deposits, Billpay, Checkbook |
| **Inflows** | INF-*, BN-* | ACH deposits, Banknotes |
| **Savings** | SFS-* | Savings accounts, SFS core |
| **Checking** | CHK-* | Checking accounts |
| **Balance Reporter** | BR-* | Balance aggregation |

## Development

### Project Structure

- **`cmd/bot/main.go`** - Application entry point
- **`internal/agent/`** - Core AI logic using Claude API
- **`internal/knowledge/`** - Knowledge base indexing and search
- **`internal/slack/`** - Slack event handling and bot logic
- **`internal/config/`** - Configuration management
- **`pkg/models/`** - Domain models and data structures

### Adding New Services

1. Update `ServiceType` enum in `pkg/models/incident.go`
2. Add detection patterns in `DetectService()` function
3. Update knowledge base parsing in `internal/knowledge/`
4. Add service-specific logic in agent response generation

### Testing

```bash
# Run tests
go test ./...

# Run with race detection
go test -race ./...

# Integration tests (requires valid tokens)
go test -tags=integration ./...
```

## Deployment

### Production Checklist

- [ ] Set up proper logging and monitoring
- [ ] Configure health checks
- [ ] Set resource limits and auto-scaling
- [ ] Set up proper secret management
- [ ] Configure alerting for bot failures
- [ ] Set up knowledge base update automation

### Monitoring

The bot logs important events:
- Incident processing starts/completions
- API errors and failures
- Knowledge base load status
- Slack connection status

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Submit a pull request

## Troubleshooting

### Common Issues

**Bot not responding:**
- Check Slack token validity
- Verify bot permissions and scopes
- Check event subscription endpoints

**Knowledge base errors:**
- Ensure banking-oncall repository path is correct
- Verify file permissions
- Check for malformed markdown

**Claude API errors:**
- Verify API key is valid
- Check rate limits
- Monitor token usage

### Logs

```bash
# View application logs
tail -f /var/log/banking-agent.log

# Check knowledge base loading
grep "Knowledge base" /var/log/banking-agent.log

# Monitor incident processing
grep "Processing incident" /var/log/banking-agent.log
```

## License

Internal Square tool - See company licensing policies.

---

Built with ❤️ by the Square Banking Engineering team