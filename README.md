# victoria-o11y-lab

A hands-on learning project for the **Victoria Metrics observability stack**, demonstrating
metrics, traces, and logs collection from a Go HTTP service.

## Architecture

```mermaid
flowchart LR
    subgraph app [Application]
        GoApp[Go App]
    end

    subgraph vector [Vector]
        OTLP[OTLP gRPC :4317]
        DockerLogs[Docker Logs]
        HostMetrics[Host Metrics]
    end

    subgraph backends [VictoriaMetrics Stack]
        VM[VictoriaMetrics :8428]
        VL[VictoriaLogs :9428]
        VT[VictoriaTraces :8427]
    end

    Grafana[Grafana :3000]

    GoApp -->|"traces (gRPC)"| OTLP
    DockerLogs -->|"container logs"| VL
    HostMetrics -->|"cpu, mem, disk, net"| VM
    OTLP -->|"OTLP/HTTP"| VT

    Grafana --> VM
    Grafana --> VL
    Grafana --> VT
```

## Quick Start

```bash
# Start app + database
make dc-up

# Start observability stack
make dc-o11y-up

# Start everything
make dc-all-up
```

## Services

| Service | Port | URL |
|---------|------|-----|
| Grafana | 3000 | http://localhost:3000 |
| VictoriaMetrics | 8428 | http://localhost:8428 |
| VictoriaLogs | 9428 | http://localhost:9428 |
| VictoriaTraces | 8427 | http://localhost:8427 |
| Vector (OTLP gRPC) | 4317 | — |
| PostgreSQL | 5432 | — |