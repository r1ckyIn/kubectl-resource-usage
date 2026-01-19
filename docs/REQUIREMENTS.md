# kubectl-resource-usage 插件 - 需求文档

> 版本: 1.0  
> 日期: 2026-01-19  
> 作者: Ricky

---

## 一、项目概述

### 1.1 项目简介

kubectl-resource-usage 是一个 kubectl 原生插件，用于显示 Kubernetes 集群中 Pod 级别的 CPU/Memory 资源使用率。相比原生的 `kubectl top pods` 命令，本插件额外计算使用率百分比（相对于 requests 和 limits），支持排序、筛选和多种输出格式，帮助 SRE 快速定位资源问题。

### 1.2 项目目标

- 解决 `kubectl top pods` 只显示原始值、无法判断资源是否耗尽的问题
- 提供资源使用率百分比，快速识别「快要爆掉」和「严重浪费」的 Pod
- 支持按使用率排序，半夜告警时快速定位问题
- 无缝集成 kubectl 生态，用户无需学习新工具

### 1.3 技术栈

| 技术 | 用途 |
|------|------|
| Go | 主要编程语言 |
| Cobra | CLI 框架 |
| k8s.io/cli-runtime | kubectl 插件标准库 |
| k8s.io/client-go | Kubernetes API 客户端 |
| k8s.io/metrics | Metrics API 客户端 |
| GoReleaser | 跨平台构建和发布 |
| GitHub Actions | CI/CD |

---

## 二、目标用户与使用场景

### 2.1 目标用户

- DevOps 工程师
- SRE 工程师
- Kubernetes 集群管理员
- 平台工程师

### 2.2 使用场景

| 场景 | 描述 | 用户故事 |
|------|------|----------|
| 半夜告警 | 集群资源告警，需快速定位问题 Pod | 作为 SRE 工程师，我希望在收到「集群内存使用率超过 80%」的告警后，快速找到内存使用率最高的 Pod，以便定位问题根源。 |
| 成本优化 | 排查资源浪费，降低云账单 | 作为 DevOps 工程师，我希望找出「申请了大量资源但实际用得很少」的 Pod，以便调整 requests/limits 降低成本。 |
| 容量规划 | 了解服务资源使用趋势 | 作为平台工程师，我希望定期检查各服务的资源使用率，以便提前扩容或缩容。 |
| 配置审计 | 检查资源配置规范性 | 作为 SRE 工程师，我希望找出没有设置 requests/limits 的 Pod，以便推动团队修复配置问题。 |
| 服务排查 | 定位特定服务的资源问题 | 作为后端开发者，我希望查看我负责的服务（特定 namespace 或 label）的资源使用情况，以便优化应用性能。 |

---

## 三、核心概念

### 3.1 Requests vs Limits

| 概念 | 定义 | 作用 |
|------|------|------|
| **requests** | 调度时 K8s 保证预留的资源量 | 影响 Pod 调度到哪个 Node |
| **limits** | 运行时的资源硬上限 | 超过会被惩罚（CPU 限流，Memory 被杀） |
| **actual usage** | Pod 实际使用的资源量 | 通过 Metrics API 获取 |

### 3.2 两种使用率

| 指标 | 公式 | 回答的问题 | 使用场景 |
|------|------|-----------|----------|
| **Request%** | `usage / requests × 100%` | 「我预留的资源用了多少？」 | 成本优化、发现浪费 |
| **Limit%** | `usage / limits × 100%` | 「我离爆掉还有多远？」 | 稳定性监控、容量预警 |

### 3.3 超过 Limits 的后果

| 资源类型 | 超过 limits 的后果 |
|----------|-------------------|
| **Memory** | Pod 被 OOM Killed（直接杀死） |
| **CPU** | Pod 被 Throttled（限流，变慢但不杀） |

---

## 四、功能需求

### 4.1 核心功能

| 功能 | 描述 | 优先级 |
|------|------|--------|
| 使用率计算 | 计算 CPU/Memory 的 Request% 和 Limit% | P0 |
| Namespace 筛选 | 按命名空间筛选 Pod | P0 |
| Label 筛选 | 按标签选择器筛选 Pod | P0 |
| 排序 | 按 CPU 或 Memory 使用率排序 | P0 |
| 多种输出格式 | 支持 Table 和 JSON 格式输出 | P0 |
| 边界处理 | 无 requests/limits 时显示 N/A | P0 |
| 使用率阈值筛选 | 只显示超过指定阈值的 Pod | P1 |
| 无 limits 筛选 | 只显示未设置 limits 的 Pod | P1 |
| YAML 输出 | 支持 YAML 格式输出 | P2 |
| Wide 输出 | 显示更多列（如 Container 级别） | P2 |

### 4.2 数据来源

| 数据 | API | 说明 |
|------|-----|------|
| CPU/Memory 实际使用量 | `metrics.k8s.io/v1beta1` PodMetrics | 需要集群安装 Metrics Server |
| requests/limits 配置 | `core/v1` Pod spec | 从 Pod 的 containers[].resources 获取 |
| Node 信息 | `core/v1` Pod spec | 从 Pod 的 spec.nodeName 获取 |

### 4.3 错误分类

| 错误类型 | 原因 | 处理方式 |
|----------|------|----------|
| Metrics Server 未安装 | 集群未部署 metrics-server | 报错退出，提示用户安装 |
| kubeconfig 无效 | 配置文件错误或过期 | 报错退出 |
| 无权限访问 namespace | RBAC 限制 | 跳过并警告，继续处理其他 namespace |
| 指定的 namespace 不存在 | 用户输入错误 | 报错退出 |
| 无 Pod 匹配筛选条件 | 筛选条件过严 | 输出空结果，退出码 0 |

---

## 五、命令行接口设计

### 5.1 命令结构

```
kubectl resource-usage [flags]
```

作为 kubectl 插件，本工具只有一个命令，通过 flags 控制行为。

### 5.2 命令详情

#### `kubectl resource-usage`

显示 Pod 级别的资源使用率。

```bash
kubectl resource-usage [flags]
```

**参数：**

| 参数 | 短写 | 类型 | 默认值 | 说明 |
|------|------|------|--------|------|
| `--namespace` | `-n` | string | 所有 | 筛选指定 namespace |
| `--selector` | `-l` | string | - | 按 label 筛选（如 app=api） |
| `--sort` | - | string | - | 排序字段：cpu 或 memory |
| `--asc` | - | bool | false | 升序排列（默认降序） |
| `--output` | `-o` | string | table | 输出格式：table 或 json |

**示例：**

```bash
# 查看所有 namespace 的资源使用率
kubectl resource-usage

# 只看 payment namespace
kubectl resource-usage -n payment

# 只看带有 app=api 标签的 Pod
kubectl resource-usage -l app=api

# 按 memory 使用率降序排列
kubectl resource-usage --sort memory

# 按 cpu 使用率升序排列
kubectl resource-usage --sort cpu --asc

# 输出 JSON 格式
kubectl resource-usage -o json

# 组合使用：payment namespace + memory 降序 + json 输出
kubectl resource-usage -n payment --sort memory -o json

# 复杂 label 筛选
kubectl resource-usage -l "app=api,environment=prod"
```

### 5.3 全局参数

| 参数 | 短写 | 说明 |
|------|------|------|
| `--help` | `-h` | 显示帮助信息 |
| `--kubeconfig` | - | kubeconfig 文件路径（继承自 kubectl） |
| `--context` | - | 使用指定的 context（继承自 kubectl） |

### 5.4 退出码

| 退出码 | 含义 | 触发场景 |
|--------|------|----------|
| 0 | 成功 | 正常执行完成 |
| 1 | 执行失败 | Metrics Server 未安装、连接失败等 |

---

## 六、输出格式

### 6.1 Table 格式（默认）

人类可读的表格格式，用于终端查看。

**正常输出：**

```
NAMESPACE     POD                  CPU_USAGE   CPU_REQ%   CPU_LIM%   MEM_USAGE   MEM_REQ%   MEM_LIM%   NODE
default       api-server-abc       250m        200%       40%        512Mi       117%       29%        node-1
default       worker-xyz           100m        50%        10%        256Mi       80%        20%        node-2
payment       checkout-abc         150m        N/A        30%        128Mi       N/A        12%        node-1
kube-system   coredns-123          50m         25%        N/A        64Mi        32%        N/A        node-2
```

**列说明：**

| 列名 | 说明 |
|------|------|
| NAMESPACE | Pod 所在命名空间 |
| POD | Pod 名称 |
| CPU_USAGE | CPU 实际使用量（如 250m = 0.25 核） |
| CPU_REQ% | CPU 使用量 / CPU requests × 100% |
| CPU_LIM% | CPU 使用量 / CPU limits × 100% |
| MEM_USAGE | Memory 实际使用量（如 512Mi） |
| MEM_REQ% | Memory 使用量 / Memory requests × 100% |
| MEM_LIM% | Memory 使用量 / Memory limits × 100% |
| NODE | Pod 运行所在节点 |

**设计规则：**

| 元素 | 设计 |
|------|------|
| 无 requests/limits | 对应列显示 `N/A` |
| 百分比超过 100% | 正常显示（如 `200%`），表示超额使用 |
| 列对齐 | 左对齐，自动计算宽度 |
| 排序时无值的处理 | N/A 的 Pod 排在最后 |

### 6.2 JSON 格式

机器可读格式，用于脚本处理和 CI/CD 集成。

```json
{
  "items": [
    {
      "namespace": "default",
      "pod": "api-server-abc",
      "node": "node-1",
      "cpu": {
        "usage": "250m",
        "requests": "125m",
        "limits": "500m",
        "requestPercent": 200,
        "limitPercent": 40
      },
      "memory": {
        "usage": "512Mi",
        "requests": "256Mi",
        "limits": "1Gi",
        "requestPercent": 117,
        "limitPercent": 29
      }
    },
    {
      "namespace": "payment",
      "pod": "checkout-abc",
      "node": "node-1",
      "cpu": {
        "usage": "150m",
        "requests": null,
        "limits": "500m",
        "requestPercent": null,
        "limitPercent": 30
      },
      "memory": {
        "usage": "128Mi",
        "requests": null,
        "limits": "1Gi",
        "requestPercent": null,
        "limitPercent": 12
      }
    }
  ]
}
```

**字段说明：**

| 字段 | 类型 | 说明 |
|------|------|------|
| items | array | Pod 列表 |
| namespace | string | Pod 所在命名空间 |
| pod | string | Pod 名称 |
| node | string | 所在节点 |
| cpu.usage | string | CPU 实际使用量 |
| cpu.requests | string/null | CPU requests 配置 |
| cpu.limits | string/null | CPU limits 配置 |
| cpu.requestPercent | int/null | Request 使用率 |
| cpu.limitPercent | int/null | Limit 使用率 |
| memory.* | - | 同 cpu，针对内存 |

---

## 七、技术实现要点

### 7.1 项目结构

```
kubectl-resource-usage/
├── cmd/
│   └── kubectl-resource_usage/
│       └── main.go           # 入口
├── pkg/
│   ├── cmd/
│   │   └── resourceusage.go  # 命令实现
│   ├── collector/
│   │   ├── metrics.go        # Metrics API 数据获取
│   │   ├── pods.go           # Pod 信息获取
│   │   └── collector_test.go
│   ├── calculator/
│   │   ├── usage.go          # 使用率计算逻辑
│   │   └── usage_test.go
│   └── output/
│       ├── formatter.go      # 输出接口定义
│       ├── table.go          # Table 格式实现
│       ├── json.go           # JSON 格式实现
│       └── output_test.go
├── go.mod
├── go.sum
├── .goreleaser.yaml          # 跨平台构建配置
├── .github/
│   └── workflows/
│       ├── ci.yml            # CI 流程
│       └── release.yml       # 发布流程
├── README.md
└── LICENSE
```

### 7.2 插件发现机制

kubectl 插件通过命名约定自动发现：

```bash
# 二进制命名规则
kubectl-resource_usage    # 对应命令 kubectl resource-usage

# 放置位置
$PATH 中任意目录

# 验证插件被发现
kubectl plugin list
```

**命名规则：**
- 命令中的 `-` 在二进制名中变为 `_`
- 例如：`kubectl resource-usage` → `kubectl-resource_usage`

### 7.3 核心数据获取

**获取 Pod Metrics：**

```go
import (
    metricsv1beta1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"
    metricsclient "k8s.io/metrics/pkg/client/clientset/versioned"
)

func GetPodMetrics(client *metricsclient.Clientset, namespace string) (*metricsv1beta1.PodMetricsList, error) {
    return client.MetricsV1beta1().PodMetricses(namespace).List(ctx, metav1.ListOptions{})
}
```

**获取 Pod Spec（requests/limits）：**

```go
import (
    corev1 "k8s.io/api/core/v1"
    "k8s.io/client-go/kubernetes"
)

func GetPods(client *kubernetes.Clientset, namespace string, selector string) (*corev1.PodList, error) {
    return client.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
        LabelSelector: selector,
    })
}
```

### 7.4 使用率计算

```go
type ResourceUsage struct {
    Usage          resource.Quantity
    Requests       *resource.Quantity  // nil 表示未设置
    Limits         *resource.Quantity  // nil 表示未设置
    RequestPercent *int                // nil 表示无法计算
    LimitPercent   *int                // nil 表示无法计算
}

func CalculatePercent(usage, base *resource.Quantity) *int {
    if base == nil || base.IsZero() {
        return nil  // 返回 N/A
    }
    percent := int(usage.MilliValue() * 100 / base.MilliValue())
    return &percent
}
```

### 7.5 排序逻辑

```go
func SortByLimitPercent(pods []PodUsage, field string, ascending bool) {
    sort.Slice(pods, func(i, j int) bool {
        var vi, vj *int
        if field == "cpu" {
            vi, vj = pods[i].CPU.LimitPercent, pods[j].CPU.LimitPercent
        } else {
            vi, vj = pods[i].Memory.LimitPercent, pods[j].Memory.LimitPercent
        }
        
        // N/A 排最后
        if vi == nil && vj == nil { return false }
        if vi == nil { return false }  // i 排后面
        if vj == nil { return true }   // j 排后面
        
        if ascending {
            return *vi < *vj
        }
        return *vi > *vj
    })
}
```

---

## 八、开发计划

| 日期 | 任务 | 产出 |
|------|------|------|
| Day 1 | 项目初始化、cli-runtime 集成、连接 K8s 集群 | 能获取并打印 PodMetrics |
| Day 2 | 获取 Pod requests/limits、实现使用率计算 | 核心计算逻辑完成 |
| Day 3 | 实现筛选（-n, -l）、排序（--sort）、Table/JSON 输出 | 全部 MVP 功能完成 |
| Day 4 | GitHub Actions CI、GoReleaser 配置、README 文档 | 可发布版本 |

---

## 九、验收标准

### 9.1 功能验收

- [ ] 能正确计算 CPU/Memory 的 Request% 和 Limit%
- [ ] `-n` 能正确筛选 namespace
- [ ] `-l` 能正确筛选 label
- [ ] `--sort cpu` 能按 CPU 使用率降序排序
- [ ] `--sort memory` 能按 Memory 使用率降序排序
- [ ] `--asc` 能切换为升序
- [ ] `-o table` 输出格式正确、列对齐
- [ ] `-o json` 输出格式正确、可被 jq 解析
- [ ] 无 requests/limits 时显示 `N/A`
- [ ] Metrics Server 未安装时给出友好提示

### 9.2 工程验收

- [ ] `kubectl plugin list` 能发现插件
- [ ] 支持 `--kubeconfig` 和 `--context` 参数
- [ ] GitHub Actions CI 通过（lint + test）
- [ ] GoReleaser 能生成 Linux/macOS/Windows 二进制
- [ ] README 包含安装说明、使用示例、截图

---

## 十、简历亮点

完成此项目后，可在简历中体现以下技能：

| 技能点 | 体现 |
|--------|------|
| Kubernetes 生态 | 理解 Metrics API、requests/limits 资源模型 |
| kubectl 插件开发 | 使用 cli-runtime 开发原生插件 |
| Go 语言 | client-go、结构体设计、接口抽象 |
| DevOps 思维 | 解决实际运维痛点、多格式输出便于集成 |
| 跨平台发布 | GoReleaser 自动构建多平台二进制 |
| CI/CD 实践 | GitHub Actions 自动化测试和发布 |

---

## 十一、面试准备要点

### 11.1 产品设计类

1. **为什么要做这个工具？**  
   `kubectl top pods` 只显示原始值，无法判断资源是否快要耗尽或存在浪费。

2. **为什么要同时计算 Request% 和 Limit%？**  
   Request% 用于发现浪费（成本优化），Limit% 用于预警风险（稳定性监控）。

3. **为什么默认降序排列？**  
   SRE 90% 的场景是找「最危险」或「最浪费」的 Pod，都需要看最高的。

4. **为什么 flag 用 `-n`、`-l`、`-o`？**  
   遵循 kubectl 惯例，降低用户学习成本。

### 11.2 技术实现类

1. **Kubernetes Metrics API 是怎么工作的？**  
   Metrics Server 从各节点的 kubelet 收集数据，暴露为 `metrics.k8s.io/v1beta1` API。

2. **kubectl 插件的发现机制是什么？**  
   kubectl 在 `$PATH` 中搜索 `kubectl-*` 命名的可执行文件，自动注册为子命令。

3. **如何处理没有设置 limits 的 Pod？**  
   使用 `*int` 类型表示百分比，`nil` 表示无法计算，输出时显示 `N/A`。

4. **排序时 N/A 怎么处理？**  
   N/A 的 Pod 排在最后，因为它们不参与使用率比较。

### 11.3 工程实践类

1. **项目结构是怎么设计的？**  
   按职责分层：cmd（命令入口）、collector（数据获取）、calculator（计算逻辑）、output（输出格式）。

2. **为什么选择 Table 和 JSON 两种输出格式？**  
   Table 给人看，JSON 给程序处理（脚本、jq、CI/CD 集成）。

3. **CI/CD 流水线是怎么配置的？**  
   PR 触发 lint + test，tag 触发 GoReleaser 构建多平台二进制并发布到 GitHub Releases。

---

## 十二、参考资料

- [Kubernetes Metrics Server](https://github.com/kubernetes-sigs/metrics-server)
- [kubectl 插件开发指南](https://kubernetes.io/docs/tasks/extend-kubectl/kubectl-plugins/)
- [cli-runtime 库文档](https://pkg.go.dev/k8s.io/cli-runtime)
- [client-go 示例](https://github.com/kubernetes/client-go/tree/master/examples)
- [sample-cli-plugin 官方示例](https://github.com/kubernetes/sample-cli-plugin)
- [GoReleaser 文档](https://goreleaser.com/)
