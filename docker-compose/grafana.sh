#!/bin/sh

# Exported dashboard can't be auto provisioned, fix dashboard before starting Grafana
mkdir -p /var/lib/grafana/dashboards
cat /dashboard.json | sed -e 's/${DS_PROMETHEUS}/1/' > /var/lib/grafana/dashboards/dashboard.json
/run.sh
