# enc-alertbuddy
Simple alert grouping and priority planning CLI tool

# How to install
`go install github.com/arho/enc-alertbuddy`

# How to use
`enc-alertbuddy -i <myalerts.json>`

# Help
```
❯ ./enc-alertbuddy -h
enc-alertbuddy - Alert Management CLI Tool

USAGE:
  enc-alertbuddy -i <input-file> [OPTIONS]

REQUIRED FLAGS:
  -i, --input <file>     Input JSON file containing alerts

OPTIONAL FLAGS:
  --groupby <field>      Group alerts by field (severity, service, component, metric, threshold, value, priority)
  --lastminutes <n>      Filter alerts from the last N minutes
  --show-all, -a         Show all alerts in detailed format
  -v, --version          Show version information
  -h, --help             Show this help message

EXAMPLES:
  enc-alertbuddy -i alerts.json
  enc-alertbuddy -i alerts.json --show-all
  enc-alertbuddy -i alerts.json -a
  enc-alertbuddy -i alerts.json --groupby=severity
  enc-alertbuddy -i alerts.json --groupby severity
  enc-alertbuddy -i alerts.json --lastminutes=30
  enc-alertbuddy -i alerts.json --groupby=service --lastminutes=60
  enc-alertbuddy -i alerts.json --show-all --lastminutes=30

VALID GROUPBY FIELDS:
  severity    - Group by alert severity (critical, warning, info)
  service     - Group by service name
  component   - Group by component name
  metric      - Group by metric type
  threshold   - Group by threshold value
  value       - Group by metric value
  priority    - Group by calculated priority score

FEATURES:
  • Automatic priority calculation based on severity, deviation, and affected components
  • Time-based filtering to focus on recent alerts
  • Flexible grouping by any alert field
  • Show all alerts in detailed format with --show-all
  • Beautiful formatted output for better readability
  • Support for large alert datasets
```

# Example
```
❯ ./enc-alertbuddy -i sample-alerts.json
📊 Loaded 148 alerts from sample-alerts.json

🔥 Top 10 Highest Priority Alerts:
============================================================

[1] Priority: 112.00 | critical | ALT-1008
    Service: analytics-platform | Component: data-processor
    Metric: processing_lag (3600.00 / 300.00)
    Description: Data processing lag severely behind schedule

[2] Priority: 112.00 | critical | ALT-1138
    Service: distributed-tracing | Component: span-collector
    Metric: dropped_spans (1250.00 / 100.00)
    Description: Distributed tracing span drop rate critically high

[3] Priority: 97.00 | critical | ALT-1096
    Service: geospatial-index | Component: proximity-calculator
    Metric: calculation_errors (95.00 / 10.00)
    Description: Geospatial proximity calculation error rate critically high

[4] Priority: 92.00 | critical | ALT-1121
    Service: data-encryption | Component: key-rotator
    Metric: rotation_failures (45.00 / 5.00)
    Description: Encryption key rotation failures critically high

[5] Priority: 92.00 | critical | ALT-1079
    Service: distributed-lock | Component: lease-manager
    Metric: lock_timeouts (45.00 / 5.00)
    Description: Distributed lock timeout rate critically high

[6] Priority: 92.00 | critical | ALT-1058
    Service: task-scheduler | Component: cron-manager
    Metric: missed_executions (45.00 / 5.00)
    Description: Critical number of scheduled task executions missed

[7] Priority: 87.00 | critical | ALT-1145
    Service: api-documentation | Component: spec-generator
    Metric: generation_failures (85.00 / 10.00)
    Description: API documentation generation failures critically high

[8] Priority: 87.00 | critical | ALT-1135
    Service: feature-store | Component: feature-server
    Metric: serving_latency (850.00 / 100.00)
    Description: Feature store serving latency critically high

[9] Priority: 77.00 | critical | ALT-1142
    Service: deployment-manager | Component: rollback-controller
    Metric: rollback_failures (15.00 / 2.00)
    Description: Deployment rollback failures critically high

[10] Priority: 76.00 | critical | ALT-1110
    Service: api-gateway-v2 | Component: request-router
    Metric: routing_errors (185.00 / 25.00)
    Description: API gateway routing errors indicate configuration issues

... and 138 more alerts (use --show-all to see all or --groupby to organize)

============================================================
📈 SUMMARY STATISTICS
============================================================

🚨 Severity Breakdown:
  Critical: 43 alerts (29.1%)
  Warning: 63 alerts (42.6%)
  Info: 42 alerts (28.4%)

⚡ Average Priority Score: 24.07

🏢 Top 5 Services by Alert Count:
  email-service: 1 alerts
  web-crawler: 1 alerts
  content-management: 1 alerts
  mobile-push: 1 alerts
  media-transcoder: 1 alerts
```


# Assingment notes

In addition to the required functions, the program was wrapped with a CLI tool, security checks with CodeQL, tests were added, and a CI pipeline is also added.

Boilerplate code is generated with  Claude, most of the formatting, wrapping and CLI code is also generated with that. Function code testing, weird decisions, and the obsession with overcomplexity was corrected manually, and most of the deprecated code was also replaced with current packages, best practices etc.
