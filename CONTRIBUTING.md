# 贡献指南

感谢您对 btcapis 项目的关注！我们欢迎所有形式的贡献。

## 开发环境

- Go 1.21+
- 支持的操作系统：Linux, macOS, Windows

## 代码规范

- 遵循 Go 官方代码规范
- 使用 `gofmt` 格式化代码
- 运行 `golangci-lint` 检查代码质量
- 所有公共 API 必须有文档注释

## 提交规范

提交信息格式：

```
type(scope): description

[optional body]

[optional footer]
```

类型说明：

- `feat`: 新功能
- `fix`: 修复 bug
- `docs`: 文档更新
- `style`: 代码格式调整
- `refactor`: 代码重构
- `test`: 测试相关
- `chore`: 构建过程或辅助工具的变动

## 测试要求

- 新功能必须包含测试
- 修复 bug 必须包含回归测试
- 集成测试需要设置环境变量

## 后端适配器开发

新增后端支持时：

1. 在 `providers/` 下创建新目录
2. 实现 `chain.Backend` 接口
3. 添加能力探测
4. 编写测试和文档

## 问题反馈

- 使用 GitHub Issues 报告问题
- 提供详细的复现步骤
- 包含环境信息和错误日志

## 拉取请求

1. Fork 项目
2. 创建功能分支
3. 提交更改
4. 创建 Pull Request
5. 等待代码审查

感谢您的贡献！
