package main

import (
	"fmt"
	"strings"
)

func (alert Alert) PrettyPrint() {
	fmt.Printf("┌─ Alert: %s ─┐\n", alert.ID)
	fmt.Printf("│ Service:     %s\n", alert.Service)
	fmt.Printf("│ Component:   %s\n", alert.Component)
	fmt.Printf("│ Severity:    %s\n", alert.Severity)
	fmt.Printf("│ Metric:      %s\n", alert.Metric)
	fmt.Printf("│ Value:       %.2f (threshold: %.2f)\n", alert.Value, alert.Threshold)
	fmt.Printf("│ Time:        %s\n", alert.Timestamp.Format("2006-01-02 15:04:05"))
	fmt.Printf("│ Description: %s\n", alert.Description)
	fmt.Println("└────────────────────────────────────────┘")
}

// PrettyPrint formats and prints all alerts in the collection
func (alerts Alerts) PrettyPrint() {
	fmt.Printf("📋 Alerts Summary: %d total\n", len(alerts.Alerts))
	fmt.Println(strings.Repeat("=", 50))

	if len(alerts.Alerts) == 0 {
		fmt.Println("No alerts to display.")
		return
	}

	for i, alert := range alerts.Alerts {
		fmt.Printf("\n[%d/%d]\n", i+1, len(alerts.Alerts))
		alert.PrettyPrint()
	}

	fmt.Printf("\n📊 Total: %d alerts displayed\n", len(alerts.Alerts))
}
