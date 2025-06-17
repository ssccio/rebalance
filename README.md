# kubectl-rebalance

[![CI](https://github.com/ssccio/rebalance/actions/workflows/ci.yml/badge.svg)](https://github.com/ssccio/rebalance/actions/workflows/ci.yml)
[![Release](https://github.com/ssccio/rebalance/actions/workflows/release.yml/badge.svg)](https://github.com/ssccio/rebalance/actions/workflows/release.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/ssccio/rebalance)](https://goreportcard.com/report/github.com/ssccio/rebalance)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A kubectl plugin to redistribute pods from overloaded Kubernetes nodes by intelligently evicting them for rescheduling on less loaded nodes. Perfect for rebalancing clusters after node failures or replacements.

## Table of Contents

<!-- START doctoc -->
<!-- END doctoc -->

## Features

- 🎯 **Targeted eviction** - Select pods by label selector
- 📊 **Smart node selection** - Sort nodes by CPU/memory usage or pod count
- 🔄 **Controlled eviction** - Configurable intervals between evictions
- 🧪 **Safe operation** - Dry-run mode to preview changes
- 📈 **Metrics-optional** - Works without metrics-server (falls back to pod count)
- ⚡ **Real-world ready** - Built for production scenarios

## Installation

### From source

```bash
# Clone the repository
git clone https://github.com/ssccio/rebalance.git
cd rebalance
# Build the plugin
go build -o kubectl-rebalance cmd/kubectl-rebalance/main.go

# Move to PATH
sudo mv kubectl-rebalance /usr/local/bin/

# Verify installation
kubectl rebalance --help
```

### From releases

Download the latest release from the [releases page](https://github.com/ssccio/rebalance/releases):

```bash
# Linux
curl -L https://github.com/ssccio/rebalance/releases/latest/download/kubectl-rebalance_Linux_x86_64.tar.gz | tar xz
sudo mv kubectl-rebalance /usr/local/bin/

# macOS
curl -L https://github.com/ssccio/rebalance/releases/latest/download/kubectl-rebalance_Darwin_x86_64.tar.gz | tar xz
sudo mv kubectl-rebalance /usr/local/bin/
```

### Using Krew (coming soon)

```bash
kubectl krew install rebalance
```

## Usage

### Basic usage

```bash
# Evict 10 pods with label 'name=php' from the 3 most loaded nodes
kubectl rebalance --count 10 --selector name=php --from-nodes 3
# Dry-run mode to see what would be evicted
kubectl rebalance --count 10 --selector name=php --dry-run

# Target nodes with over 90% utilization
kubectl rebalance --count 20 --selector app=web --threshold 90

# Ensure at least 5 pods are found before proceeding
kubectl rebalance --count 10 --selector name=php --min-pods 5

# Focus on memory pressure only
kubectl rebalance --count 10 --selector name=php --metric memory

# Target specific namespace
kubectl rebalance --count 10 --selector name=php --namespace production
```

### For clusters without metrics-server

```bash
# Explicitly use pod count as the metric
kubectl rebalance --count 10 --selector name=php --metric pod-count

# The plugin automatically falls back to pod-count if metrics-server is unavailable
kubectl rebalance --count 10 --selector name=php --dry-run
```

## Flags

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--selector` | `-l` | Label selector for pods **(REQUIRED)** | - |
| `--count` | `-c` | Number of pods to evict | 10 |
| `--from-nodes` | | Number of most loaded nodes to target | 3 |
| `--threshold` | | CPU/Memory threshold percentage | 80.0 |
| `--min-pods` | | Minimum pods required to proceed | 0 |
| `--namespace` | `-n` | Target namespace (empty for all) | "" || `--interval` | | Interval between evictions | 10s |
| `--dry-run` | | Show what would be evicted without doing it | false |
| `--metric` | | Metric type: cpu, memory, both, or pod-count | both |
| `--output` | `-o` | Output format: table, json, yaml | table |

## Use Cases

### Post-node-replacement rebalancing

When nodes rejoin a cluster after maintenance or failure:

```bash
# Check node utilization
kubectl top nodes

# Preview the rebalancing plan
kubectl rebalance --count 50 --selector name=php --from-nodes 5 --dry-run

# Execute the rebalancing
kubectl rebalance --count 50 --selector name=php --from-nodes 5
```

### Regular cluster maintenance

Periodically rebalance to prevent hot spots:

```bash
# Target the top 10% most loaded nodes
kubectl rebalance --count 20 --selector app=web --threshold 90 --from-nodes 10
```

## Requirements

- kubectl configured with cluster access
- Metrics Server installed in the cluster (optional - falls back to pod count)
- Appropriate RBAC permissions:
  - List nodes and pods
  - Get node metrics (if using metrics-server)
  - Create pod evictions

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.