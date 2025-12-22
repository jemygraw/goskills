# GoSkills - Claude Skills 管理工具

[English](README.md) | 简体中文

一个强大的命令行工具，用于解析、管理和执行 Claude Skill 包。GoSkills 是根据 [官方 Claude 文档](https://docs.claude.com/en/docs/agents-and-tools/agent-skills/) 中的规范设计的。

[![License](https://img.shields.io/:license-MIT-blue.svg)](https://opensource.org/licenses/MIT) [![github actions](https://github.com/smallnest/goskills/actions/workflows/go.yml/badge.svg)](https://github.com/smallnest/goskills/actions) [![Go Report Card](https://goreportcard.com/badge/github.com/smallnest/goskills)](https://goreportcard.com/report/github.com/smallnest/goskills)

[![YouTube Video](https://img.youtube.com/vi/Lod9DnAfd9c/maxresdefault.jpg)](https://www.youtube.com/watch?v=Lod9DnAfd9c)

## 特性

- **技能管理**: 从本地目录列出、搜索、解析和检查 Claude 技能
- **运行时执行**: 通过 LLM 集成执行技能（OpenAI、Claude 和兼容 API）
- **Web 界面**: 交互式聊天 UI，支持实时更新、会话回放和丰富的工件渲染（PPT、播客）
- **富内容生成**: 生成 PowerPoint 演示文稿（通过 Slidev）和播客音频
- **深度研究**: 用于深度调查的递归分析和自我纠正能力
- **内置工具**: Shell 命令、Python 执行、文件操作、web 获取和搜索
- **MCP 支持**: 模型上下文协议 (MCP) 服务器集成
- **国际化**: 全面支持英文和中文
- **全面测试**: 完整的测试套件和覆盖率报告

## 安装

### 从源码编译

```shell
git clone https://github.com/smallnest/goskills.git
cd goskills
make
```

### 使用 Homebrew

```shell
brew install smallnest/goskills/goskills
```

或者：

```shell
# 添加 tap
brew tap smallnest/goskills

# 安装 goskills
brew install goskills
```

## 快速开始

```shell
# 设置你的 OpenAI API 密钥
export OPENAI_API_KEY="YOUR_OPENAI_API_KEY"

# 启动 Web 界面
./agent-web

# 列出可用技能
./goskills-cli list ./skills

# 运行技能
./goskills run "创建一个待办事项应用的 React 组件"
```

## 内置工具

GoSkills 包含一套全面的内置工具，用于技能执行：

- **Shell 工具**：执行 shell 命令和脚本
- **Python 工具**：运行 Python 代码和脚本
- **文件工具**：读取、写入和管理文件
- **Web 工具**：获取和处理 Web 内容
- **搜索工具**：Wikipedia 和 Tavily 搜索集成
- **MCP 工具**：与模型上下文协议服务器集成

## CLI 工具

GoSkills 提供用于不同目的的工具套件：

### 1. Web 界面 (`agent-web`)

用于与 GoSkills Agent 交互的现代 Web 界面。
- **聊天**: 与 Agent 进行实时对话。
- **工件**: 直接在浏览器中查看生成的报告、PPT 和播客。
- **历史**: 回放和查看过去的会话。
- **本地化**: 在英文和中文界面之间切换。

### 2. 技能管理 CLI (`goskills-cli`)

位于 `cmd/goskills-cli`，此工具帮助你检查和管理本地 Claude 技能。

#### 构建 `goskills-cli`

```shell
make cli
# 或
go build -o goskills-cli ./cmd/goskills-cli
```

#### 可用命令

- **list**: 列出给定目录中的所有有效技能。
- **parse**: 解析单个技能并显示其结构摘要。
- **detail**: 显示单个技能的完整详细信息，包括完整的正文内容。
- **files**: 列出组成技能包的所有文件。
- **search**: 在目录中按名称或描述搜索技能。

### 3. 技能运行器 CLI (`goskills`)

位于 `cmd/goskills`，此工具通过集成像 OpenAI 模型这样的大型语言模型 (LLM) 来模拟 Claude 技能使用工作流。

#### 构建 `goskills` 运行器

```shell
make runner
# 或
go build -o goskills ./cmd/goskills
```

#### 可用命令

#### download
从 GitHub 目录 URL 下载技能包到 `~/.goskills/skills`。

```shell
# 从 GitHub 下载技能
./goskills download https://github.com/ComposioHQ/awesome-claude-skills/tree/master/meeting-insights-analyzer

# 下载包含子目录的技能
./goskills download https://github.com/ComposioHQ/awesome-claude-skills/tree/master/artifacts-builder
```

download 命令功能：
- 如果 `~/.goskills/skills` 目录不存在，自动创建
- 递归下载所有文件和子目录
- 从 URL 中提取技能名称并将其用作目标目录名
- 通过错误消息防止重复下载

#### run
处理用户请求，首先发现可用技能，然后要求 LLM 选择最合适的技能，最后通过将所选技能的内容作为系统提示提供给 LLM 来执行该技能。

**需要设置 `OPENAI_API_KEY` 环境变量。**

```shell
# 使用默认 OpenAI 模型 (gpt-4o) 的示例
export OPENAI_API_KEY="YOUR_OPENAI_API_KEY"
./goskills run "create an algorithm that generates abstract art"

# 使用自定义 OpenAI 兼容模型和 API 基础 URL（使用环境变量）的示例
export OPENAI_API_KEY="YOUR_OPENAI_API_KEY"
export OPENAI_API_BASE="https://qianfan.baidubce.com/v2"
export OPENAI_MODEL="deepseek-v3"
./goskills run "create an algorithm that generates abstract art"

# 使用自定义 OpenAI 兼容模型和 API 基础 URL（使用命令行标志）的示例
export OPENAI_API_KEY="YOUR_OPENAI_API_KEY"
./goskills run --model deepseek-v3 --api-base https://qianfan.baidubce.com/v2 "create an algorithm that generates abstract art"

# 使用自定义 OpenAI 兼容模型和 API 基础 URL（使用命令行标志），无人工介入自动批准的示例
./goskills run --auto-approve --model deepseek-v3 --api-base https://qianfan.baidubce.com/v2 "使用markitdown 工具解析网页 https://baike.baidu.com/item/%E5%AD%94%E5%AD%90/1584"

# 使用自定义 OpenAI 兼容模型和 API 基础 URL（使用命令行标志），在循环模式下且不自动退出的示例
./goskills run --auto-approve --model deepseek-v3 --api-base https://qianfan.baidubce.com/v2 --skills-dir=./testdata/skills "使用markitdown 工具解析网 页 https://baike.baidu.com/item/%E5%AD%94%E5%AD%90/1584" -l
```


## 开发

### Make 命令

项目包含一个全面的 Makefile 用于开发任务：

```shell
# 帮助 - 显示所有可用命令
make help

# 构建
make build          # 构建 CLI 和运行器
make cli            # 仅构建 CLI
make runner         # 仅构建运行器

# 测试
make test           # 运行所有测试
make test-race      # 运行带竞态检测的测试
make test-coverage  # 运行测试并生成覆盖率报告
make benchmark      # 运行基准测试

# 代码质量
make check          # 运行 fmt-check、vet 和 lint
make fmt            # 格式化所有 Go 文件
make vet            # 运行 go vet
make lint           # 运行 golangci-lint

# 依赖管理
make deps           # 下载依赖
make tidy           # 整理并验证依赖

# 其他
make clean          # 清理构建产物
make install-tools  # 安装开发工具
make info           # 显示项目信息
```

### 运行所有测试

```shell
# 运行全面的测试套件
make test-coverage

# 运行特定工具测试
cd tool && ./test_all.sh
```

## 配置

### 环境变量

- `OPENAI_API_KEY`: LLM 集成的 OpenAI API 密钥
- `OPENAI_API_BASE`: 自定义 API 基础 URL（可选）
- `OPENAI_MODEL`: 自定义模型名称（可选）
- `TAVILY_API_KEY`: Tavily 搜索 API 密钥
- `MCP_CONFIG`: MCP 配置文件路径

### MCP 集成

通过创建 `mcp.json` 文件来配置模型上下文协议 (MCP) 服务器：

```json
{
  "mcpServers": {
    "server-name": {
      "command": "path/to/server",
      "args": ["arg1", "arg2"]
    }
  }
}
```

## 贡献

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交你的更改
4. 运行测试 (`make test`)
5. 运行代码检查 (`make check`)
6. 提交更改 (`git commit -m 'Add amazing feature'`)
7. 推送到分支 (`git push origin feature/amazing-feature`)
8. 打开 Pull Request

## 许可证

本项目在 MIT 许可证下授权 - 详情请参阅 [LICENSE](LICENSE) 文件。