package main
import "time"
// Alerts represents a collection of alerts
type Alerts struct {
	Alerts []Alert `json:"alerts"`
}

// Alert represents a single alert with all its properties
type Alert struct {
	ID          string    `json:"id"`
	Timestamp   time.Time `json:"timestamp"`
	Service     string    `json:"service"`
	Component   string    `json:"component"`
	Severity    string    `json:"severity"`
	Metric      string    `json:"metric"`
	Value       float64   `json:"value"`
	Threshold   float64   `json:"threshold"`
	Description string    `json:"description"`
}
