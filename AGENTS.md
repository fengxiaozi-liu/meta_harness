# Agent 指令

## 项目背景

Meta Harness 是一个代码审查和测试框架，用于验证 AI 生成的代码变更是否符合质量标准。

## 技术栈

- **语言**: Go 1.23.2
- **依赖**: 纯标准库，无外部依赖
- **后端集成**: Claude CLI、Codex CLI

## 核心组件

| 组件 | 路径 | 职责 |
|------|------|------|
| Controller | `internal/controller/` | 主编排器，协调整个工作流程 |
| Executor | `internal/executor/` | 执行 AI 后端命令 |
| SpecGate | `internal/gate/` | 验证任务规范合规性 |
| Reviewer | `internal/review/` | 代码审查与问题发现 |
| Validator | `internal/validate/` | 文件验证与测试运行 |
| StateMachine | `internal/state/` | 工作流状态管理 |
| Adapter | `internal/adapter/` | Claude/Codex CLI 适配器 |
| Schema | `internal/schema/` | 数据结构定义 |

## 工作流程

```
TaskSpec → SpecGate → Executor → Reviewer → Validator → Accepted/Failed
                ↓                      ↓
            (失败则终止)          (高严重度问题 → 返工)
```

### 流程说明

1. **Spec Gate** - 验证任务规范是否符合要求
2. **Execute** - 调用 AI 后端执行代码生成
3. **Review** - 分析生成的代码，标记高严重度问题
4. **Validate** - 运行测试并检查覆盖率阈值

## 关键文件

- `unit_test/unit_test.yml` - 测试覆盖率阈值配置
- `internal/gate/spec_gate.go` - 规范验证逻辑
- `internal/controller/controller.go` - 主编排逻辑

## 代码规范

- 使用中文注释
- 注释简洁，一行概括
- 遵循 Go 规范（gofmt, golint）
- 保持测试覆盖率高于配置的阈值
