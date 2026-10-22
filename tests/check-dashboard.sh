#!/usr/bin/env bash
set -euo pipefail

temp_dir=$(mktemp -d)
trap 'rm -rf "$temp_dir"' EXIT

# Test the actual dashboard queries, replacing Grafana's cluster filter with All.
yq -p=json -o=yaml '.templating.list[0].allValue as $clusters |
{"groups": [{"name": "vault-up-dashboard", "rules": [
  .panels[] | . as $panel | .targets[] | {
    "record": ("vault_up_dashboard_panel_" + ($panel.id | tostring) + "_" + .refId),
    "expr": (.expr | sub("\\$\\{cluster:regex\\}"; $clusters))
  }
]}]}' grafana/vault-up.json > "$temp_dir/dashboard-rules.yaml"
promtool check rules "$temp_dir/dashboard-rules.yaml"
cp tests/dashboard.yaml "$temp_dir/tests.yaml"
promtool test rules "$temp_dir/tests.yaml"
