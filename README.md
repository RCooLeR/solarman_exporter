<p style="text-align: center">
<img src="./solarman.png" alt="SolarMan 2 Prometheus" width="70%">
</p>

# ☀️ solarman_exporter

A simple **Prometheus exporter** for **SOLARMAN Smart (Solarman OpenAPI)**.  
It polls Solarman Cloud on a fixed interval and exposes metrics for **Prometheus + Grafana**.

---

## ✨ Highlights

- 🔐 Solarman OpenAPI token auth (password is **SHA256-hashed**)
- 🔎 Auto-discovery: stations → station devices (or set device SNs manually)
- ⏱️ Polling loop with caching (does **not** hit the API per scrape)
- 📊 Prometheus-ready metrics for Grafana dashboards
- 🧩 Metric groups: **pv, inverter, grid, battery, bms, gen, load (ups), house**
- 🐳 Docker Compose guide included (Exporter + Prometheus + Grafana)

---

## 📈 Metrics

### ✅ Health

- `solarman_device_up{device_sn}`  
  `1` if last poll for device succeeded, else `0`.

- `solarman_last_update_timestamp_seconds{device_sn}`  
  Unix timestamp of last successful refresh.

### 🧩 Grouped metrics

Each group exposes:

- `solarman_<group>_metric{device_sn,key,name,unit}`

Groups:
- `solarman_pv_metric`
- `solarman_inverter_metric`
- `solarman_grid_metric`
- `solarman_battery_metric`
- `solarman_bms_metric`
- `solarman_gen_metric`
- `solarman_load_metric` — UPS/backup/output
- `solarman_house_metric` — house/home consumption

### 🧺 Generic (optional)

- `solarman_metric{device_sn,key,name,unit}`

Enable/disable with `--enable-generic` / `SOLARMAN_ENABLE_GENERIC`.

> ⚠️ Label-heavy metrics can grow cardinality. Prefer grouped metrics + filters.

---

## 🔑 Credentials & prerequisites

You typically need **two things**:

### 1) Solarman Smart account
- Your Solarman Smart **email** + **password**
- Your plant/device must already be visible in the Solarman Smart app/portal

### 2) Solarman OpenAPI access
You must request OpenAPI credentials:
- `appId`
- `appSecret`

These are not normally visible in the Solarman Smart UI. [How to request them ](https://doc.solarmanpv.com/en/Documentation%20and%20Quick%20Guide#access-process)

### 🧾 Device serial number (SN)
You can find the datalogger/inverter SN:
- on the physical device sticker/label, or
- in Solarman Smart device list (after adding it)

---

## ⚙️ Install

### Requirements
- Go **1.22+**

### Build
```bash
go mod tidy
go build -o solarman_exporter .
```

---

## 🚀 Run

### Using plaintext password

```bash
./solarman_exporter   --app-id "APPID"   --app-secret "APPSECRET"   --email "you@example.com"   --password "yourpassword"   --poll-interval 60s
```

### Using SHA256 password hash (recommended)

```bash
./solarman_exporter   --app-id "APPID"   --app-secret "APPSECRET"   --email "you@example.com"   --password-sha256 "0123abcd..."   --poll-interval 60s
```

### Set device serial numbers manually

```bash
./solarman_exporter   --app-id "APPID"   --app-secret "APPSECRET"   --email "you@example.com"   --password "yourpassword"   --device-sn 1234567890   --device-sn 0987654321
```

If you do **not** specify `--device-sn`, the exporter will try to auto-discover devices from the first station  
(or from `--station-id`).

---

## 🧰 Configuration

All CLI flags can be set via environment variables.

| Flag | Env | Default | Description |
|------|-----|---------|-------------|
| `--listen` | `SOLARMAN_EXPORTER_LISTEN` | `:9876` | HTTP listen address |
| `--metrics-path` | `SOLARMAN_EXPORTER_METRICS_PATH` | `/metrics` | Metrics path |
| `--log-level` | `SOLARMAN_EXPORTER_LOG_LEVEL` | `info` | `debug\|info\|warn\|error` |
| `--base-url` | `SOLARMAN_BASE_URL` | `https://globalapi.solarmanpv.com` | API base URL |
| `--api-version` | `SOLARMAN_API_VERSION` | `v1.0` | API version segment |
| `--language` | `SOLARMAN_LANGUAGE` | `en` | Language param |
| `--app-id` | `SOLARMAN_APP_ID` | (required) | Solarman OpenAPI appId |
| `--app-secret` | `SOLARMAN_APP_SECRET` | (required) | Solarman OpenAPI appSecret |
| `--email` | `SOLARMAN_EMAIL` | (required) | Solarman account email |
| `--password` | `SOLARMAN_PASSWORD` | | Password (plain) |
| `--password-sha256` | `SOLARMAN_PASSWORD_SHA256` | | Password SHA256 hex |
| `--device-sn` | `SOLARMAN_DEVICE_SN` | | Device serials (repeat or comma-separated env) |
| `--station-id` | `SOLARMAN_STATION_ID` | `0` | Station ID for auto-discovery |
| `--poll-interval` | `SOLARMAN_POLL_INTERVAL` | `60s` | Poll interval |
| `--http-timeout` | `SOLARMAN_HTTP_TIMEOUT` | `15s` | HTTP timeout |
| `--enable-generic` | `SOLARMAN_ENABLE_GENERIC` | `true` | Export `solarman_metric` |

---

## 🌍 Base URL (region notes)

Regional API domains may vary. Common options:
- `https://globalapi.solarmanpv.com`
- `https://api.solarmanpv.com`

If you see auth errors, try switching `SOLARMAN_BASE_URL`.

---

## 🐳 Docker Compose (Exporter + Prometheus + Grafana)

This guide runs:
- `solarman_exporter`
- Prometheus
- Grafana

### 📁 Files

Create:

```text
.
├─ docker-compose.yml
└─ prometheus/
   └─ prometheus.yml
```

### 🧩 docker-compose.yml

Create `docker-compose.yml` in the repository root:

```yaml
services:
  solarman_exporter:
    build:
      context: .
    container_name: solarman_exporter
    restart: unless-stopped
    environment:
      # --- Solarman OpenAPI ---
      SOLARMAN_BASE_URL: "https://globalapi.solarmanpv.com"
      SOLARMAN_API_VERSION: "v1.0"
      SOLARMAN_LANGUAGE: "en"

      SOLARMAN_APP_ID: "${SOLARMAN_APP_ID}"
      SOLARMAN_APP_SECRET: "${SOLARMAN_APP_SECRET}"
      SOLARMAN_EMAIL: "${SOLARMAN_EMAIL}"

      # Use ONE of these:
      SOLARMAN_PASSWORD_SHA256: "${SOLARMAN_PASSWORD_SHA256}"
      # SOLARMAN_PASSWORD: "${SOLARMAN_PASSWORD}"

      # Optional (if auto-discovery doesn't work):
      # SOLARMAN_DEVICE_SN: "1234567890,0987654321"
      # SOLARMAN_STATION_ID: "0"

      # Exporter runtime
      SOLARMAN_POLL_INTERVAL: "60s"
      SOLARMAN_HTTP_TIMEOUT: "15s"
      SOLARMAN_ENABLE_GENERIC: "true"
      SOLARMAN_EXPORTER_LOG_LEVEL: "info"
      SOLARMAN_EXPORTER_METRICS_PATH: "/metrics"
      SOLARMAN_EXPORTER_LISTEN: ":9876"

    ports:
      - "9876:9876"

  prometheus:
    image: prom/prometheus:latest
    container_name: prometheus
    restart: unless-stopped
    volumes:
      - ./prometheus/prometheus.yml:/etc/prometheus/prometheus.yml:ro
      - prometheus_data:/prometheus
    ports:
      - "9090:9090"
    depends_on:
      - solarman_exporter

  grafana:
    image: grafana/grafana:latest
    container_name: grafana
    restart: unless-stopped
    environment:
      GF_SECURITY_ADMIN_USER: "admin"
      GF_SECURITY_ADMIN_PASSWORD: "admin"
    volumes:
      - grafana_data:/var/lib/grafana
    ports:
      - "3000:3000"
    depends_on:
      - prometheus

volumes:
  prometheus_data:
  grafana_data:
```

### 🧲 Prometheus scrape config

Create `prometheus/prometheus.yml`:

```yaml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: solarman
    metrics_path: /metrics
    static_configs:
      - targets: ["solarman_exporter:9876"]
```

### 🔒 Environment variables (.env)

Create `.env` in the repository root (**do not commit it**):

```env
SOLARMAN_APP_ID=your_app_id
SOLARMAN_APP_SECRET=your_app_secret
SOLARMAN_EMAIL=you@example.com
SOLARMAN_PASSWORD_SHA256=lowercase_hex_sha256_of_password
# or SOLARMAN_PASSWORD=plaintext_password
```

### ▶️ Run the stack

```bash
docker compose up -d --build
```

### 🔎 Check services

- Exporter metrics: `http://localhost:9876/metrics`
- Prometheus UI: `http://localhost:9090`
- Grafana UI: `http://localhost:3000` (default: `admin` / `admin`)

### 🎛️ Configure Grafana

1. Open Grafana → **Connections / Data sources**
2. Add **Prometheus**
3. URL: `http://prometheus:9090`
4. **Save & test**

### 🧪 Example PromQL queries

Health:

```promql
solarman_device_up
```

PV metrics:

```promql
solarman_pv_metric
```

Totals (energy):

```promql
solarman_totals_metric{unit=~"kwh|wh|mwh"}
```

House/consumption:

```promql
solarman_house_metric
```

---

## 🧠 Notes / Caveats

- Solarman API payloads vary by inverter model; grouping is keyword/unit-based.
- If you want **strict non-overlapping groups**, use “first match wins” logic (put `totals` first).
- Label-heavy metrics can increase Prometheus cardinality; disable generic metrics if needed.
