# ow-exporter

A Prometheus exporter for overwatch stats.

## Quickstart

Start container

```bash
docker run -p9420:9420 ghcr.io/lexfrei/tools/ow-exporter:master https://playoverwatch.com/en-gb/career/pc/LexFrei-21715/
```

Edit `prometheus.yml`:

```yaml
global:
  scrape_interval:     15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: 'overwatch'
    scrape_interval: 15m
    dns_sd_configs:
      - names:
        - 'tasks.ow-exporter'
        type: 'A'
        port: 9420
    relabel_configs:
      - source_labels: [__address__]
        regex: '.*'
        target_label: instance
        replacement: 'overwatch'
```
