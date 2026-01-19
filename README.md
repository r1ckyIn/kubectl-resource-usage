# kubectl-resource-usage

<div align="center">

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go&logoColor=white)](https://go.dev)
[![Kubernetes](https://img.shields.io/badge/Kubernetes-Plugin-326CE5?style=flat-square&logo=kubernetes&logoColor=white)](https://kubernetes.io)
[![License](https://img.shields.io/badge/License-MIT-blue?style=flat-square)](LICENSE)

**A kubectl plugin to display Pod CPU/Memory usage percentage relative to requests and limits**

[English](#english) | [中文](#中文)

</div>

---

## English

### Overview

`kubectl-resource-usage` is a native kubectl plugin that displays Pod-level CPU/Memory resource usage rates. Unlike the native `kubectl top pods` command, this plugin calculates usage percentages relative to requests and limits, supports sorting, filtering, and multiple output formats, helping SREs quickly identify resource issues.

### Why This Tool?

| Problem | Solution |
|---------|----------|
| `kubectl top pods` only shows raw values | Shows usage percentage (Request% and Limit%) |
| Hard to identify pods about to exhaust resources | Sort by usage rate, quickly find "dangerous" pods |
| Hard to find resource waste | Request% shows how much of reserved resources are actually used |
| Need to manually calculate percentages | Automatic calculation, saving time during incidents |

### Features

- **Usage Percentage Calculation**: Calculate CPU/Memory Request% and Limit%
- **Namespace Filtering**: Filter pods by namespace (`-n`)
- **Label Filtering**: Filter pods by label selector (`-l`)
- **Sorting**: Sort by CPU or Memory usage rate (`--sort`)
- **Multiple Output Formats**: Support Table and JSON formats (`-o`)
- **Edge Case Handling**: Display `N/A` when requests/limits are not set

### Prerequisites

- Kubernetes cluster with [Metrics Server](https://github.com/kubernetes-sigs/metrics-server) installed
- `kubectl` configured with cluster access
- Go 1.21+ (for building from source)

### Installation

#### From GitHub Releases

```bash
# Download the latest release (Linux/amd64)
curl -LO https://github.com/r1ckyIn/kubectl-resource-usage/releases/latest/download/kubectl-resource_usage-linux-amd64.tar.gz
tar -xzf kubectl-resource_usage-linux-amd64.tar.gz
sudo mv kubectl-resource_usage /usr/local/bin/

# Verify installation
kubectl plugin list
```

#### Build from Source

```bash
git clone https://github.com/r1ckyIn/kubectl-resource-usage.git
cd kubectl-resource-usage
go build -o kubectl-resource_usage ./cmd/kubectl-resource_usage
sudo mv kubectl-resource_usage /usr/local/bin/
```

### Usage

```bash
# View resource usage for all namespaces
kubectl resource-usage

# Filter by namespace
kubectl resource-usage -n payment

# Filter by label
kubectl resource-usage -l app=api

# Sort by memory usage (descending)
kubectl resource-usage --sort memory

# Sort by CPU usage (ascending)
kubectl resource-usage --sort cpu --asc

# Output as JSON
kubectl resource-usage -o json

# Combined: payment namespace + memory sort + json output
kubectl resource-usage -n payment --sort memory -o json
```

### Output Example

**Table Format (default):**

```
NAMESPACE     POD                  CPU_USAGE   CPU_REQ%   CPU_LIM%   MEM_USAGE   MEM_REQ%   MEM_LIM%   NODE
default       api-server-abc       250m        200%       40%        512Mi       117%       29%        node-1
default       worker-xyz           100m        50%        10%        256Mi       80%        20%        node-2
payment       checkout-abc         150m        N/A        30%        128Mi       N/A        12%        node-1
```

### Command Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--namespace` | `-n` | string | all | Filter by namespace |
| `--selector` | `-l` | string | - | Filter by label selector |
| `--sort` | - | string | - | Sort field: cpu or memory |
| `--asc` | - | bool | false | Sort ascending (default: descending) |
| `--output` | `-o` | string | table | Output format: table or json |

### Project Structure

```
kubectl-resource-usage/
├── cmd/
│   └── kubectl-resource_usage/
│       └── main.go           # Entry point
├── pkg/
│   ├── cmd/
│   │   └── resourceusage.go  # Command implementation
│   ├── collector/
│   │   └── metrics.go        # Metrics API data fetching
│   ├── calculator/
│   │   └── usage.go          # Usage calculation logic
│   └── output/
│       ├── table.go          # Table format output
│       └── json.go           # JSON format output
├── go.mod
├── README.md
└── LICENSE
```

---

## 中文

### 项目概述

`kubectl-resource-usage` 是一个 kubectl 原生插件，用于显示 Pod 级别的 CPU/Memory 资源使用率。相比原生的 `kubectl top pods` 命令，本插件额外计算使用率百分比（相对于 requests 和 limits），支持排序、筛选和多种输出格式，帮助 SRE 快速定位资源问题。

### 为什么需要这个工具？

| 问题 | 解决方案 |
|------|----------|
| `kubectl top pods` 只显示原始值 | 显示使用率百分比（Request% 和 Limit%） |
| 难以识别即将耗尽资源的 Pod | 按使用率排序，快速找到"危险"的 Pod |
| 难以发现资源浪费 | Request% 显示预留资源的实际使用情况 |
| 需要手动计算百分比 | 自动计算，告警时节省时间 |

### 功能特性

- **使用率计算**：计算 CPU/Memory 的 Request% 和 Limit%
- **Namespace 筛选**：按命名空间筛选 Pod（`-n`）
- **Label 筛选**：按标签选择器筛选 Pod（`-l`）
- **排序**：按 CPU 或 Memory 使用率排序（`--sort`）
- **多种输出格式**：支持 Table 和 JSON 格式（`-o`）
- **边界处理**：无 requests/limits 时显示 `N/A`

### 前置条件

- Kubernetes 集群已安装 [Metrics Server](https://github.com/kubernetes-sigs/metrics-server)
- `kubectl` 已配置集群访问权限
- Go 1.21+（从源码构建时需要）

### 安装

#### 从 GitHub Releases 下载

```bash
# 下载最新版本（Linux/amd64）
curl -LO https://github.com/r1ckyIn/kubectl-resource-usage/releases/latest/download/kubectl-resource_usage-linux-amd64.tar.gz
tar -xzf kubectl-resource_usage-linux-amd64.tar.gz
sudo mv kubectl-resource_usage /usr/local/bin/

# 验证安装
kubectl plugin list
```

#### 从源码构建

```bash
git clone https://github.com/r1ckyIn/kubectl-resource-usage.git
cd kubectl-resource-usage
go build -o kubectl-resource_usage ./cmd/kubectl-resource_usage
sudo mv kubectl-resource_usage /usr/local/bin/
```

### 使用方法

```bash
# 查看所有 namespace 的资源使用率
kubectl resource-usage

# 按 namespace 筛选
kubectl resource-usage -n payment

# 按 label 筛选
kubectl resource-usage -l app=api

# 按内存使用率降序排序
kubectl resource-usage --sort memory

# 按 CPU 使用率升序排序
kubectl resource-usage --sort cpu --asc

# 输出 JSON 格式
kubectl resource-usage -o json
```

### 技术栈

- **语言**：Go 1.21+
- **主要库**：Cobra, cli-runtime, client-go, metrics
- **构建工具**：GoReleaser
- **CI/CD**：GitHub Actions

---

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## Author

**Ricky** - CS Student @ University of Sydney

[![GitHub](https://img.shields.io/badge/GitHub-r1ckyIn-181717?style=flat-square&logo=github)](https://github.com/r1ckyIn)

Interested in Cloud Engineering & DevOps
