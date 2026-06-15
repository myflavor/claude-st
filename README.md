# claude-cli

Claude CLI 配置管理工具。

## 安装

```bash
go build -o claude-cli .
```

## 使用方法

### 启动 Claude（推荐）

```bash
./claude-cli start [选项] [claude参数...]
```

**claude-cli 短参数（会先解析）：**
- `-P <provider>`: 注入 provider 环境变量（从 `~/.claude-cli/<provider>/settings.json` 读取）
- `-S`: 跳过权限检查（添加 `-dangerously-skip-permissions`）
- `-H`: 显示帮助

**参数传递规则：**
- `claude-cli` 的短参数（`-P`, `-S`, `-H`）会先被解析
- 其他所有参数原样透传给 claude 官方
- **不修改任何配置文件**，只是临时注入环境变量启动

### 示例

```bash
# 直接启动（不注入环境变量）
./claude-cli start

# 注入 minimax 环境变量并启动
./claude-cli start -P minimax

# 注入环境变量 + 跳过权限检查
./claude-cli start -P minimax -S

# 透传参数给 claude
./claude-cli start -P minimax --model opus
./claude-cli start -S --model sonnet         # 跳过权限 + 指定model
./claude-cli start -P minimax -S -c          # 切换+跳过权限+continue

# 同时启动多个不同 provider（不冲突）
./claude-cli start -P minimax --model opus
./claude-cli start -P deepseek --model sonnet
```

### 传统方式：复制配置文件

```bash
./claude-cli config <provider>
```

将 `~/.claude-cli/<provider>/*` 中的配置文件复制到 `~/.claude/`。

> 注意：这种方式会覆盖原配置，不推荐。推荐使用 `start -P` 临时注入环境变量。

### Provider 配置示例

在 `~/.claude-cli/minimax/settings.json` 中：

```json
{
  "env": {
    "ANTHROPIC_BASE_URL": "https://api.minimaxi.com/anthropic",
    "ANTHROPIC_AUTH_TOKEN": "sk-cp-xxx",
    "ANTHROPIC_MODEL": "MiniMax-M3[1m]",
    "ANTHROPIC_DEFAULT_SONNET_MODEL": "MiniMax-M3[1m]",
    "ANTHROPIC_DEFAULT_OPUS_MODEL": "MiniMax-M3[1m]",
    "ANTHROPIC_DEFAULT_HAIKU_MODEL": "MiniMax-M3[1m]"
  }
}
```

## 目录结构

```
~/.claude-cli/
├── minimax/
│   └── settings.json  # env 字段会被注入到环境变量
├── deepseek/
│   └── settings.json
└── ...
```