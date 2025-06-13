package main

import (
	"flag"
	"os"
	"strings"
	"testing"
)

// Helper function to reset flags for testing
func resetFlags() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
}

// Helper function to create a temporary test file
func createTestFile(t *testing.T, content string) string {
	tmpFile, err := os.CreateTemp("", "test_alerts_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	
	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}
	
	return tmpFile.Name()
}

// Test JSON content
const testJSONContent = `{
  "alerts": [
    {
      "id": "ALT-001",
      "timestamp": "2024-04-28T10:26:19Z",
      "service": "test-service",
      "component": "test-component",
      "severity": "critical",
      "metric": "latency",
      "value": 2300,
      "threshold": 1000,
      "description": "Test alert description"
    },
    {
      "id": "ALT-002",
      "timestamp": "2024-04-28T10:27:19Z",
      "service": "test-service-2",
      "component": "test-component-2",
      "severity": "warning",
      "metric": "cpu_usage",
      "value": 85,
      "threshold": 80,
      "description": "Test warning alert"
    }
  ]
}`

func TestParseFlags_ValidInput(t *testing.T) {
	resetFlags()
	
	// Create temporary test file
	testFile := createTestFile(t, testJSONContent)
	defer os.Remove(testFile)
	
	// Set up command line args
	os.Args = []string{"enc-alertbuddy", "-i", testFile}
	
	config, err := parseFlags()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	if config.InputFile != testFile {
		t.Errorf("Expected InputFile to be %s, got %s", testFile, config.InputFile)
	}
	
	if config.GroupBy != "" {
		t.Errorf("Expected GroupBy to be empty, got %s", config.GroupBy)
	}
	
	if config.LastMinutes != 0 {
		t.Errorf("Expected LastMinutes to be 0, got %d", config.LastMinutes)
	}
	
	if config.ShowAll {
		t.Errorf("Expected ShowAll to be false, got %t", config.ShowAll)
	}
}

func TestParseFlags_WithGroupBy(t *testing.T) {
	resetFlags()
	
	testFile := createTestFile(t, testJSONContent)
	defer os.Remove(testFile)
	
	os.Args = []string{"enc-alertbuddy", "-i", testFile, "--groupby=severity"}
	
	config, err := parseFlags()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	if config.GroupBy != "severity" {
		t.Errorf("Expected GroupBy to be 'severity', got %s", config.GroupBy)
	}
}

func TestParseFlags_WithLastMinutes(t *testing.T) {
	resetFlags()
	
	testFile := createTestFile(t, testJSONContent)
	defer os.Remove(testFile)
	
	os.Args = []string{"enc-alertbuddy", "-i", testFile, "--lastminutes=30"}
	
	config, err := parseFlags()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	if config.LastMinutes != 30 {
		t.Errorf("Expected LastMinutes to be 30, got %d", config.LastMinutes)
	}
}

func TestParseFlags_WithShowAll(t *testing.T) {
	resetFlags()
	
	testFile := createTestFile(t, testJSONContent)
	defer os.Remove(testFile)
	
	// Test --show-all flag
	os.Args = []string{"enc-alertbuddy", "-i", testFile, "--show-all"}
	
	config, err := parseFlags()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	if !config.ShowAll {
		t.Errorf("Expected ShowAll to be true, got %t", config.ShowAll)
	}
}

func TestParseFlags_WithShowAllShortFlag(t *testing.T) {
	resetFlags()
	
	testFile := createTestFile(t, testJSONContent)
	defer os.Remove(testFile)
	
	// Test -a flag
	os.Args = []string{"enc-alertbuddy", "-i", testFile, "-a"}
	
	config, err := parseFlags()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	if !config.ShowAll {
		t.Errorf("Expected ShowAll to be true with -a flag, got %t", config.ShowAll)
	}
}

func TestParseFlags_CombinedFlags(t *testing.T) {
	resetFlags()
	
	testFile := createTestFile(t, testJSONContent)
	defer os.Remove(testFile)
	
	os.Args = []string{"enc-alertbuddy", "-i", testFile, "--groupby=service", "--lastminutes=60"}
	
	config, err := parseFlags()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	if config.GroupBy != "service" {
		t.Errorf("Expected GroupBy to be 'service', got %s", config.GroupBy)
	}
	
	if config.LastMinutes != 60 {
		t.Errorf("Expected LastMinutes to be 60, got %d", config.LastMinutes)
	}
}

func TestParseFlags_MissingInputFile(t *testing.T) {
	resetFlags()
	
	os.Args = []string{"enc-alertbuddy", "--groupby=severity"}
	
	_, err := parseFlags()
	if err == nil {
		t.Error("Expected error for missing input file, got nil")
	}
	
	if !strings.Contains(err.Error(), "input file is required") {
		t.Errorf("Expected error message about missing input file, got: %v", err)
	}
}

func TestParseFlags_NonExistentFile(t *testing.T) {
	resetFlags()
	
	os.Args = []string{"enc-alertbuddy", "-i", "nonexistent_file.json"}
	
	_, err := parseFlags()
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
	
	if !strings.Contains(err.Error(), "does not exist") {
		t.Errorf("Expected error message about non-existent file, got: %v", err)
	}
}

func TestParseFlags_InvalidGroupByField(t *testing.T) {
	resetFlags()
	
	testFile := createTestFile(t, testJSONContent)
	defer os.Remove(testFile)
	
	os.Args = []string{"enc-alertbuddy", "-i", testFile, "--groupby=invalid_field"}
	
	_, err := parseFlags()
	if err == nil {
		t.Error("Expected error for invalid groupby field, got nil")
	}
	
	if !strings.Contains(err.Error(), "invalid groupby field") {
		t.Errorf("Expected error message about invalid groupby field, got: %v", err)
	}
}

func TestParseFlags_NegativeLastMinutes(t *testing.T) {
	resetFlags()
	
	testFile := createTestFile(t, testJSONContent)
	defer os.Remove(testFile)
	
	os.Args = []string{"enc-alertbuddy", "-i", testFile, "--lastminutes=-10"}
	
	_, err := parseFlags()
	if err == nil {
		t.Error("Expected error for negative lastminutes, got nil")
	}
	
	if !strings.Contains(err.Error(), "must be a positive number") {
		t.Errorf("Expected error message about positive number, got: %v", err)
	}
}

func TestParseFlags_ValidGroupByFields(t *testing.T) {
	testFile := createTestFile(t, testJSONContent)
	defer os.Remove(testFile)
	
	validFields := []string{"severity", "service", "component", "metric", "threshold", "value", "priority"}
	
	for _, field := range validFields {
		resetFlags()
		os.Args = []string{"enc-alertbuddy", "-i", testFile, "--groupby=" + field}
		
		config, err := parseFlags()
		if err != nil {
			t.Errorf("Expected no error for valid field '%s', got: %v", field, err)
		}
		
		if config.GroupBy != field {
			t.Errorf("Expected GroupBy to be '%s', got %s", field, config.GroupBy)
		}
	}
}

func TestLoadAlertsFromFile_ValidFile(t *testing.T) {
	testFile := createTestFile(t, testJSONContent)
	defer os.Remove(testFile)
	
	alerts, err := loadAlertsFromFile(testFile)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	if alerts == nil {
		t.Error("Expected alerts to not be nil")
	}
	
	if len(alerts.Alerts) != 2 {
		t.Errorf("Expected 2 alerts, got %d", len(alerts.Alerts))
	}
	
	// Check first alert
	if alerts.Alerts[0].ID != "ALT-001" {
		t.Errorf("Expected first alert ID to be ALT-001, got %s", alerts.Alerts[0].ID)
	}
	
	if alerts.Alerts[0].Severity != "critical" {
		t.Errorf("Expected first alert severity to be critical, got %s", alerts.Alerts[0].Severity)
	}
}

func TestLoadAlertsFromFile_NonExistentFile(t *testing.T) {
	_, err := loadAlertsFromFile("nonexistent_file.json")
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
	
	if !strings.Contains(err.Error(), "error reading file") {
		t.Errorf("Expected error message about reading file, got: %v", err)
	}
}

func TestLoadAlertsFromFile_InvalidJSON(t *testing.T) {
	invalidJSON := `{"alerts": [{"id": "test", "invalid": json}]}`
	testFile := createTestFile(t, invalidJSON)
	defer os.Remove(testFile)
	
	_, err := loadAlertsFromFile(testFile)
	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
	
	if !strings.Contains(err.Error(), "error parsing JSON") {
		t.Errorf("Expected error message about parsing JSON, got: %v", err)
	}
}

func TestLoadAlertsFromFile_EmptyAlerts(t *testing.T) {
	emptyJSON := `{"alerts": []}`
	testFile := createTestFile(t, emptyJSON)
	defer os.Remove(testFile)
	
	_, err := loadAlertsFromFile(testFile)
	if err == nil {
		t.Error("Expected error for empty alerts, got nil")
	}
	
	if !strings.Contains(err.Error(), "no alerts found") {
		t.Errorf("Expected error message about no alerts found, got: %v", err)
	}
}

func TestContains(t *testing.T) {
	slice := []string{"severity", "service", "component"}
	
	// Test existing item
	if !contains(slice, "severity") {
		t.Error("Expected contains to return true for 'severity'")
	}
	
	// Test case insensitive
	if !contains(slice, "SEVERITY") {
		t.Error("Expected contains to return true for 'SEVERITY' (case insensitive)")
	}
	
	// Test non-existing item
	if contains(slice, "nonexistent") {
		t.Error("Expected contains to return false for 'nonexistent'")
	}
	
	// Test empty slice
	emptySlice := []string{}
	if contains(emptySlice, "severity") {
		t.Error("Expected contains to return false for empty slice")
	}
}

// Integration test for the entire process
func TestProcessAlerts_Integration(t *testing.T) {
	testFile := createTestFile(t, testJSONContent)
	defer os.Remove(testFile)
	
	// Load alerts
	alerts, err := loadAlertsFromFile(testFile)
	if err != nil {
		t.Fatalf("Failed to load alerts: %v", err)
	}
	
	// Test with basic config
	config := &Config{
		InputFile:   testFile,
		GroupBy:     "",
		LastMinutes: 0,
		ShowAll:     false,
	}
	
	// This should not panic or error
	// Note: In real tests, you might want to capture output
	processAlerts(alerts, config)
	
	// Test with show-all
	config.ShowAll = true
	processAlerts(alerts, config)
	
	// Test with groupby
	config.ShowAll = false
	config.GroupBy = "severity"
	processAlerts(alerts, config)
}

// Benchmark test for CLI parsing
func BenchmarkParseFlags(b *testing.B) {
	testFile := createTestFile(&testing.T{}, testJSONContent)
	defer os.Remove(testFile)
	
	for i := 0; i < b.N; i++ {
		resetFlags()
		os.Args = []string{"enc-alertbuddy", "-i", testFile, "--groupby=severity", "--lastminutes=30"}
		
		_, err := parseFlags()
		if err != nil {
			b.Fatalf("Unexpected error: %v", err)
		}
	}
}

// Benchmark test for loading alerts
func BenchmarkLoadAlertsFromFile(b *testing.B) {
	testFile := createTestFile(&testing.T{}, testJSONContent)
	defer os.Remove(testFile)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := loadAlertsFromFile(testFile)
		if err != nil {
			b.Fatalf("Unexpected error: %v", err)
		}
	}
}