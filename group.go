package main

import (
	"fmt"
	"reflect"
	"strings"
)

// Group groups alerts by any field using reflection
func (alerts Alerts) Group(field string) map[string]Alerts {
	grouped := make(map[string]Alerts)

	for _, alert := range alerts.Alerts {
		var key string

		// Use reflection to get the field value
		alertValue := reflect.ValueOf(alert)
		fieldValue := alertValue.FieldByName(strings.Title(field))

		if !fieldValue.IsValid() {
			// If field doesn't exist, skip this alert
			continue
		}

		// Convert field value to string for use as map key
		switch fieldValue.Kind() {
		case reflect.String:
			key = fieldValue.String()
		case reflect.Float64, reflect.Float32:
			key = fmt.Sprintf("%.2f", fieldValue.Float())
		case reflect.Int, reflect.Int64, reflect.Int32:
			key = fmt.Sprintf("%d", fieldValue.Int())
		default:
			key = fmt.Sprintf("%v", fieldValue.Interface())
		}

		// Initialize group if it doesn't exist
		if _, exists := grouped[key]; !exists {
			grouped[key] = Alerts{Alerts: []Alert{}}
		}

		// Append alert to the group
		groupAlerts := grouped[key]
		groupAlerts.Alerts = append(groupAlerts.Alerts, alert)
		grouped[key] = groupAlerts
	}

	return grouped
}

// PrettyPrintGrouped prints grouped alerts in a nice format
func prettyPrintGrouped(grouped map[string]Alerts, groupName string) {
	fmt.Printf("📊 Alerts grouped by %s:\n", groupName)
	fmt.Println(strings.Repeat("=", 60))

	for key, alertGroup := range grouped {
		fmt.Printf("\n🏷️  %s: %d alerts\n", key, len(alertGroup.Alerts))
		fmt.Println(strings.Repeat("-", 40))

		for i, alert := range alertGroup.Alerts {
			fmt.Printf("  [%d] %s - %s (%s)\n",
				i+1, alert.ID, alert.Description, alert.Severity)
		}
	}

	fmt.Printf("\n📈 Total groups: %d\n", len(grouped))
}

// PrettyPrintGroupedBy groups alerts by the specified field and prints them
func (alerts Alerts) PrettyPrintGroupedBy(field string) {
	grouped := alerts.Group(field)
	prettyPrintGrouped(grouped, strings.Title(field))
}