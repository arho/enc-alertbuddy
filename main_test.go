package main

import (
	"testing"
	"time"
)

// Helper function to create test alerts
func createTestAlerts() Alerts {
	return Alerts{
		Alerts: []Alert{
			{
				ID:          "ALT-001",
				Timestamp:   time.Now().Add(-10 * time.Minute),
				Service:     "payment-processor",
				Component:   "api-gateway",
				Severity:    "critical",
				Metric:      "latency",
				Value:       2300,
				Threshold:   1000,
				Description: "API response time exceeded threshold",
			},
			{
				ID:          "ALT-002",
				Timestamp:   time.Now().Add(-5 * time.Minute),
				Service:     "payment-processor",
				Component:   "database",
				Severity:    "warning",
				Metric:      "cpu_usage",
				Value:       85,
				Threshold:   80,
				Description: "Database CPU usage approaching critical threshold",
			},
			{
				ID:          "ALT-003",
				Timestamp:   time.Now().Add(-60 * time.Minute),
				Service:     "user-authentication",
				Component:   "auth-service",
				Severity:    "info",
				Metric:      "memory_usage",
				Value:       68.5,
				Threshold:   70.0,
				Description: "Memory usage within normal range",
			},
			{
				ID:          "ALT-004",
				Timestamp:   time.Now().Add(-2 * time.Minute),
				Service:     "user-authentication",
				Component:   "session-manager",
				Severity:    "critical",
				Metric:      "latency",
				Value:       5200,
				Threshold:   2000,
				Description: "Session management latency critically high",
			},
		},
	}
}

func TestFilterBySeverity(t *testing.T) {
	alerts := createTestAlerts()
	
	// Test filtering by critical severity
	criticalAlerts := alerts.FilterBySeverity("critical")
	if len(criticalAlerts.Alerts) != 2 {
		t.Errorf("Expected 2 critical alerts, got %d", len(criticalAlerts.Alerts))
	}
	
	// Test filtering by warning severity
	warningAlerts := alerts.FilterBySeverity("warning")
	if len(warningAlerts.Alerts) != 1 {
		t.Errorf("Expected 1 warning alert, got %d", len(warningAlerts.Alerts))
	}
	
	// Test filtering by non-existent severity
	noneAlerts := alerts.FilterBySeverity("nonexistent")
	if len(noneAlerts.Alerts) != 0 {
		t.Errorf("Expected 0 alerts for non-existent severity, got %d", len(noneAlerts.Alerts))
	}
}

func TestFilterByService(t *testing.T) {
	alerts := createTestAlerts()
	
	// Test filtering by payment-processor service
	paymentAlerts := alerts.FilterByService("payment-processor")
	if len(paymentAlerts.Alerts) != 2 {
		t.Errorf("Expected 2 payment-processor alerts, got %d", len(paymentAlerts.Alerts))
	}
	
	// Test filtering by user-authentication service
	authAlerts := alerts.FilterByService("user-authentication")
	if len(authAlerts.Alerts) != 2 {
		t.Errorf("Expected 2 user-authentication alerts, got %d", len(authAlerts.Alerts))
	}
	
	// Test filtering by non-existent service
	noneAlerts := alerts.FilterByService("nonexistent")
	if len(noneAlerts.Alerts) != 0 {
		t.Errorf("Expected 0 alerts for non-existent service, got %d", len(noneAlerts.Alerts))
	}
}

func TestFilterByLastMinutes(t *testing.T) {
	alerts := createTestAlerts()
	
	// Test filtering by last 15 minutes (should get 3 alerts)
	recent15 := alerts.FilterByLastMinutes(15)
	if len(recent15.Alerts) != 3 {
		t.Errorf("Expected 3 alerts in last 15 minutes, got %d", len(recent15.Alerts))
	}
	
	// Test filtering by last 3 minutes (should get 1 alert)
	recent3 := alerts.FilterByLastMinutes(3)
	if len(recent3.Alerts) != 1 {
		t.Errorf("Expected 1 alert in last 3 minutes, got %d", len(recent3.Alerts))
	}
	
	// Test filtering by last 120 minutes (should get all 4 alerts)
	recent120 := alerts.FilterByLastMinutes(120)
	if len(recent120.Alerts) != 4 {
		t.Errorf("Expected 4 alerts in last 120 minutes, got %d", len(recent120.Alerts))
	}
}

func TestGroup(t *testing.T) {
	alerts := createTestAlerts()
	
	// Test grouping by severity
	severityGroups := alerts.Group("severity")
	if len(severityGroups) != 3 {
		t.Errorf("Expected 3 severity groups, got %d", len(severityGroups))
	}
	
	if len(severityGroups["critical"].Alerts) != 2 {
		t.Errorf("Expected 2 critical alerts in group, got %d", len(severityGroups["critical"].Alerts))
	}
	
	// Test grouping by service
	serviceGroups := alerts.Group("service")
	if len(serviceGroups) != 2 {
		t.Errorf("Expected 2 service groups, got %d", len(serviceGroups))
	}
	
	if len(serviceGroups["payment-processor"].Alerts) != 2 {
		t.Errorf("Expected 2 payment-processor alerts in group, got %d", len(serviceGroups["payment-processor"].Alerts))
	}
	
	// Test grouping by non-existent field
	invalidGroups := alerts.Group("nonexistent")
	if len(invalidGroups) != 0 {
		t.Errorf("Expected 0 groups for non-existent field, got %d", len(invalidGroups))
	}
}

func TestCalculateDeviationPercentage(t *testing.T) {
	tests := []struct {
		value     float64
		threshold float64
		expected  float64
	}{
		{100, 80, 25.0},    // 25% deviation
		{80, 100, 20.0},    // 20% deviation
		{200, 100, 100.0},  // 100% deviation
		{100, 0, 0.0},      // Division by zero case
		{50, 50, 0.0},      // No deviation
	}
	
	for _, test := range tests {
		result := calculateDeviationPercentage(test.value, test.threshold)
		if result != test.expected {
			t.Errorf("calculateDeviationPercentage(%.1f, %.1f) = %.1f, expected %.1f", 
				test.value, test.threshold, result, test.expected)
		}
	}
}

func TestCountAffectedComponents(t *testing.T) {
	alerts := createTestAlerts()
	
	// Test counting affected components for latency metric
	targetAlert := alerts.Alerts[0] // payment-processor, latency
	count := alerts.countAffectedComponents(targetAlert)
	
	// Should count only itself = 1 (no other payment-processor + latency alerts)
	if count != 1 {
		t.Errorf("Expected 1 affected component for payment-processor latency metric, got %d", count)
	}
	
	// Test counting for unique metric
	targetAlert2 := alerts.Alerts[1] // payment-processor, cpu_usage
	count2 := alerts.countAffectedComponents(targetAlert2)
	
	// Should count only itself = 1 (no other payment-processor + cpu_usage alerts)
	if count2 != 1 {
		t.Errorf("Expected 1 affected component for cpu_usage metric, got %d", count2)
	}
	
	// Test with alerts that have matching service + metric
	// Add test alerts with same service and metric to verify counting works
	testAlerts := Alerts{
		Alerts: []Alert{
			{
				ID: "ALT-100", Service: "test-service", Component: "component-1", 
				Metric: "latency", Value: 100, Threshold: 50,
			},
			{
				ID: "ALT-101", Service: "test-service", Component: "component-2", 
				Metric: "latency", Value: 120, Threshold: 50,
			},
			{
				ID: "ALT-102", Service: "test-service", Component: "component-3", 
				Metric: "cpu_usage", Value: 80, Threshold: 70,
			},
		},
	}
	
	// Test first alert (should find 1 other + itself = 2)
	count3 := testAlerts.countAffectedComponents(testAlerts.Alerts[0])
	if count3 != 2 {
		t.Errorf("Expected 2 affected components for test-service latency, got %d", count3)
	}
	
	// Test third alert (should find 0 others + itself = 1)
	count4 := testAlerts.countAffectedComponents(testAlerts.Alerts[2])
	if count4 != 1 {
		t.Errorf("Expected 1 affected component for test-service cpu_usage, got %d", count4)
	}
}

func TestCalculatePriority(t *testing.T) {
	alerts := createTestAlerts()
	
	// Calculate priorities for all alerts
	alerts.CalculateAllPriorities()
	
	// Test that priorities are calculated (not zero)
	for i, alert := range alerts.Alerts {
		if alert.Priority == 0 {
			t.Errorf("Alert %d priority should not be zero", i)
		}
	}
	
	// Test that critical alerts have higher priority than warning/info
	criticalPriority := alerts.Alerts[0].Priority // critical alert
	warningPriority := alerts.Alerts[1].Priority  // warning alert
	infoPriority := alerts.Alerts[2].Priority     // info alert
	
	if criticalPriority <= warningPriority {
		t.Errorf("Critical priority (%.2f) should be higher than warning priority (%.2f)", 
			criticalPriority, warningPriority)
	}
	
	if warningPriority <= infoPriority {
		t.Errorf("Warning priority (%.2f) should be higher than info priority (%.2f)", 
			warningPriority, infoPriority)
	}
}

func TestSortByPriority(t *testing.T) {
	alerts := createTestAlerts()
	alerts.CalculateAllPriorities()
	
	// Get initial order
	firstAlertPriority := alerts.Alerts[0].Priority
	
	// Sort by priority
	alerts.SortByPriority()
	
	// Check that alerts are sorted in descending order
	for i := 0; i < len(alerts.Alerts)-1; i++ {
		current := alerts.Alerts[i].Priority
		next := alerts.Alerts[i+1].Priority
		
		if current < next {
			t.Errorf("Alerts not sorted correctly: position %d (%.2f) < position %d (%.2f)", 
				i, current, i+1, next)
		}
	}
	
	// The first alert after sorting should have the highest priority
	if alerts.Alerts[0].Priority < firstAlertPriority {
		t.Errorf("First alert after sorting should have highest priority")
	}
}

func TestAlertPrettyPrint(t *testing.T) {
	alert := Alert{
		ID:          "TEST-001",
		Service:     "test-service",
		Component:   "test-component",
		Severity:    "critical",
		Priority:    25.5,
		Metric:      "test-metric",
		Value:       100.0,
		Threshold:   50.0,
		Description: "Test alert description",
		Timestamp:   time.Now(),
	}
	
	// This test mainly ensures the function doesn't panic
	// In a real scenario, you might want to capture output and verify content
	alert.PrettyPrint()
}

func TestAlertsPrettyPrint(t *testing.T) {
	alerts := createTestAlerts()
	
	// This test mainly ensures the function doesn't panic
	// In a real scenario, you might want to capture output and verify content
	alerts.PrettyPrint()
	
	// Test with empty alerts
	emptyAlerts := Alerts{}
	emptyAlerts.PrettyPrint()
}

// Benchmark tests
func BenchmarkFilterBySeverity(b *testing.B) {
	alerts := createTestAlerts()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		alerts.FilterBySeverity("critical")
	}
}

func BenchmarkCalculateAllPriorities(b *testing.B) {
	alerts := createTestAlerts()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		alerts.CalculateAllPriorities()
	}
}

func BenchmarkGroup(b *testing.B) {
	alerts := createTestAlerts()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		alerts.Group("severity")
	}
}