# claude-st / claude-cfg

Claude Code 的 provider 配置启动与管理工具，拆成两个单职责的二进制：

- **`claude-st`** — 用某个 provider 启动 claude（非破坏性，走 `--settings`）
- **`claude-cfg`** — 把某个 provider 复制覆盖到 `~/.claude/settings.json`（破坏性，不推荐）

两者共用 `~/.claude-cli/` 下的 provider 配置。

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
1. 从 GitHub Releases 下载对应平台的 tarball（`claude-cli-<os>-<arch>.tar.gz`，内含 `claude-st` 和 `claude-cfg`）
2. 解压放到和 `claude` 同一个目录（自动在 PATH 里）

> 解压用的是系统自带的 `tar`（Linux/macOS/Windows 10+ 都内置），无需额外安装任何工具。

### 从源码编译

```bash
go build -o claude-st  ./cmd/claude-st
go build -o claude-cfg ./cmd/claude-cfg
```

## 使用方法

### claude-st —— 用 provider 启动（推荐）

```bash
claude-st <provider> [claude 参数...]
```

`<provider>` 必传，选择 `~/.claude-cli/<provider>.json`。

**claude-st 解析的短参数：**
- `-S`：跳过权限检查（添加 `--dangerously-skip-permissions`）
- `-H`：显示帮助

其他所有参数原样透传给 claude。`<provider>` 是第一个非 flag 参数，可以放在 `-S` 之后。

### 工作原理

`claude-st <provider>` 会：
1. 读取 `~/.claude-cli/<provider>.json`
2. 通过 `claude --settings <provider.json>` 启动

这样可以：
- ✅ 完全覆盖 `~/.claude/settings.json` 的配置
- ✅ 不污染原配置文件
- ✅ 多个 provider 可以共存
- ✅ 多个终端可以同时运行不同 provider

### 示例

```bash
# 用 minimax provider 启动
claude-st minimax

# 跳过权限检查
claude-st minimax -S

# 透传参数给 claude
claude-st minimax --model opus
claude-st minimax -S -c          # 跳过权限 + continue
claude-st -S deepseek            # flag 可以在 provider 之前

# 同时启动多个不同 provider（不冲突）
claude-st claude     # 终端1
claude-st deepseek   # 终端2
```

### claude-cfg —— 复制覆盖配置（不推荐）

```bash
claude-cfg <provider>
```

把 `~/.claude-cli/<provider>.json` 复制到 `~/.claude/settings.json`，会覆盖原配置。

> 一般用 `claude-st` 就够了；`claude-cfg` 只在确实需要写回 `~/.claude/settings.json` 时用。

### Provider 配置示例

在 `~/.claude-cli/claude.json` 中：

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
├── claude.json
├── deepseek.json
└── ...
```

### 从旧版迁移

旧版用文件夹 `~/.claude-cli/<provider>/settings.json`，新版改用单文件 `~/.claude-cli/<provider>.json`。手动迁移：

```bash
mv ~/.claude-cli/<provider>/settings.json ~/.claude-cli/<provider>.json
rmdir ~/.claude-cli/<provider>
```

旧版的 `claude-cli start -P <provider>` 对应新版的 `claude-st <provider>`；`claude-cli config <provider>` 对应 `claude-cfg <provider>`。

## 发布

推送 tag 即可自动构建并发布 Release：

```bash
git tag v0.2.0
git push origin v0.2.0
```

GitHub Action 会自动：
- 在 6 个平台编译 `claude-st` + `claude-cfg` 并打包成一个 tarball（共 6 个 asset）
- 创建 GitHub Release 并上传所有 tarball
- 生成 release notes

### 发布产物

每个 tarball 内含 `claude-st` 和 `claude-cfg`（Windows 带 `.exe`）。

| 平台 | 下载文件 |
|------|----------|
| Linux x86_64 | `claude-cli-linux-amd64.tar.gz` |
| Linux ARM64 | `claude-cli-linux-arm64.tar.gz` |
| macOS Intel | `claude-cli-darwin-amd64.tar.gz` |
| macOS Apple Silicon | `claude-cli-darwin-arm64.tar.gz` |
| Windows x86_64 | `claude-cli-windows-amd64.tar.gz` |
| Windows ARM64 | `claude-cli-windows-arm64.tar.gz` |
