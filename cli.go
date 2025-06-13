package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
)

const (
	Version = "1.0.0"
	AppName = "enc-alertbuddy"
)

type Config struct {
	InputFile   string
	GroupBy     string
	LastMinutes int
	ShowVersion bool
	ShowHelp    bool
	ShowAll     bool
}

func parseFlags() (*Config, error) {
	config := &Config{}
	
	// Define flags
	flag.StringVar(&config.InputFile, "i", "", "Input JSON file containing alerts")
	flag.StringVar(&config.InputFile, "input", "", "Input JSON file containing alerts")
	flag.StringVar(&config.GroupBy, "groupby", "", "Group alerts by field (severity, service, component, metric, etc.)")
	flag.IntVar(&config.LastMinutes, "lastminutes", 0, "Filter alerts from the last N minutes")
	flag.BoolVar(&config.ShowAll, "show-all", false, "Show all alerts in detailed format")
	flag.BoolVar(&config.ShowAll, "a", false, "Show all alerts in detailed format")
	flag.BoolVar(&config.ShowVersion, "version", false, "Show version information")
	flag.BoolVar(&config.ShowVersion, "v", false, "Show version information")
	flag.BoolVar(&config.ShowHelp, "help", false, "Show help information")
	flag.BoolVar(&config.ShowHelp, "h", false, "Show help information")
	
	// Custom usage function
	flag.Usage = func() {
		showHelp()
	}
	
	// Parse flags
	flag.Parse()
	
	// Handle version
	if config.ShowVersion {
		showVersion()
		os.Exit(0)
	}
	
	// Handle help
	if config.ShowHelp {
		showHelp()
		os.Exit(0)
	}
	
	// Validate required flags
	if config.InputFile == "" {
		return nil, fmt.Errorf("input file is required. Use -i or --input to specify the JSON file")
	}
	
	// Check if file exists
	if _, err := os.Stat(config.InputFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("input file '%s' does not exist", config.InputFile)
	}
	
	// Validate groupby field if provided
	if config.GroupBy != "" {
		validFields := []string{"severity", "service", "component", "metric", "threshold", "value", "priority"}
		if !contains(validFields, strings.ToLower(config.GroupBy)) {
			return nil, fmt.Errorf("invalid groupby field '%s'. Valid fields: %s", 
				config.GroupBy, strings.Join(validFields, ", "))
		}
	}
	
	// Validate lastminutes if provided
	if config.LastMinutes < 0 {
		return nil, fmt.Errorf("lastminutes must be a positive number")
	}
	
	return config, nil
}

func showVersion() {
	fmt.Printf("%s version %s\n", AppName, Version)
	fmt.Println("A powerful alert management and analysis tool")
}

func showHelp() {
	fmt.Printf("%s - Alert Management CLI Tool\n\n", AppName)
	fmt.Println("USAGE:")
	fmt.Printf("  %s -i <input-file> [OPTIONS]\n\n", AppName)
	
	fmt.Println("REQUIRED FLAGS:")
	fmt.Println("  -i, --input <file>     Input JSON file containing alerts")
	
	fmt.Println("\nOPTIONAL FLAGS:")
	fmt.Println("  --groupby <field>      Group alerts by field (severity, service, component, metric, threshold, value, priority)")
	fmt.Println("  --lastminutes <n>      Filter alerts from the last N minutes")
	fmt.Println("  --show-all, -a         Show all alerts in detailed format")
	fmt.Println("  -v, --version          Show version information")
	fmt.Println("  -h, --help             Show this help message")
	
	fmt.Println("\nEXAMPLES:")
	fmt.Printf("  %s -i alerts.json\n", AppName)
	fmt.Printf("  %s -i alerts.json --show-all\n", AppName)
	fmt.Printf("  %s -i alerts.json -a\n", AppName)
	fmt.Printf("  %s -i alerts.json --groupby=severity\n", AppName)
	fmt.Printf("  %s -i alerts.json --groupby severity\n", AppName)
	fmt.Printf("  %s -i alerts.json --lastminutes=30\n", AppName)
	fmt.Printf("  %s -i alerts.json --groupby=service --lastminutes=60\n", AppName)
	fmt.Printf("  %s -i alerts.json --show-all --lastminutes=30\n", AppName)
	
	fmt.Println("\nVALID GROUPBY FIELDS:")
	fmt.Println("  severity    - Group by alert severity (critical, warning, info)")
	fmt.Println("  service     - Group by service name")
	fmt.Println("  component   - Group by component name")
	fmt.Println("  metric      - Group by metric type")
	fmt.Println("  threshold   - Group by threshold value")
	fmt.Println("  value       - Group by metric value")
	fmt.Println("  priority    - Group by calculated priority score")
	
	fmt.Println("\nFEATURES:")
	fmt.Println("  • Automatic priority calculation based on severity, deviation, and affected components")
	fmt.Println("  • Time-based filtering to focus on recent alerts")
	fmt.Println("  • Flexible grouping by any alert field")
	fmt.Println("  • Show all alerts in detailed format with --show-all")
	fmt.Println("  • Beautiful formatted output for better readability")
	fmt.Println("  • Support for large alert datasets")
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if strings.ToLower(s) == strings.ToLower(item) {
			return true
		}
	}
	return false
}

func loadAlertsFromFile(filename string) (*Alerts, error) {
	// Read the JSON file
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading file '%s': %v", filename, err)
	}
	
	// Unmarshal JSON into alerts struct
	var alerts Alerts
	err = json.Unmarshal(data, &alerts)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON from '%s': %v", filename, err)
	}
	
	if len(alerts.Alerts) == 0 {
		return nil, fmt.Errorf("no alerts found in file '%s'", filename)
	}
	
	return &alerts, nil
}

func processAlerts(alerts *Alerts, config *Config) {
	// Calculate priorities for all alerts
	alerts.CalculateAllPriorities()
	
	// Sort by priority (highest first)
	alerts.SortByPriority()
	
	// Apply time filter if specified
	if config.LastMinutes > 0 {
		fmt.Printf("🕒 Filtering alerts from the last %d minutes...\n", config.LastMinutes)
		filtered := alerts.FilterByLastMinutes(config.LastMinutes)
		alerts = &filtered
		
		if len(alerts.Alerts) == 0 {
			fmt.Printf("⚠️  No alerts found in the last %d minutes.\n", config.LastMinutes)
			return
		}
	}
	
	// Show summary
	fmt.Printf("📊 Loaded %d alerts from %s\n", len(alerts.Alerts), config.InputFile)
	if config.LastMinutes > 0 {
		fmt.Printf("🔍 Showing alerts from the last %d minutes\n", config.LastMinutes)
	}
	fmt.Println()
	
	// Group and display if groupby is specified
	if config.GroupBy != "" {
		fmt.Printf("📋 Grouping alerts by: %s\n", config.GroupBy)
		alerts.PrettyPrintGroupedBy(config.GroupBy)
	} else if config.ShowAll {
		// Show all alerts in detailed format
		fmt.Printf("📋 Showing all %d alerts in detailed format:\n", len(alerts.Alerts))
		alerts.PrettyPrint()
	} else {
		// Show top 10 highest priority alerts
		maxDisplay := 10
		if len(alerts.Alerts) < maxDisplay {
			maxDisplay = len(alerts.Alerts)
		}
		
		fmt.Printf("🔥 Top %d Highest Priority Alerts:\n", maxDisplay)
		fmt.Println(strings.Repeat("=", 60))
		
		for i := 0; i < maxDisplay; i++ {
			alert := alerts.Alerts[i]
			fmt.Printf("\n[%d] Priority: %.2f | %s | %s\n", 
				i+1, alert.Priority, alert.Severity, alert.ID)
			fmt.Printf("    Service: %s | Component: %s\n", 
				alert.Service, alert.Component)
			fmt.Printf("    Metric: %s (%.2f / %.2f)\n", 
				alert.Metric, alert.Value, alert.Threshold)
			fmt.Printf("    Description: %s\n", alert.Description)
		}
		
		if len(alerts.Alerts) > maxDisplay {
			fmt.Printf("\n... and %d more alerts (use --show-all to see all or --groupby to organize)\n", 
				len(alerts.Alerts)-maxDisplay)
		}
		
		// Show summary statistics
		showSummaryStats(alerts)
	}
}

func showSummaryStats(alerts *Alerts) {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("📈 SUMMARY STATISTICS")
	fmt.Println(strings.Repeat("=", 60))
	
	// Count by severity
	severityCounts := make(map[string]int)
	serviceCounts := make(map[string]int)
	var totalPriority float64
	
	for _, alert := range alerts.Alerts {
		severityCounts[alert.Severity]++
		serviceCounts[alert.Service]++
		totalPriority += alert.Priority
	}
	
	// Severity breakdown
	fmt.Println("\n🚨 Severity Breakdown:")
	for _, severity := range []string{"critical", "warning", "info"} {
		if count, exists := severityCounts[severity]; exists {
			percentage := float64(count) / float64(len(alerts.Alerts)) * 100
			fmt.Printf("  %s: %d alerts (%.1f%%)\n", 
				strings.Title(severity), count, percentage)
		}
	}
	
	// Average priority
	avgPriority := totalPriority / float64(len(alerts.Alerts))
	fmt.Printf("\n⚡ Average Priority Score: %.2f\n", avgPriority)
	
	// Top services
	fmt.Printf("\n🏢 Top 5 Services by Alert Count:\n")
	// Simple approach - show first 5 services
	count := 0
	for service, alertCount := range serviceCounts {
		if count < 5 {
			fmt.Printf("  %s: %d alerts\n", service, alertCount)
			count++
		}
	}
}

func handleCLIError(err error) {
	fmt.Fprintf(os.Stderr, "❌ Error: %v\n\n", err)
	showHelp()
	os.Exit(1)
}

func runCLI() {
	config, err := parseFlags()
	if err != nil {
		handleCLIError(err)
	}
	
	// Load alerts from file
	alerts, err := loadAlertsFromFile(config.InputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
		os.Exit(1)
	}
	
	// Process and display alerts
	processAlerts(alerts, config)
}