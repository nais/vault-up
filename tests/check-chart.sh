#!/usr/bin/env bash
set -euo pipefail

temp_dir=$(mktemp -d)
trap 'rm -rf "$temp_dir"' EXIT

for alerts in true false; do
  for canary in true false; do
    helm lint --strict charts/vault-up --set "alerts.enabled=$alerts,vault_up.enabled=$canary"
    helm template vault-up charts/vault-up --namespace nais-system \
      --set "alerts.enabled=$alerts,vault_up.enabled=$canary" > "$temp_dir/rendered.yaml"
    expected=0
    if [[ "$alerts" == true ]]; then expected=$((expected + 1)); fi
    if [[ "$canary" == true ]]; then expected=$((expected + 2)); fi
    actual=$(yq ea '[select(.kind != null)] | length' "$temp_dir/rendered.yaml")
    [[ "$actual" -eq "$expected" ]]
  done
done

helm template vault-up charts/vault-up --namespace nais-system \
  --set vault_up.enabled=true,env=dev-fss > "$temp_dir/rendered.yaml"
yq ea --exit-status 'select(.kind == "Naisjob") |
  .apiVersion == "nais.io/v1" and .spec.schedule == "* * * * *" and
  .spec.activeDeadlineSeconds == 120 and .spec.backoffLimit == 0 and
  .spec.restartPolicy == "Never" and .spec.concurrencyPolicy == "Forbid" and
  .spec.vault.enabled == true and .spec.prometheus == null and
  .spec.replicas == null' "$temp_dir/rendered.yaml" > /dev/null
yq ea '[select(.kind == "PrometheusRule") | .spec.groups[]] | {"groups": .}' "$temp_dir/rendered.yaml" \
  > "$temp_dir/all-rules.yaml"
promtool check rules "$temp_dir/all-rules.yaml"
yq ea 'select(.metadata.name == "vault-up-alerts") | .spec' "$temp_dir/rendered.yaml" \
  > "$temp_dir/rules.yaml"
cp tests/alerts.yaml "$temp_dir/tests.yaml"
promtool test rules "$temp_dir/tests.yaml"

# Match the metadata updates performed by the reusable publishing workflow.
cp -R charts/vault-up "$temp_dir/chart"
export VERSION=20261022000000-test
yq -i '.version = strenv(VERSION) | .appVersion = strenv(VERSION)' "$temp_dir/chart/Chart.yaml"
yq -i '.image.tag = strenv(VERSION)' "$temp_dir/chart/values.yaml"
helm package "$temp_dir/chart" --destination "$temp_dir"
helm template vault-up "$temp_dir/vault-up-$VERSION.tgz" --namespace nais-system \
  --set vault_up.enabled=true > "$temp_dir/packaged.yaml"
yq ea --exit-status 'select(.kind == "Naisjob") |
  .spec.image == "europe-north1-docker.pkg.dev/nais-io/nais/images/vault-up:20261022000000-test"' \
  "$temp_dir/packaged.yaml" > /dev/null
