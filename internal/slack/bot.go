package slack

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"banking-oncall-agent/internal/agent"
	"banking-oncall-agent/pkg/models"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

// Bot represents the Slack bot instance
type Bot struct {
	client      *slack.Client
	socketMode  *socketmode.Client
	agent       *agent.BankingAgent
	botUserID   string
	triggerWords []string
}

// NewBot creates a new Slack bot
func NewBot(botToken, appToken string, bankingAgent *agent.BankingAgent) (*Bot, error) {
	api := slack.New(botToken, slack.OptionDebug(false))

	socketClient := socketmode.New(
		api,
		socketmode.OptionDebug(false),
		socketmode.OptionLog(log.New(log.Writer(), "socketmode: ", log.Lshortfile|log.LstdFlags)),
	)

	// Get bot user ID
	authResp, err := api.AuthTest()
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate: %w", err)
	}

	bot := &Bot{
		client:     api,
		socketMode: socketClient,
		agent:      bankingAgent,
		botUserID:  authResp.UserID,
		triggerWords: []string{
			"banking-agent",
			"oncall",
			"incident",
			"error",
			"outage",
			"down",
			"issue",
			"help",
		},
	}

	return bot, nil
}

// Start begins listening for Slack events
func (b *Bot) Start() error {
	log.Printf("Starting Banking On-Call Bot (User ID: %s)", b.botUserID)

	go b.handleEvents()

	return b.socketMode.Run()
}

// handleEvents processes incoming Slack events
func (b *Bot) handleEvents() {
	for envelope := range b.socketMode.Events {
		switch envelope.Type {
		case socketmode.EventTypeEventsAPI:
			b.handleEventsAPI(envelope)
		case socketmode.EventTypeInteractive:
			b.handleInteractiveEvent(envelope)
		case socketmode.EventTypeSlashCommand:
			b.handleSlashCommand(envelope)
		}
	}
}

// handleEventsAPI processes Events API events
func (b *Bot) handleEventsAPI(envelope *socketmode.Event) {
	eventsAPIEvent, ok := envelope.Data.(slackevents.EventsAPIEvent)
	if !ok {
		log.Printf("Could not type cast the event to the EventsAPIEvent: %v", envelope.Data)
		b.socketMode.Ack(*envelope.Request)
		return
	}

	b.socketMode.Ack(*envelope.Request)

	switch eventsAPIEvent.Type {
	case slackevents.CallbackEvent:
		innerEvent := eventsAPIEvent.InnerEvent
		switch event := innerEvent.Data.(type) {
		case *slackevents.MessageEvent:
			b.handleMessageEvent(event)
		case *slackevents.AppMentionEvent:
			b.handleAppMentionEvent(event)
		}
	}
}

// handleMessageEvent processes message events
func (b *Bot) handleMessageEvent(event *slackevents.MessageEvent) {
	// Skip bot messages and empty messages
	if event.User == b.botUserID || event.Text == "" {
		return
	}

	// Check if message should trigger the bot
	if b.shouldRespondToMessage(event.Text, event.Channel) {
		b.processIncidentMessage(event.Channel, event.User, event.Text, event.TimeStamp)
	}
}

// handleAppMentionEvent processes app mention events
func (b *Bot) handleAppMentionEvent(event *slackevents.AppMentionEvent) {
	// Remove the bot mention from the text
	text := b.removeMentions(event.Text)
	b.processIncidentMessage(event.Channel, event.User, text, event.TimeStamp)
}

// shouldRespondToMessage determines if the bot should respond to a message
func (b *Bot) shouldRespondToMessage(text, channelID string) bool {
	text = strings.ToLower(text)

	// Always respond in DMs
	if strings.HasPrefix(channelID, "D") {
		return true
	}

	// Check for trigger words in channels
	for _, word := range b.triggerWords {
		if strings.Contains(text, word) {
			return true
		}
	}

	// Check for error code patterns
	errorPattern := regexp.MustCompile(`\b[A-Z]{2,4}-\d{3,4}\b`)
	return errorPattern.MatchString(strings.ToUpper(text))
}

// processIncidentMessage analyzes a message as a potential incident
func (b *Bot) processIncidentMessage(channelID, userID, text, timestamp string) {
	log.Printf("Processing incident message from %s in %s: %s", userID, channelID, text)

	// Post initial acknowledgment
	b.postMessage(channelID, "🔍 Analyzing incident... one moment please.", "")

	// Parse incident from message
	incident := b.parseIncident(text, userID, channelID, timestamp)

	// Process with agent
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	response, err := b.agent.ProcessIncident(ctx, incident)
	if err != nil {
		log.Printf("Failed to process incident: %v", err)
		b.postMessage(channelID, "❌ Sorry, I encountered an error analyzing this incident. Please try again or contact the on-call team directly.", "")
		return
	}

	// Post the response
	b.postIncidentResponse(channelID, response)

	log.Printf("Successfully processed incident %s", incident.ID)
}

// parseIncident extracts incident information from message text
func (b *Bot) parseIncident(text, userID, channelID, timestamp string) *models.Incident {
	// Extract error code if present
	errorPattern := regexp.MustCompile(`\b([A-Z]{2,4}-\d{3,4})\b`)
	errorMatches := errorPattern.FindStringSubmatch(strings.ToUpper(text))
	errorCode := ""
	if len(errorMatches) > 1 {
		errorCode = errorMatches[1]
	}

	// Parse timestamp
	eventTime := time.Now()
	if timestamp != "" {
		if ts, err := time.Parse("1136239445.000000", timestamp); err == nil {
			eventTime = ts
		}
	}

	incident := &models.Incident{
		ID:          fmt.Sprintf("inc_%s_%s", channelID, timestamp),
		Description: text,
		ErrorCode:   errorCode,
		UserID:      userID,
		ChannelID:   channelID,
		Timestamp:   eventTime,
	}

	// Auto-detect service and severity
	incident.Service = models.DetectService(text)
	incident.Severity = models.DetectSeverity(text)

	return incident
}

// postIncidentResponse formats and posts the agent's response
func (b *Bot) postIncidentResponse(channelID string, response *agent.IncidentResponse) {
	// Use the formatted response from the agent
	b.postMessage(channelID, response.FormattedResponse, "")

	// If there are many actions, post them as a thread
	if len(response.ImmediateActions) > 3 {
		actionText := "🔧 **Complete Action List:**\n"
		for i, action := range response.ImmediateActions {
			actionText += fmt.Sprintf("%d. %s\n", i+1, action)
		}
		b.postMessage(channelID, actionText, "")
	}
}

// postMessage sends a message to a Slack channel
func (b *Bot) postMessage(channelID, text, threadTS string) {
	options := []slack.MsgOption{
		slack.MsgOptionText(text, false),
	}

	if threadTS != "" {
		options = append(options, slack.MsgOptionTS(threadTS))
	}

	_, _, err := b.client.PostMessage(channelID, options...)
	if err != nil {
		log.Printf("Failed to post message to %s: %v", channelID, err)
	}
}

// removeMentions removes bot mentions from text
func (b *Bot) removeMentions(text string) string {
	mentionPattern := regexp.MustCompile(`<@[A-Z0-9]+>`)
	return strings.TrimSpace(mentionPattern.ReplaceAllString(text, ""))
}

// handleInteractiveEvent processes interactive component events
func (b *Bot) handleInteractiveEvent(envelope *socketmode.Event) {
	// For future interactive features like buttons or select menus
	b.socketMode.Ack(*envelope.Request)
}

// handleSlashCommand processes slash commands
func (b *Bot) handleSlashCommand(envelope *socketmode.Event) {
	cmd, ok := envelope.Data.(slack.SlashCommand)
	if !ok {
		b.socketMode.Ack(*envelope.Request)
		return
	}

	switch cmd.Command {
	case "/banking-help":
		b.handleHelpCommand(cmd)
	case "/banking-status":
		b.handleStatusCommand(cmd)
	default:
		b.socketMode.Ack(*envelope.Request, map[string]interface{}{
			"text": "Unknown command. Use `/banking-help` for assistance.",
		})
	}
}

// handleHelpCommand shows help information
func (b *Bot) handleHelpCommand(cmd slack.SlashCommand) {
	helpText := `🏦 **Banking On-Call Assistant Help**

**How to use:**
• Mention me with incident details: @banking-agent payroll error PT-401
• Use trigger words: "incident", "error", "outage", "down"
• Include error codes for specific guidance: DEP-500, BP-401, etc.

**Commands:**
• \`/banking-help\` - Show this help
• \`/banking-status\` - Show service status

**Supported Services:**
• Payroll (PT-, PR- errors)
• Outflows: Deposits, Billpay, Checkbook (DEP-, BP-, CB- errors)
• Inflows: ACH, Banknotes (INF-, BN- errors)
• Savings: SFS (SFS- errors)
• Checking accounts
• Balance Reporter (BR- errors)

**Example messages:**
• "Seeing timeouts in deposits service"
• "PT-401 error in payroll processing"
• "Checkbook service appears down"
• "Need help with billpay incident"`

	b.socketMode.Ack(*envelope.Request, map[string]interface{}{
		"text": helpText,
	})
}

// handleStatusCommand shows service status summary
func (b *Bot) handleStatusCommand(cmd slack.SlashCommand) {
	// This could be extended to check actual service health
	statusText := `🏦 **Banking Services Overview**

**Available Services:**
• 🏦 **Payroll** - Processing, tax calculations, reckoner
• 💸 **Outflows** - Deposits, Billpay, Checkbook
• 💰 **Inflows** - ACH deposits, instant payouts, banknotes
• 💳 **Savings** - SFS core, account management
• ✅ **Checking** - Square checking accounts
• 📊 **Balance Reporter** - Balance aggregation

**Knowledge Base:** Loaded with runbooks and procedures
**Status:** Ready to assist with incidents

Use @banking-agent with your incident details for help!`

	b.socketMode.Ack(*envelope.Request, map[string]interface{}{
		"text": statusText,
	})
}