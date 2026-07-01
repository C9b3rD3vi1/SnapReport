package templates

type Template struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	Category        string   `json:"category"`
	DefaultTitle    string   `json:"default_title"`
	Classification  string   `json:"classification"`
	Categories      []string `json:"categories"`
	Priorities      []string `json:"priorities"`
	Statuses        []string `json:"statuses"`
	DefaultSeverity string   `json:"default_severity"`
}

var All = []Template{
	{
		ID: "qa-bug-report", Name: "QA Bug Report", Category: "Quality Assurance",
		Description: "Detailed bug report with screenshots, environment info, and severity ratings.",
		DefaultTitle: "QA Bug Report", Classification: "Internal",
		Categories: []string{"Bug", "UI", "Backend", "Performance", "Security", "Enhancement"},
		Priorities: []string{"Critical", "High", "Medium", "Low"},
		Statuses:   []string{"Open", "In Progress", "Resolved", "Closed"},
		DefaultSeverity: "Medium",
	},
	{
		ID: "client-report", Name: "Client Report", Category: "Consulting",
		Description: "Professional consulting report suitable for direct client delivery.",
		DefaultTitle: "Technical Assessment Report", Classification: "Confidential",
		Categories: []string{"Finding", "Observation", "Recommendation", "Compliance"},
		Priorities: []string{"Critical", "High", "Medium", "Low"},
		Statuses:   []string{"Open", "Resolved"},
		DefaultSeverity: "Medium",
	},
	{
		ID: "technical-documentation", Name: "Technical Documentation", Category: "Engineering",
		Description: "System documentation with architecture screenshots and configuration details.",
		DefaultTitle: "Technical Documentation", Classification: "Internal",
		Categories: []string{"Architecture", "Configuration", "Deployment", "Integration", "API"},
		Priorities: []string{"High", "Medium", "Low"},
		Statuses:   []string{"Draft", "Review", "Approved"},
		DefaultSeverity: "Low",
	},
	{
		ID: "security-assessment", Name: "Security Assessment", Category: "Security",
		Description: "Security audit report with vulnerability findings and remediation steps.",
		DefaultTitle: "Security Assessment Report", Classification: "Confidential",
		Categories: []string{"Authentication", "Authorization", "Encryption", "Network", "Input Validation", "Configuration"},
		Priorities: []string{"Critical", "High", "Medium", "Low"},
		Statuses:   []string{"Open", "In Progress", "Resolved", "Accepted"},
		DefaultSeverity: "High",
	},
	{
		ID: "network-audit", Name: "Network Audit", Category: "Infrastructure",
		Description: "Network infrastructure audit with topology screenshots and compliance checks.",
		DefaultTitle: "Network Audit Report", Classification: "Confidential",
		Categories: []string{"Topology", "Firewall", "Routing", "DNS", "Load Balancing", "Monitoring"},
		Priorities: []string{"Critical", "High", "Medium", "Low"},
		Statuses:   []string{"Pass", "Fail", "Requires Action"},
		DefaultSeverity: "Medium",
	},
	{
		ID: "incident-report", Name: "Incident Report", Category: "Operations",
		Description: "Post-incident analysis with timeline, screenshots, and remediation steps.",
		DefaultTitle: "Incident Report", Classification: "Confidential",
		Categories: []string{"Detection", "Response", "Recovery", "Root Cause", "Prevention"},
		Priorities: []string{"Critical", "High", "Medium", "Low"},
		Statuses:   []string{"Investigating", "Resolved", "Post-Mortem"},
		DefaultSeverity: "Critical",
	},
}

func FindByID(id string) *Template {
	for _, t := range All {
		if t.ID == id {
			return &t
		}
	}
	return nil
}
