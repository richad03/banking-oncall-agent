package knowledge

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"banking-oncall-agent/pkg/models"
)

// KnowledgeBase manages the banking-oncall documentation
type KnowledgeBase struct {
	basePath      string
	runbooks      map[models.ServiceType][]*models.RunbookSection
	dependencies  map[models.ServiceType]*models.ServiceDependency
	errorPatterns map[string]*models.RunbookSection
}

// NewKnowledgeBase creates and initializes a knowledge base
func NewKnowledgeBase(basePath string) *KnowledgeBase {
	return &KnowledgeBase{
		basePath:      basePath,
		runbooks:      make(map[models.ServiceType][]*models.RunbookSection),
		dependencies:  make(map[models.ServiceType]*models.ServiceDependency),
		errorPatterns: make(map[string]*models.RunbookSection),
	}
}

// Load indexes all documentation from the repository
func (kb *KnowledgeBase) Load() error {
	fmt.Printf("Loading knowledge base from: %s\n", kb.basePath)

	// Walk through all markdown files
	return filepath.WalkDir(kb.basePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !strings.HasSuffix(path, ".md") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", path, err)
		}

		return kb.parseDocument(path, string(content))
	})
}

// parseDocument extracts knowledge from a markdown document
func (kb *KnowledgeBase) parseDocument(filePath, content string) error {
	// Determine service from file path
	service := kb.detectServiceFromPath(filePath)
	if service == "" {
		return nil // Skip non-service files
	}

	// Parse runbook sections if it's a runbook
	if strings.Contains(filePath, "runbook") {
		sections := kb.parseRunbook(service, content)
		kb.runbooks[service] = append(kb.runbooks[service], sections...)

		// Index error patterns
		for _, section := range sections {
			for _, errorCode := range section.ErrorCodes {
				kb.errorPatterns[strings.ToUpper(errorCode)] = section
			}
		}
	}

	// Parse service dependencies if it's a service overview
	if strings.Contains(filePath, "service-overview") || strings.Contains(filePath, "README") {
		deps := kb.parseDependencies(service, content)
		if deps != nil {
			kb.dependencies[service] = deps
		}
	}

	return nil
}

// detectServiceFromPath determines service type from file path
func (kb *KnowledgeBase) detectServiceFromPath(path string) models.ServiceType {
	path = strings.ToLower(path)

	if strings.Contains(path, "/payroll/") {
		return models.ServicePayroll
	}
	if strings.Contains(path, "/outflows/") {
		return models.ServiceOutflows
	}
	if strings.Contains(path, "/inflows/") {
		return models.ServiceInflows
	}
	if strings.Contains(path, "/savings/") {
		return models.ServiceSavings
	}
	if strings.Contains(path, "/checking/") {
		return models.ServiceChecking
	}
	if strings.Contains(path, "/balance-reporter/") {
		return models.ServiceBalanceReporter
	}

	return ""
}

// parseRunbook extracts runbook sections from markdown content
func (kb *KnowledgeBase) parseRunbook(service models.ServiceType, content string) []*models.RunbookSection {
	var sections []*models.RunbookSection

	// Split by headers (## or ###)
	headerRegex := regexp.MustCompile(`(?m)^(#{2,3})\s+(.+)$`)
	matches := headerRegex.FindAllStringSubmatchIndex(content, -1)

	for i, match := range matches {
		start := match[0]
		var end int
		if i+1 < len(matches) {
			end = matches[i+1][0]
		} else {
			end = len(content)
		}

		sectionContent := content[start:end]
		title := content[match[4]:match[5]]

		section := &models.RunbookSection{
			Service:     service,
			Section:     fmt.Sprintf("%d", i+1),
			Title:       strings.TrimSpace(title),
			Description: kb.extractDescription(sectionContent),
			Steps:       kb.extractSteps(sectionContent),
			ErrorCodes:  kb.extractErrorCodes(sectionContent),
			AdminLinks:  kb.extractLinks(sectionContent, "go/"),
			References:  kb.extractLinks(sectionContent, "http"),
		}

		sections = append(sections, section)
	}

	return sections
}

// parseDependencies extracts service dependencies from content
func (kb *KnowledgeBase) parseDependencies(service models.ServiceType, content string) *models.ServiceDependency {
	// Look for dependency mentions in the content
	upstream := kb.findServiceMentions(content, "depends on", "calls", "upstream")
	downstream := kb.findServiceMentions(content, "called by", "downstream", "consumers")

	if len(upstream) == 0 && len(downstream) == 0 {
		return nil
	}

	return &models.ServiceDependency{
		Service:      service,
		Upstream:     upstream,
		Downstream:   downstream,
		CriticalPath: strings.Contains(strings.ToLower(content), "critical"),
	}
}

// Helper methods for parsing

func (kb *KnowledgeBase) extractDescription(content string) string {
	lines := strings.Split(content, "\n")
	var description []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") && !strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "*") {
			description = append(description, line)
			if len(description) >= 3 { // First few non-header lines
				break
			}
		}
	}

	return strings.Join(description, " ")
}

func (kb *KnowledgeBase) extractSteps(content string) []string {
	var steps []string
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* ") || strings.HasPrefix(line, "1. ") {
			steps = append(steps, strings.TrimPrefix(strings.TrimPrefix(strings.TrimPrefix(line, "- "), "* "), "1. "))
		}
	}

	return steps
}

func (kb *KnowledgeBase) extractErrorCodes(content string) []string {
	// Common error code patterns: PT-401, DEP-500, etc.
	errorRegex := regexp.MustCompile(`\b[A-Z]{2,4}-\d{3,4}\b`)
	matches := errorRegex.FindAllString(content, -1)

	// Deduplicate
	seen := make(map[string]bool)
	var unique []string
	for _, match := range matches {
		if !seen[match] {
			seen[match] = true
			unique = append(unique, match)
		}
	}

	return unique
}

func (kb *KnowledgeBase) extractLinks(content, prefix string) []string {
	var links []string
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		if strings.Contains(line, prefix) {
			// Extract URLs or go/ links
			if prefix == "go/" {
				goRegex := regexp.MustCompile(`go/[a-zA-Z0-9-_]+`)
				matches := goRegex.FindAllString(line, -1)
				links = append(links, matches...)
			} else {
				urlRegex := regexp.MustCompile(`https?://[^\s)]+`)
				matches := urlRegex.FindAllString(line, -1)
				links = append(links, matches...)
			}
		}
	}

	return links
}

func (kb *KnowledgeBase) findServiceMentions(content string, keywords ...string) []models.ServiceType {
	var services []models.ServiceType
	content = strings.ToLower(content)

	serviceNames := map[string]models.ServiceType{
		"payroll":          models.ServicePayroll,
		"reckoner":         models.ServicePayroll,
		"deposit":          models.ServiceOutflows,
		"outflow":          models.ServiceOutflows,
		"billpay":          models.ServiceOutflows,
		"checkbook":        models.ServiceOutflows,
		"inflow":           models.ServiceInflows,
		"banknote":         models.ServiceInflows,
		"saving":           models.ServiceSavings,
		"sfs":              models.ServiceSavings,
		"checking":         models.ServiceChecking,
		"balance":          models.ServiceBalanceReporter,
		"balance-reporter": models.ServiceBalanceReporter,
	}

	for serviceName, serviceType := range serviceNames {
		if strings.Contains(content, serviceName) {
			services = append(services, serviceType)
		}
	}

	return services
}

// Query methods

// FindRunbookSection searches for relevant runbook sections
func (kb *KnowledgeBase) FindRunbookSection(service models.ServiceType, errorCode string) *models.RunbookSection {
	// First try exact error code match
	if errorCode != "" {
		if section, exists := kb.errorPatterns[strings.ToUpper(errorCode)]; exists {
			return section
		}
	}

	// Fall back to service runbook
	if sections, exists := kb.runbooks[service]; exists && len(sections) > 0 {
		return sections[0] // Return first section as default
	}

	return nil
}

// GetServiceDependencies returns dependency information for a service
func (kb *KnowledgeBase) GetServiceDependencies(service models.ServiceType) *models.ServiceDependency {
	return kb.dependencies[service]
}

// SearchContent performs full-text search across all documentation
func (kb *KnowledgeBase) SearchContent(query string) []*models.RunbookSection {
	var results []*models.RunbookSection
	query = strings.ToLower(query)

	for _, sections := range kb.runbooks {
		for _, section := range sections {
			score := 0

			// Check title match
			if strings.Contains(strings.ToLower(section.Title), query) {
				score += 10
			}

			// Check description match
			if strings.Contains(strings.ToLower(section.Description), query) {
				score += 5
			}

			// Check error codes
			for _, errorCode := range section.ErrorCodes {
				if strings.Contains(strings.ToLower(errorCode), query) {
					score += 15
				}
			}

			// Check steps
			for _, step := range section.Steps {
				if strings.Contains(strings.ToLower(step), query) {
					score += 3
				}
			}

			if score >= 5 { // Minimum relevance threshold
				results = append(results, section)
			}
		}
	}

	return results
}

// GetServiceOverview returns a summary of available services
func (kb *KnowledgeBase) GetServiceOverview() map[models.ServiceType]int {
	overview := make(map[models.ServiceType]int)
	for service, sections := range kb.runbooks {
		overview[service] = len(sections)
	}
	return overview
}
