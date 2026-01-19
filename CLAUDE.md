# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

kubectl-resource-usage is a kubectl native plugin that displays Pod-level CPU/Memory resource usage percentages. Unlike `kubectl top pods` which only shows raw values, this plugin calculates usage percentages relative to requests and limits, helping SREs quickly identify resource issues.

## Build and Run Commands

```bash
# Build the plugin
go build -o kubectl-resource_usage ./cmd/kubectl-resource_usage

# Install locally (move to PATH)
sudo mv kubectl-resource_usage /usr/local/bin/

# Verify kubectl discovers the plugin
kubectl plugin list

# Run the plugin
kubectl resource-usage
kubectl resource-usage -n <namespace>
kubectl resource-usage -l app=api
kubectl resource-usage --sort memory
kubectl resource-usage -o json
kubectl resource-usage -o yaml
kubectl resource-usage -o wide
kubectl resource-usage --above 80          # Filter pods with usage >= 80%
kubectl resource-usage --below 50          # Filter pods with usage <= 50%
kubectl resource-usage --no-limits         # Filter pods without limits configured
```

## Test Commands

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests for a specific package
go test -v ./pkg/calculator/...

# Run tests with coverage
go test -cover ./...
```

## Architecture

```
cmd/kubectl-resource_usage/main.go    # Entry point
pkg/
├── cmd/resourceusage.go              # Cobra command, flag handling, orchestrates collectors/calculator/output
├── collector/
│   ├── metrics.go                    # Fetches PodMetrics from metrics.k8s.io API
│   └── pods.go                       # Fetches Pod specs (requests/limits) from core API
├── calculator/
│   ├── usage.go                      # Calculates Request% and Limit%, sorting
│   ├── filter.go                     # Filtering by usage percentage (--above, --below, --no-limits)
│   └── *_test.go                     # Unit tests
└── output/
    ├── formatter.go                  # Output interface and factory function
    ├── table.go                      # Default table format
    ├── json.go                       # JSON format
    ├── yaml.go                       # YAML format (reuses JSON structures)
    ├── wide.go                       # Wide table with raw request/limit values
    └── output_test.go                # Unit tests
```

### Key Dependencies

- `k8s.io/cli-runtime`: kubectl plugin standard library for kubeconfig/context handling
- `k8s.io/client-go`: Kubernetes API client
- `k8s.io/metrics`: Metrics API client for PodMetrics
- `github.com/spf13/cobra`: CLI framework

### Data Flow

1. **collector** fetches data from two sources:
   - `metrics.k8s.io/v1beta1 PodMetrics` for actual CPU/Memory usage
   - `core/v1 Pod` spec for requests/limits configuration
2. **calculator** computes percentages: `usage / requests × 100%` and `usage / limits × 100%`
3. **calculator/filter** applies threshold filters (--above, --below, --no-limits)
4. **output** formats results as Table, JSON, YAML, or Wide

### Plugin Naming Convention

kubectl plugins use a specific naming pattern:
- Binary name: `kubectl-resource_usage` (underscore)
- Command invocation: `kubectl resource-usage` (hyphen)

## CI/CD

- `.github/workflows/ci.yml`: Runs tests and build on push/PR
- `.github/workflows/release.yml`: Builds and publishes releases via GoReleaser on tag push
