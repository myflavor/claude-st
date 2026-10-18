# claude-st

Claude Code 的 provider 配置启动工具。用某个 provider 启动 claude（非破坏性，走 `--settings`），不覆盖 `~/.claude/settings.json`。

配置放在 `~/.claude-st/`。

## 安装

### Linux / macOS

```bash
curl -sSL https://raw.githubusercontent.com/myflavor/claude-st/main/scripts/install.sh | bash
```

### Windows (PowerShell)

```powershell
irm https://raw.githubusercontent.com/myflavor/claude-st/main/scripts/install.ps1 | iex
```

安装脚本会：
1. 从 GitHub Releases 下载对应平台的二进制（`claude-st-<os>-<arch>`，Windows 带 `.exe`）
2. 放到和 `claude` 同一个目录（自动在 PATH 里）
3. 若还没有 `~/.claude-st/`，创建一个

### 从源码编译

```bash
go build -o claude-st ./cmd/claude-st
```

## 使用方法

```bash
claude-st <provider> [claude 参数...]
```

`<provider>` 必传，选择 `~/.claude-st/<provider>.json`。

**claude-st 解析的短参数：**
- `-S`：跳过权限检查（添加 `--dangerously-skip-permissions`）
- `-H`：显示帮助

其他所有参数原样透传给 claude。`<provider>` 是第一个非 flag 参数，可以放在 `-S` 之后。

### 工作原理

`claude-st <provider>` 会：
1. 读取 `~/.claude-st/<provider>.json`
2. 若存在 `~/.claude-st/.mcp.json`，通过 `--mcp-config` 注入（所有 provider 共用）
3. 若存在 `~/.claude-st/.env`，把其中 `KEY=VALUE` 导入环境（shell 已设的不覆盖）
4. 若存在 `~/.claude-st/AGENTS.md`，通过 `--append-system-prompt-file` 追加到系统提示（所有 provider 共用）
5. 通过 `claude --settings <provider.json>` 启动

这样可以：
- ✅ 完全覆盖 `~/.claude/settings.json` 的配置
- ✅ 不污染原配置文件
- ✅ 多个 provider 可以共存
- ✅ 多个终端可以同时运行不同 provider
- ✅ 全局 MCP 对每个 provider 生效，且不写项目 `.mcp.json`
- ✅ 全局指令对每个 provider 生效，且不写项目 `AGENTS.md` / `CLAUDE.md`

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

# 额外再叠一层临时 MCP（仍会带上全局 .mcp.json）
claude-st minimax --mcp-config /tmp/extra-mcp.json
```

### 全局 env 文件（~/.claude-st/.env）

若存在 `~/.claude-st/.env`，启动时会逐行导入 `KEY=VALUE` 到环境（所有 provider 共用，且不会写进任何 settings 文件）。shell 里已 export 的同名变量不会被覆盖。

```bash
# ~/.claude-st/.env
CLAUDE_CODE_NO_FLICKER=1                     # 强制 flicker-free 渲染（恢复鼠标点击）
CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC=1   # 禁遥测
# DISABLE_TELEMETRY=1                        # 或只禁遥测
```

> 适合放"每个 provider 都想要、但不想在每个 settings.json 里重复"的通用变量。provider 自己的 `env`（在 `<provider>.json` 里）仍会生效，可覆盖这里设的值。

### 全局指令文件（~/.claude-st/AGENTS.md）

若存在 `~/.claude-st/AGENTS.md`，启动时通过 `--append-system-prompt-file` 追加到 claude 的默认系统提示（所有 provider 共用）。内容不写进项目 `AGENTS.md` / `CLAUDE.md`，也不写进 `~/.claude/CLAUDE.md`，只对 claude-st 启动的会话生效。

```markdown
<!-- ~/.claude-st/AGENTS.md -->
- 回复用中文
- 改完代码先跑 gofmt 和 go test
```

注意事项：
- 它走的是系统提示，不是 memory 体系：不会出现在 `/context` 的 Memory files 里，也不支持 `@path` 导入。
- claude 同一时间只接受一个文件；你自己再传 `--append-system-prompt` 或 `--append-system-prompt-file` 时，claude-st 会跳过这份全局文件，以你的参数为准。

### Provider 配置示例

在 `~/.claude-st/claude.json` 中：

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

### 全局 MCP 配置示例

在 `~/.claude-st/.mcp.json` 中（所有 provider 共用，文件不存在则不注入）：

```json
{
  "mcpServers": {
    "filesystem": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-filesystem", "/tmp"]
    }
  }
}
```

## 目录结构

```
~/.claude-st/
├── .mcp.json         # 可选，全局 MCP（所有 provider 自动注入）
├── .env              # 可选，全局环境变量
├── AGENTS.md         # 可选，全局指令（所有 provider 自动追加到系统提示）
├── claude.json
├── deepseek.json
└── ...
```

## 发布

推送 tag 即可自动构建并发布 Release：

```bash
git tag v0.3.0
git push origin v0.3.0
```

GitHub Action 会自动：
- 在 6 个平台编译 `claude-st`（共 6 个 asset）
- 创建 GitHub Release 并上传所有二进制
- 生成 release notes

### 发布产物

每个 asset 就是对应平台的 `claude-st`（Windows 带 `.exe`）。

| 平台 | 下载文件 |
|------|----------|
| Linux x86_64 | `claude-st-linux-amd64` |
| Linux ARM64 | `claude-st-linux-arm64` |
| macOS Intel | `claude-st-darwin-amd64` |
| macOS Apple Silicon | `claude-st-darwin-arm64` |
| Windows x86_64 | `claude-st-windows-amd64.exe` |
| Windows ARM64 | `claude-st-windows-arm64.exe` |
