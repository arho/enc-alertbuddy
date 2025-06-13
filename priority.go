package main
import (
	"math"
	"sort"
	"strings"
)

// calculateDeviationPercentage calculates the percentage deviation from threshold
func calculateDeviationPercentage(value, threshold float64) float64 {
	if threshold == 0 {
		return 0 // Avoid division by zero
	}

	deviation := math.Abs(value - threshold)
	percentage := (deviation / threshold) * 100

	// Cap the deviation at 1000% to avoid extreme scores
	if percentage > 1000 {
		percentage = 1000
	}

	return percentage
}

// countAffectedComponents counts unique components for the same service and metric
func (alerts Alerts) countAffectedComponents(targetAlert Alert) int {
	componentSet := make(map[string]bool)

	for _, alert := range alerts.Alerts {
		// Count components with same service and metric that are also alerting
		if alert.Service == targetAlert.Service &&
			alert.Metric == targetAlert.Metric &&
			alert.ID != targetAlert.ID { // Don't count itself
			componentSet[alert.Component] = true
		}
	}

	// Always count at least 1 (the alert itself)
	return len(componentSet) + 1
}

// CalculatePriority calculates and sets the priority score for an alert
func (alert *Alert) CalculatePriority(allAlerts Alerts) {
	// 1. Severity score (critical=10, warning=5, info=1)
	var severityScore float64
	switch strings.ToLower(alert.Severity) {
	case "critical":
		severityScore = 10.0
	case "warning":
		severityScore = 5.0
	case "info":
		severityScore = 1.0
	default:
		severityScore = 1.0 // Default to info level
	}

	// 2. Deviation from threshold (percentage)
	deviationPercentage := calculateDeviationPercentage(alert.Value, alert.Threshold)

	// 3. Number of affected components
	affectedComponents := float64(allAlerts.countAffectedComponents(*alert))

	// Calculate priority score using weighted formula
	// Priority = (Severity * 1.0) + (Deviation% * 0.1) + (Components * 2.0)
	priority := (severityScore * 1.0) + (deviationPercentage * 0.1) + (affectedComponents * 2.0)

	alert.Priority = math.Round(priority*100) / 100 // Round to 2 decimal places
}

// CalculateAllPriorities calculates priority scores for all alerts
func (alerts *Alerts) CalculateAllPriorities() {
	for i := range alerts.Alerts {
		alerts.Alerts[i].CalculatePriority(*alerts)
	}
}

// SortByPriority sorts alerts by priority score in descending order (highest first)
func (alerts *Alerts) SortByPriority() {
	sort.Slice(alerts.Alerts, func(i, j int) bool {
		return alerts.Alerts[i].Priority > alerts.Alerts[j].Priority
	})
}
