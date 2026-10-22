# vault-up

Checks that the Vault init container works in on-premises Nais clusters.
Every minute, a Naisjob starts a fresh pod and verifies the injected secret.
It exits successfully if the value is `expected`; otherwise it fails.

## Development

Install [mise](https://mise.jdx.dev/), then run:

```sh
mise install
mise run build
mise run test
mise run check
```

## Alerts

- **No recent success:** no successful run for over five minutes, sustained for one minute.
- **Missing canary:** the CronJob is missing from kube-state-metrics for five minutes.

See [the alert rules](charts/vault-up/templates/alerts.yaml) for details.

## Dashboard

Import [grafana/vault-up.json](grafana/vault-up.json) in Grafana and select the shared Prometheus datasource.
