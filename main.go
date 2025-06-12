package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func main() {
	// Read the JSON file
	data, err := os.ReadFile("sample-alerts.json")
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	// Unmarshal JSON into alerts struct
	var alerts Alerts
	err = json.Unmarshal(data, &alerts)
	if err != nil {
		log.Fatalf("Error unmarshaling JSON: %v", err)
	}

	// Print the alerts
	fmt.Printf("Loaded %d alerts:\n\n", len(alerts.Alerts))

	for _, alert := range alerts.Alerts {
		fmt.Printf("ID: %s\n", alert.ID)
		fmt.Printf("Timestamp: %s\n", alert.Timestamp.Format("2006-01-02 15:04:05"))
		fmt.Printf("Service: %s\n", alert.Service)
		fmt.Printf("Component: %s\n", alert.Component)
		fmt.Printf("Severity: %s\n", alert.Severity)
		fmt.Printf("Metric: %s\n", alert.Metric)
		fmt.Printf("Value: %.2f\n", alert.Value)
		fmt.Printf("Threshold: %.2f\n", alert.Threshold)
		fmt.Printf("Description: %s\n", alert.Description)
		fmt.Println("---")
	}
}
