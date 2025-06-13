package main

import "time"

func (a Alerts) FilterBySeverity(s string) Alerts {
	var filtered []Alert

	for _, alert := range a.Alerts {
		if alert.Severity == s {
			filtered = append(filtered, alert)
		}
	}

	return Alerts{
		Alerts: filtered,
	}
}

func (a Alerts) FilterByLastMinutes(minutes int) Alerts {
	var filtered []Alert

	// Calculate the cutoff time (current time minus X minutes)
	cutoffTime := time.Now().Add(-time.Duration(minutes) * time.Minute)

	for _, alert := range a.Alerts {
		if alert.Timestamp.After(cutoffTime) {
			filtered = append(filtered, alert)
		}
	}

	return Alerts{
		Alerts: filtered,
	}
}

func (a Alerts) FilterByService(s string) Alerts {
	var filtered []Alert

	for _, alert := range a.Alerts {
		if alert.Service == s {
			filtered = append(filtered, alert)
		}
	}

	return Alerts{
		Alerts: filtered,
	}
}
