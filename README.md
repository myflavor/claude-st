# claude-cli

Claude CLI 配置管理工具。

## 安装

### Linux / macOS

```bash
curl -sSL https://raw.githubusercontent.com/myflavor/claude-cli/main/scripts/install.sh | bash
```

### Windows (PowerShell)

```powershell
irm https://raw.githubusercontent.com/myflavor/claude-cli/main/scripts/install.ps1 | iex
```

安装脚本会：
1. 从 GitHub Releases 下载最新版本的二进制
2. 放到和 `claude` 同一个目录（自动放到 PATH 里）

### 从源码编译

```bash
go build -o claude-cli .
```

## 使用方法

### 启动 Claude（推荐）

```bash
./claude-cli start [选项] [claude参数...]
```

**claude-cli 短参数（会先解析）：**
- `-P <provider>`: 切换到 provider 配置（通过 `--settings` 传给 claude）
- `-S`: 跳过权限检查（添加 `--dangerously-skip-permissions`）
- `-H`: 显示帮助

**参数传递规则：**
- `claude-cli` 的短参数（`-P`, `-S`, `-H`）会先被解析
- 其他所有参数原样透传给 claude 官方

### 工作原理

`-P <provider>` 会：
1. 直接使用 `~/.claude-cli/<provider>/settings.json`
2. 通过 `claude --settings <provider_settings>` 启动

这样可以：
- ✅ 完全覆盖 `~/.claude/settings.json` 的配置
- ✅ 不污染原配置文件
- ✅ 多个 provider 可以共存
- ✅ 多个终端可以同时运行不同 provider

### 示例

```bash
# 直接启动（不切换配置）
./claude-cli start

# 切换到 claude provider 并启动
./claude-cli start -P claude

# 切换 + 跳过权限检查
./claude-cli start -P claude -S

# 透传参数给 claude
./claude-cli start -P claude --model opus
./claude-cli start -S --model sonnet         # 跳过权限 + 指定model
./claude-cli start -P claude -S -c          # 切换+跳过权限+continue

# 同时启动多个不同 provider（不冲突）
./claude-cli start -P claude     # 终端1
./claude-cli start -P deepseek   # 终端2
```

### 传统方式：复制配置文件

```bash
./claude-cli config <provider>
```

将 `~/.claude-cli/<provider>/*` 中的配置文件复制到 `~/.claude/`。

> 注意：这种方式会覆盖原配置，不推荐。推荐使用 `start -P`。

### Provider 配置示例

在 `~/.claude-cli/claude/settings.json` 中：

```json
{
  "env": {
    "ANTHROPIC_BASE_URL": "https://muyuan.do",
    "ANTHROPIC_AUTH_TOKEN": "sk-xxx",
    "ANTHROPIC_MODEL": "claude-opus-4-6[1m]",
    "ANTHROPIC_DEFAULT_SONNET_MODEL": "claude-opus-4-6[1m]",
    "ANTHROPIC_DEFAULT_OPUS_MODEL": "claude-opus-4-6[1m]",
    "ANTHROPIC_DEFAULT_HAIKU_MODEL": "claude-opus-4-6[1m]"
  }
}
```

## 目录结构

```
~/.claude-cli/
├── claude/
│   └── settings.json
├── deepseek/
│   └── settings.json
└── ...
```

## 发布

推送 tag 即可自动构建并发布 Release：

```bash
git tag v0.1.0
git push origin v0.1.0
```

GitHub Action 会自动：
- 在 6 个平台编译（Linux/macOS/Windows × amd64/arm64）
- 创建 GitHub Release 并上传所有二进制
- 生成 release notes

### 支持的二进制

| 平台 | 文件名 |
|------|--------|
| Linux x86_64 | `claude-cli` |
| Linux ARM64 | `claude-cli-linux-arm64` |
| macOS Intel | `claude-cli-darwin-amd64` |
| macOS Apple Silicon | `claude-cli-darwin-arm64` |
| Windows x86_64 | `claude-cli.exe` |
| Windows ARM64 | `claude-cli-arm64.exe` |