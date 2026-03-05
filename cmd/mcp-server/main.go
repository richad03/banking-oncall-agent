package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"banking-oncall-agent/internal/config"
	"banking-oncall-agent/internal/knowledge"
	"banking-oncall-agent/pkg/models"
)

// MCP Server for Banking On-Call Agent
// Exposes banking knowledge and incident analysis as MCP tools

type MCPServer struct {
	knowledgeBase *knowledge.KnowledgeBase
}

type MCPRequest struct {
	Jsonrpc string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

type MCPResponse struct {
	Jsonrpc string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

type MCPTool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize knowledge base
	log.Println("Loading banking knowledge base for MCP...")
	kb := knowledge.NewKnowledgeBase(cfg.KnowledgeBasePath)
	if err := kb.Load(); err != nil {
		log.Fatalf("Failed to load knowledge base: %v", err)
	}

	server := &MCPServer{
		knowledgeBase: kb,
	}

	log.Println("Banking On-Call MCP Server starting...")
	server.Run()
}

func (s *MCPServer) Run() {
	decoder := json.NewDecoder(os.Stdin)
	encoder := json.NewEncoder(os.Stdout)

	for {
		var request MCPRequest
		if err := decoder.Decode(&request); err != nil {
			log.Printf("Error decoding request: %v", err)
			continue
		}

		response := s.handleRequest(request)
		if err := encoder.Encode(response); err != nil {
			log.Printf("Error encoding response: %v", err)
		}
	}
}

func (s *MCPServer) handleRequest(req MCPRequest) MCPResponse {
	switch req.Method {
	case "initialize":
		return MCPResponse{
			Jsonrpc: "2.0",
			ID:      req.ID,
			Result: map[string]interface{}{
				"protocolVersion": "2024-11-05",
				"capabilities": map[string]interface{}{
					"tools": map[string]interface{}{},
				},
				"serverInfo": map[string]interface{}{
					"name":    "banking-oncall-agent",
					"version": "1.0.0",
				},
			},
		}

	case "tools/list":
		return MCPResponse{
			Jsonrpc: "2.0",
			ID:      req.ID,
			Result: map[string]interface{}{
				"tools": []MCPTool{
					{
						Name:        "analyze_banking_incident",
						Description: "Analyze a banking service incident and provide runbook guidance, error code mapping, and suggested actions",
						InputSchema: map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"description": map[string]interface{}{
									"type":        "string",
									"description": "Description of the banking incident or error",
								},
								"error_code": map[string]interface{}{
									"type":        "string",
									"description": "Specific error code if known (e.g., PT-401, DEP-500)",
								},
								"service": map[string]interface{}{
									"type":        "string",
									"description": "Banking service if known (payroll, outflows, inflows, savings, checking, balance-reporter)",
								},
							},
							"required": []string{"description"},
						},
					},
					{
						Name:        "search_banking_runbooks",
						Description: "Search banking on-call runbooks and documentation for specific procedures or error codes",
						InputSchema: map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"query": map[string]interface{}{
									"type":        "string",
									"description": "Search query for runbook procedures or error codes",
								},
								"service": map[string]interface{}{
									"type":        "string",
									"description": "Limit search to specific service (optional)",
								},
							},
							"required": []string{"query"},
						},
					},
					{
						Name:        "get_service_dependencies",
						Description: "Get dependency information for a banking service to understand impact scope",
						InputSchema: map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"service": map[string]interface{}{
									"type":        "string",
									"description": "Banking service name",
								},
							},
							"required": []string{"service"},
						},
					},
					{
						Name:        "list_banking_services",
						Description: "Get overview of all supported banking services and their runbook coverage",
						InputSchema: map[string]interface{}{
							"type":       "object",
							"properties": map[string]interface{}{},
						},
					},
				},
			},
		}

	case "tools/call":
		return s.handleToolCall(req)

	default:
		return MCPResponse{
			Jsonrpc: "2.0",
			ID:      req.ID,
			Error: map[string]interface{}{
				"code":    -32601,
				"message": "Method not found",
			},
		}
	}
}

func (s *MCPServer) handleToolCall(req MCPRequest) MCPResponse {
	params, ok := req.Params.(map[string]interface{})
	if !ok {
		return s.errorResponse(req.ID, "Invalid parameters")
	}

	toolName, ok := params["name"].(string)
	if !ok {
		return s.errorResponse(req.ID, "Missing tool name")
	}

	arguments, ok := params["arguments"].(map[string]interface{})
	if !ok {
		return s.errorResponse(req.ID, "Missing arguments")
	}

	switch toolName {
	case "analyze_banking_incident":
		return s.analyzeIncident(req.ID, arguments)

	case "search_banking_runbooks":
		return s.searchRunbooks(req.ID, arguments)

	case "get_service_dependencies":
		return s.getServiceDependencies(req.ID, arguments)

	case "list_banking_services":
		return s.listServices(req.ID)

	default:
		return s.errorResponse(req.ID, "Unknown tool")
	}
}

func (s *MCPServer) analyzeIncident(id interface{}, args map[string]interface{}) MCPResponse {
	description, _ := args["description"].(string)
	errorCode, _ := args["error_code"].(string)
	serviceName, _ := args["service"].(string)

	// Create incident object
	incident := &models.Incident{
		Description: description,
		ErrorCode:   errorCode,
	}

	// Detect service if not provided
	if serviceName != "" {
		incident.Service = models.ServiceType(serviceName)
	} else {
		incident.Service = models.DetectService(description)
	}

	// Detect severity
	incident.Severity = models.DetectSeverity(description)

	// Find relevant runbook section
	runbook := s.knowledgeBase.FindRunbookSection(incident.Service, errorCode)

	// Get service dependencies
	deps := s.knowledgeBase.GetServiceDependencies(incident.Service)

	// Search for additional context
	searchResults := s.knowledgeBase.SearchContent(description + " " + errorCode)

	// Build response
	result := map[string]interface{}{
		"incident_analysis": map[string]interface{}{
			"service":     string(incident.Service),
			"severity":    string(incident.Severity),
			"error_code":  errorCode,
			"description": description,
		},
	}

	if runbook != nil {
		result["runbook_guidance"] = map[string]interface{}{
			"section":     runbook.Section,
			"title":       runbook.Title,
			"description": runbook.Description,
			"steps":       runbook.Steps,
			"error_codes": runbook.ErrorCodes,
			"admin_links": runbook.AdminLinks,
			"references":  runbook.References,
		}
	}

	if deps != nil {
		result["service_dependencies"] = map[string]interface{}{
			"upstream":      deps.Upstream,
			"downstream":    deps.Downstream,
			"critical_path": deps.CriticalPath,
		}
	}

	if len(searchResults) > 0 {
		var relatedSections []map[string]interface{}
		for _, section := range searchResults {
			if len(relatedSections) < 3 { // Limit to top 3
				relatedSections = append(relatedSections, map[string]interface{}{
					"title":       section.Title,
					"description": section.Description,
					"service":     string(section.Service),
				})
			}
		}
		result["related_procedures"] = relatedSections
	}

	return MCPResponse{
		Jsonrpc: "2.0",
		ID:      id,
		Result: map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": s.formatIncidentAnalysis(result),
				},
			},
		},
	}
}

func (s *MCPServer) searchRunbooks(id interface{}, args map[string]interface{}) MCPResponse {
	query, _ := args["query"].(string)
	serviceName, _ := args["service"].(string)

	if query == "" {
		return s.errorResponse(id, "Query is required")
	}

	results := s.knowledgeBase.SearchContent(query)

	// Filter by service if specified
	if serviceName != "" {
		serviceType := models.ServiceType(serviceName)
		var filteredResults []*models.RunbookSection
		for _, result := range results {
			if result.Service == serviceType {
				filteredResults = append(filteredResults, result)
			}
		}
		results = filteredResults
	}

	var searchResults []map[string]interface{}
	for i, result := range results {
		if i >= 5 { // Limit to top 5 results
			break
		}
		searchResults = append(searchResults, map[string]interface{}{
			"service":     string(result.Service),
			"section":     result.Section,
			"title":       result.Title,
			"description": result.Description,
			"steps":       result.Steps,
			"error_codes": result.ErrorCodes,
			"admin_links": result.AdminLinks,
		})
	}

	return MCPResponse{
		Jsonrpc: "2.0",
		ID:      id,
		Result: map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": s.formatSearchResults(query, searchResults),
				},
			},
		},
	}
}

func (s *MCPServer) getServiceDependencies(id interface{}, args map[string]interface{}) MCPResponse {
	serviceName, _ := args["service"].(string)
	if serviceName == "" {
		return s.errorResponse(id, "Service name is required")
	}

	serviceType := models.ServiceType(serviceName)
	deps := s.knowledgeBase.GetServiceDependencies(serviceType)

	if deps == nil {
		return MCPResponse{
			Jsonrpc: "2.0",
			ID:      id,
			Result: map[string]interface{}{
				"content": []map[string]interface{}{
					{
						"type": "text",
						"text": fmt.Sprintf("No dependency information found for service: %s", serviceName),
					},
				},
			},
		}
	}

	result := map[string]interface{}{
		"service":       string(deps.Service),
		"upstream":      deps.Upstream,
		"downstream":    deps.Downstream,
		"critical_path": deps.CriticalPath,
	}

	return MCPResponse{
		Jsonrpc: "2.0",
		ID:      id,
		Result: map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": s.formatDependencies(result),
				},
			},
		},
	}
}

func (s *MCPServer) listServices(id interface{}) MCPResponse {
	overview := s.knowledgeBase.GetServiceOverview()

	var services []map[string]interface{}
	for service, count := range overview {
		services = append(services, map[string]interface{}{
			"service":          string(service),
			"runbook_sections": count,
		})
	}

	return MCPResponse{
		Jsonrpc: "2.0",
		ID:      id,
		Result: map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": s.formatServiceOverview(services),
				},
			},
		},
	}
}

func (s *MCPServer) errorResponse(id interface{}, message string) MCPResponse {
	return MCPResponse{
		Jsonrpc: "2.0",
		ID:      id,
		Error: map[string]interface{}{
			"code":    -32000,
			"message": message,
		},
	}
}

func (s *MCPServer) formatIncidentAnalysis(result map[string]interface{}) string {
	analysis := result["incident_analysis"].(map[string]interface{})

	output := fmt.Sprintf(`# Banking Incident Analysis

**Service:** %s
**Severity:** %s
**Error Code:** %s
**Description:** %s

`, analysis["service"], analysis["severity"], analysis["error_code"], analysis["description"])

	if runbook, ok := result["runbook_guidance"]; ok {
		rb := runbook.(map[string]interface{})
		output += fmt.Sprintf(`## Runbook Guidance

**Section %s:** %s

%s

### Resolution Steps:
`, rb["section"], rb["title"], rb["description"])

		if steps, ok := rb["steps"].([]string); ok {
			for i, step := range steps {
				output += fmt.Sprintf("%d. %s\n", i+1, step)
			}
		}

		if adminLinks, ok := rb["admin_links"].([]string); ok && len(adminLinks) > 0 {
			output += "\n### Admin Consoles:\n"
			for _, link := range adminLinks {
				output += fmt.Sprintf("- %s\n", link)
			}
		}
	}

	if deps, ok := result["service_dependencies"]; ok {
		d := deps.(map[string]interface{})
		output += "\n## Service Dependencies\n\n"

		if upstream, ok := d["upstream"].([]models.ServiceType); ok && len(upstream) > 0 {
			output += "**Upstream Services:**\n"
			for _, svc := range upstream {
				output += fmt.Sprintf("- %s\n", svc)
			}
		}

		if downstream, ok := d["downstream"].([]models.ServiceType); ok && len(downstream) > 0 {
			output += "\n**Downstream Impact:**\n"
			for _, svc := range downstream {
				output += fmt.Sprintf("- %s\n", svc)
			}
		}
	}

	if related, ok := result["related_procedures"]; ok {
		procedures := related.([]map[string]interface{})
		if len(procedures) > 0 {
			output += "\n## Related Procedures\n\n"
			for _, proc := range procedures {
				output += fmt.Sprintf("**%s** (%s): %s\n\n", proc["title"], proc["service"], proc["description"])
			}
		}
	}

	return output
}

func (s *MCPServer) formatSearchResults(query string, results []map[string]interface{}) string {
	output := fmt.Sprintf("# Runbook Search Results for: \"%s\"\n\n", query)

	if len(results) == 0 {
		return output + "No matching procedures found."
	}

	for i, result := range results {
		output += fmt.Sprintf("## %d. %s (%s)\n\n", i+1, result["title"], result["service"])
		output += fmt.Sprintf("%s\n\n", result["description"])

		if steps, ok := result["steps"].([]string); ok && len(steps) > 0 {
			output += "**Steps:**\n"
			for _, step := range steps {
				output += fmt.Sprintf("- %s\n", step)
			}
		}

		if errorCodes, ok := result["error_codes"].([]string); ok && len(errorCodes) > 0 {
			output += fmt.Sprintf("\n**Error Codes:** %s\n", fmt.Sprintf("%v", errorCodes))
		}

		output += "\n---\n\n"
	}

	return output
}

func (s *MCPServer) formatDependencies(deps map[string]interface{}) string {
	output := fmt.Sprintf("# Service Dependencies: %s\n\n", deps["service"])

	if upstream, ok := deps["upstream"].([]models.ServiceType); ok && len(upstream) > 0 {
		output += "## Upstream Dependencies\n"
		output += "These services must be working for this service to function:\n\n"
		for _, svc := range upstream {
			output += fmt.Sprintf("- %s\n", svc)
		}
		output += "\n"
	}

	if downstream, ok := deps["downstream"].([]models.ServiceType); ok && len(downstream) > 0 {
		output += "## Downstream Impact\n"
		output += "These services will be affected if this service fails:\n\n"
		for _, svc := range downstream {
			output += fmt.Sprintf("- %s\n", svc)
		}
		output += "\n"
	}

	if critical, ok := deps["critical_path"].(bool); ok && critical {
		output += "⚠️ **Critical Path Service** - High impact on overall system\n"
	}

	return output
}

func (s *MCPServer) formatServiceOverview(services []map[string]interface{}) string {
	output := "# Banking Services Overview\n\n"

	for _, svc := range services {
		output += fmt.Sprintf("## %s\n", svc["service"])
		output += fmt.Sprintf("- **Runbook Sections:** %v\n\n", svc["runbook_sections"])
	}

	return output
}
