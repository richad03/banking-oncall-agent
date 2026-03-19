package models

import (
	"strings"
	"time"
)

// ServiceType represents the different banking services
type ServiceType string

const (
	ServiceInflows         ServiceType = "inflows"
	ServiceOutflows        ServiceType = "outflows"
	ServiceSavings         ServiceType = "savings"
	ServiceChecking        ServiceType = "checking"
	ServicePayroll         ServiceType = "payroll"
	ServiceBalanceReporter ServiceType = "balance-reporter"
)

// IncidentSeverity levels for triage
type IncidentSeverity string

const (
	SeverityLow      IncidentSeverity = "low"
	SeverityMedium   IncidentSeverity = "medium"
	SeverityHigh     IncidentSeverity = "high"
	SeverityCritical IncidentSeverity = "critical"
)

// Incident represents a banking service incident
type Incident struct {
	ID          string           `json:"id"`
	Service     ServiceType      `json:"service"`
	ErrorCode   string           `json:"error_code"`
	Description string           `json:"description"`
	Severity    IncidentSeverity `json:"severity"`
	Timestamp   time.Time        `json:"timestamp"`
	UserID      string           `json:"user_id"`
	ChannelID   string           `json:"channel_id"`
}

// RunbookSection represents a section of incident response documentation
type RunbookSection struct {
	Service     ServiceType `json:"service"`
	Section     string      `json:"section"`
	Title       string      `json:"title"`
	ErrorCodes  []string    `json:"error_codes"`
	Description string      `json:"description"`
	Steps       []string    `json:"steps"`
	AdminLinks  []string    `json:"admin_links"`
	References  []string    `json:"references"`
}

// ServiceDependency represents service relationships
type ServiceDependency struct {
	Service      ServiceType   `json:"service"`
	Upstream     []ServiceType `json:"upstream"`
	Downstream   []ServiceType `json:"downstream"`
	CriticalPath bool          `json:"critical_path"`
}

// DetectService attempts to identify the service from incident text
func DetectService(text string) ServiceType {
	text = strings.ToLower(text)

	// Direct service mentions
	if strings.Contains(text, "payroll") || strings.Contains(text, "reckoner") || strings.Contains(text, "tax") {
		return ServicePayroll
	}
	if strings.Contains(text, "deposit") || strings.Contains(text, "outflow") || strings.Contains(text, "billpay") || strings.Contains(text, "checkbook") {
		return ServiceOutflows
	}
	if strings.Contains(text, "inflow") || strings.Contains(text, "ach") || strings.Contains(text, "banknote") {
		return ServiceInflows
	}
	if strings.Contains(text, "saving") || strings.Contains(text, "sfs") {
		return ServiceSavings
	}
	if strings.Contains(text, "checking") {
		return ServiceChecking
	}
	if strings.Contains(text, "balance") || strings.Contains(text, "reporter") {
		return ServiceBalanceReporter
	}

	// Error code patterns
	if strings.Contains(text, "pt-") || strings.Contains(text, "pr-") {
		return ServicePayroll
	}
	if strings.Contains(text, "dep-") || strings.Contains(text, "bp-") || strings.Contains(text, "cb-") {
		return ServiceOutflows
	}
	if strings.Contains(text, "inf-") || strings.Contains(text, "bn-") {
		return ServiceInflows
	}

	return ""
}

// DetectSeverity estimates incident severity based on keywords
func DetectSeverity(text string) IncidentSeverity {
	text = strings.ToLower(text)

	// Critical indicators
	criticalKeywords := []string{"down", "outage", "critical", "sev1", "sev-1", "p1", "urgent", "emergency"}
	for _, keyword := range criticalKeywords {
		if strings.Contains(text, keyword) {
			return SeverityCritical
		}
	}

	// High severity indicators
	highKeywords := []string{"error", "failure", "timeout", "sev2", "sev-2", "p2", "stuck", "broken"}
	for _, keyword := range highKeywords {
		if strings.Contains(text, keyword) {
			return SeverityHigh
		}
	}

	// Medium severity indicators
	mediumKeywords := []string{"slow", "latency", "warning", "sev3", "sev-3", "p3", "degraded"}
	for _, keyword := range mediumKeywords {
		if strings.Contains(text, keyword) {
			return SeverityMedium
		}
	}

	return SeverityLow
}
