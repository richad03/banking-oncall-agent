package agent

import (
	"context"
	"fmt"
	"strings"
	"time"

	"banking-oncall-agent/internal/knowledge"
	"banking-oncall-agent/pkg/models"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

// BankingAgent provides intelligent incident assistance
type BankingAgent struct {
	client        *anthropic.Client
	knowledgeBase *knowledge.KnowledgeBase
}

// NewBankingAgent creates a new incident assistant agent
func NewBankingAgent(apiKey string, kb *knowledge.KnowledgeBase) *BankingAgent {
	client := anthropic.NewClient(option.WithAPIKey(apiKey))

	return &BankingAgent{
		client:        client,
		knowledgeBase: kb,
	}
}

// ProcessIncident analyzes an incident and provides guidance
func (ba *BankingAgent) ProcessIncident(ctx context.Context, incident *models.Incident) (*IncidentResponse, error) {
	// 1. Detect service and severity if not provided
	if incident.Service == "" {
		incident.Service = models.DetectService(incident.Description)
	}
	if incident.Severity == "" {
		incident.Severity = models.DetectSeverity(incident.Description)
	}

	// 2. Find relevant runbook sections
	runbookSection := ba.knowledgeBase.FindRunbookSection(incident.Service, incident.ErrorCode)

	// 3. Get service dependencies for impact analysis
	dependencies := ba.knowledgeBase.GetServiceDependencies(incident.Service)

	// 4. Search for additional context
	searchResults := ba.knowledgeBase.SearchContent(incident.Description + " " + incident.ErrorCode)

	// 5. Generate intelligent response using Claude
	response, err := ba.generateResponse(ctx, incident, runbookSection, dependencies, searchResults)
	if err != nil {
		return nil, fmt.Errorf("failed to generate response: %w", err)
	}

	return response, nil
}

// IncidentResponse contains the agent's analysis and recommendations
type IncidentResponse struct {
	Summary           string                    `json:"summary"`
	Severity          models.IncidentSeverity   `json:"severity"`
	Service           models.ServiceType        `json:"service"`
	ImmediateActions  []string                  `json:"immediate_actions"`
	RunbookSection    *models.RunbookSection    `json:"runbook_section,omitempty"`
	Dependencies      *models.ServiceDependency `json:"dependencies,omitempty"`
	AdminConsoles     []string                  `json:"admin_consoles"`
	RecentIncidents   []string                  `json:"recent_incidents"`
	EscalationNeeded  bool                      `json:"escalation_needed"`
	FormattedResponse string                    `json:"formatted_response"`
}

// generateResponse uses Claude to create an intelligent incident response
func (ba *BankingAgent) generateResponse(ctx context.Context, incident *models.Incident, runbook *models.RunbookSection, deps *models.ServiceDependency, searchResults []*models.RunbookSection) (*IncidentResponse, error) {
	// Build context from knowledge base
	context := ba.buildContext(incident, runbook, deps, searchResults)

	// Create the prompt for Claude
	prompt := ba.buildPrompt(incident, context)

	// Call Claude API
	resp, err := ba.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.ModelClaude3_5Sonnet20241022,
		MaxTokens: anthropic.F(1500),
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		},
	})

	if err != nil {
		return nil, fmt.Errorf("Claude API error: %w", err)
	}

	// Parse Claude's response and structure it
	claudeResponse := ""
	for _, block := range resp.Content {
		if textBlock, ok := block.(*anthropic.TextBlock); ok {
			claudeResponse += textBlock.Text
		}
	}

	// Structure the response
	response := &IncidentResponse{
		Summary:           ba.extractSummary(claudeResponse),
		Severity:          incident.Severity,
		Service:           incident.Service,
		ImmediateActions:  ba.extractActions(claudeResponse),
		RunbookSection:    runbook,
		Dependencies:      deps,
		AdminConsoles:     ba.extractAdminLinks(runbook),
		EscalationNeeded:  ba.shouldEscalate(incident.Severity, claudeResponse),
		FormattedResponse: ba.formatSlackResponse(incident, claudeResponse, runbook),
	}

	return response, nil
}

// buildContext creates relevant context from the knowledge base
func (ba *BankingAgent) buildContext(incident *models.Incident, runbook *models.RunbookSection, deps *models.ServiceDependency, searchResults []*models.RunbookSection) string {
	var contextParts []string

	// Service information
	contextParts = append(contextParts, fmt.Sprintf("SERVICE: %s", incident.Service))

	// Runbook section if available
	if runbook != nil {
		contextParts = append(contextParts, fmt.Sprintf("RUNBOOK SECTION: %s", runbook.Title))
		if len(runbook.Steps) > 0 {
			contextParts = append(contextParts, "RESOLUTION STEPS:")
			for i, step := range runbook.Steps {
				if i < 5 { // Limit to first 5 steps
					contextParts = append(contextParts, fmt.Sprintf("- %s", step))
				}
			}
		}
		if len(runbook.ErrorCodes) > 0 {
			contextParts = append(contextParts, fmt.Sprintf("RELATED ERROR CODES: %s", strings.Join(runbook.ErrorCodes, ", ")))
		}
	}

	// Dependencies for impact analysis
	if deps != nil {
		if len(deps.Upstream) > 0 {
			contextParts = append(contextParts, fmt.Sprintf("UPSTREAM DEPENDENCIES: %s", strings.Join(convertServicesToStrings(deps.Upstream), ", ")))
		}
		if len(deps.Downstream) > 0 {
			contextParts = append(contextParts, fmt.Sprintf("DOWNSTREAM IMPACT: %s", strings.Join(convertServicesToStrings(deps.Downstream), ", ")))
		}
	}

	// Additional relevant sections
	if len(searchResults) > 0 {
		contextParts = append(contextParts, "RELATED DOCUMENTATION:")
		for i, result := range searchResults {
			if i < 3 { // Limit to top 3 results
				contextParts = append(contextParts, fmt.Sprintf("- %s: %s", result.Title, result.Description))
			}
		}
	}

	return strings.Join(contextParts, "\n")
}

// buildPrompt creates the prompt for Claude
func (ba *BankingAgent) buildPrompt(incident *models.Incident, context string) string {
	return fmt.Sprintf(`You are a Square Banking On-Call Assistant. You help engineers respond to banking service incidents quickly and effectively.

INCIDENT DETAILS:
Service: %s
Error Code: %s
Description: %s
Severity: %s
Timestamp: %s

KNOWLEDGE BASE CONTEXT:
%s

Please provide:
1. A brief summary of the incident (1-2 sentences)
2. Immediate actions to take (prioritized list)
3. Whether this needs escalation based on severity
4. Any relevant admin consoles or tools to check

Format your response to be clear and actionable for an on-call engineer. Focus on practical next steps rather than lengthy explanations.

If this is a known error pattern from the context, reference the specific runbook section and steps. If it affects multiple services, mention the impact scope.

Keep the tone professional but concise - this is for incident response where time matters.`,
		incident.Service,
		incident.ErrorCode,
		incident.Description,
		incident.Severity,
		incident.Timestamp.Format(time.RFC3339),
		context)
}

// Helper methods for response parsing and formatting

func (ba *BankingAgent) extractSummary(response string) string {
	lines := strings.Split(response, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") && !strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "1.") {
			if len(line) > 20 { // Ensure it's substantial
				return line
			}
		}
	}
	return "Incident analysis in progress..."
}

func (ba *BankingAgent) extractActions(response string) []string {
	var actions []string
	lines := strings.Split(response, "\n")
	inActionSection := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Look for action list indicators
		if strings.Contains(strings.ToLower(line), "immediate") || strings.Contains(strings.ToLower(line), "action") || strings.Contains(strings.ToLower(line), "step") {
			inActionSection = true
			continue
		}

		// Extract bulleted items in action section
		if inActionSection && (strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* ") || strings.HasPrefix(line, "1. ") || strings.HasPrefix(line, "2. ")) {
			action := strings.TrimPrefix(strings.TrimPrefix(strings.TrimPrefix(line, "- "), "* "), "1. ")
			action = strings.TrimPrefix(action, "2. ")
			if action != "" {
				actions = append(actions, action)
			}
		} else if inActionSection && line == "" {
			break // End of action section
		}
	}

	// Fallback: extract any bulleted items
	if len(actions) == 0 {
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* ") {
				action := strings.TrimPrefix(strings.TrimPrefix(line, "- "), "* ")
				if action != "" && len(actions) < 5 {
					actions = append(actions, action)
				}
			}
		}
	}

	return actions
}

func (ba *BankingAgent) extractAdminLinks(runbook *models.RunbookSection) []string {
	if runbook == nil {
		return nil
	}
	return runbook.AdminLinks
}

func (ba *BankingAgent) shouldEscalate(severity models.IncidentSeverity, response string) bool {
	if severity == models.SeverityCritical {
		return true
	}

	escalationKeywords := []string{"escalate", "manager", "senior", "critical", "urgent"}
	responseLower := strings.ToLower(response)

	for _, keyword := range escalationKeywords {
		if strings.Contains(responseLower, keyword) {
			return true
		}
	}

	return false
}

func (ba *BankingAgent) formatSlackResponse(incident *models.Incident, claudeResponse string, runbook *models.RunbookSection) string {
	var parts []string

	// Header with service and severity
	severityEmoji := ba.getSeverityEmoji(incident.Severity)
	parts = append(parts, fmt.Sprintf("%s **%s Incident - %s**", severityEmoji, strings.Title(string(incident.Service)), strings.Title(string(incident.Severity))))

	if incident.ErrorCode != "" {
		parts = append(parts, fmt.Sprintf("🔍 **Error Code:** %s", incident.ErrorCode))
	}

	// Claude's analysis
	parts = append(parts, fmt.Sprintf("🤖 **Analysis:**\n%s", ba.extractSummary(claudeResponse)))

	// Immediate actions
	actions := ba.extractActions(claudeResponse)
	if len(actions) > 0 {
		parts = append(parts, "🔧 **Immediate Actions:**")
		for i, action := range actions {
			if i < 3 { // Show top 3 actions
				parts = append(parts, fmt.Sprintf("• %s", action))
			}
		}
	}

	// Runbook reference
	if runbook != nil {
		parts = append(parts, fmt.Sprintf("📖 **Runbook:** Section %s - %s", runbook.Section, runbook.Title))
	}

	// Admin consoles
	adminLinks := ba.extractAdminLinks(runbook)
	if len(adminLinks) > 0 {
		parts = append(parts, "🔗 **Admin Consoles:**")
		for _, link := range adminLinks {
			parts = append(parts, fmt.Sprintf("• %s", link))
		}
	}

	// Escalation notice
	if ba.shouldEscalate(incident.Severity, claudeResponse) {
		parts = append(parts, "⚠️ **Escalation recommended** - Contact on-call manager")
	}

	return strings.Join(parts, "\n")
}

func (ba *BankingAgent) getSeverityEmoji(severity models.IncidentSeverity) string {
	switch severity {
	case models.SeverityCritical:
		return "🚨"
	case models.SeverityHigh:
		return "🔥"
	case models.SeverityMedium:
		return "⚠️"
	default:
		return "ℹ️"
	}
}

// convertServicesToStrings is a helper to convert ServiceType slice to string slice
func convertServicesToStrings(services []models.ServiceType) []string {
	result := make([]string, len(services))
	for i, service := range services {
		result[i] = string(service)
	}
	return result
}
