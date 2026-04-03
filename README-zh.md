# Meta Harness

AI 辅助开发的代码审查和测试框架。

## 概述

Meta Harness 是一个基于 Go 的harness 框架，集成了代码审查工具，用于验证 AI 生成的代码变更。

## 项目结构

```
meta_harness/
├── cmd/run/              # 程序入口
├── internal/
│   ├── adapter/          # 外部工具适配器 (Claude CLI, Codex CLI)
│   ├── controller/       # 主编排逻辑
│   ├── executor/         # 测试执行器
│   ├── gate/            # 质量门禁 (规范验证)
│   ├── review/          # 代码审查集成
│   ├── schema/          # 数据结构定义
│   ├── state/           # 状态机
│   ├── testutil/        # 测试工具
│   └── validate/        # 验证逻辑
├── unit_test/           # 单元测试配置
└── artifacts/           # 输出产物
```

## 快速开始

```bash
go run ./cmd/run/main.go
```

## 配置说明

单元测试阈值配置位于 `unit_test/unit_test.yml`。
