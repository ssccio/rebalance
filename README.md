# kubectl-rebalance

A kubectl plugin to rebalance pods from overloaded nodes by evicting them and allowing Kubernetes to reschedule them on less loaded nodes.

## Installation

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

## Usage

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

## Flags

- `--count, -c`: Number of pods to evict (default: 10)
- `--selector, -l`: Label selector for pods (required)
- `--from-nodes`: Number of most loaded nodes to target (default: 3)
- `--threshold`: CPU/Memory threshold percentage (default: 80.0)
- `--min-pods`: Minimum pods required to proceed (default: 0)
- `--namespace, -n`: Target namespace (empty for all)
- `--interval`: Interval between evictions (default: 10s)- `--dry-run`: Show what would be evicted without doing it
- `--metric`: Metric type: cpu, memory, or both (default: both)
- `--output, -o`: Output format: table, json, yaml (default: table)

## Requirements

- kubectl configured with cluster access
- Metrics Server installed in the cluster
- Appropriate RBAC permissions to:
  - List nodes and pods
  - Get node metrics
  - Create pod evictions

## Example Scenario

After replacing failed nodes in a cluster:

```bash
# Check node utilization
kubectl top nodes

# Rebalance PHP application pods from overloaded nodes
kubectl rebalance --count 50 --selector name=php --from-nodes 5 --dry-run

# If the plan looks good, execute it
kubectl rebalance --count 50 --selector name=php --from-nodes 5
```

## License

MIT