# 第三方笔记 · 中文学习导航

本目录收录 Go 第三方库（Bubble Tea、Cobra）的教程，并配中文翻译，供中国初学者使用。

## 一、目录内容

| 内容                  | 位置                                                      | 说明                             |
| ------------------- | ------------------------------------------------------- | ------------------------------ |
| Bubble Tea 基础教程（中文） | `bubbletea/basics/README.zh-CN.md`                      | Elm 架构、Model/Init/Update/View  |
| Bubble Tea 命令教程（中文） | `bubbletea/commands/README.zh-CN.md`                    | Cmd 异步 I/O（HTTP 请求）            |
| Bubble Tea 可运行代码    | `bubbletea/basics/main.go`、`bubbletea/commands/main.go` | 对照教程运行                         |
| Cobra 核心章节（中文）      | `cobra/user_guide.zh.md`                                | rootCmd、子命令、flags、RunE         |
| 术语表（中文）             | `术语表.md`                                                | Bubble Tea / Cobra / Go 术语对照   |

> 英文原文保留在同目录（`README.md`、`user_guide.md`），便于对照官方最新版。

## 二、建议学习路径

1. **先打 Go 基础**：变量、struct、方法、error 接口、goroutine 概念、module。
2. **读标准库中文书**的相关章节打底：第 1 章（`fmt`）、第 4 章（时间）、第 6 章（`os`）、第 9 章（测试）、第 13.1 章（`flag`）。
3. **Bubble Tea**：先 `bubbletea/basics` 再 `bubbletea/commands`，边读边跑 `main.go`。
4. **Cobra**：以 `cobra/user_guide.zh.md` 为手册，配合 `cobra-cli` 生成器动手搭一个 CLI。

## 三、环境准备（重要，原文未写）

- 安装 Go（≥ 1.24，见 `bubbletea/go.mod` 要求）。
- 中国网络需设置代理，否则拉不到 GitHub 上的包：
  ```bash
  go env -w GOPROXY=https://goproxy.cn,direct
  ```
- **bubbletea 的 `go.mod` 有个坑**：里面有一行 `replace charm.land/bubbletea/v2 => ../`，指向仓库上层目录。单独拷贝该文件夹会**编译失败**。请改为真实模块路径后执行：
  ```bash
  # 编辑 bubbletea/go.mod，删掉 replace 行，并确保 require 指向真实模块
  go mod tidy
  go run main.go
  ```
- Cobra 安装：`go get github.com/spf13/cobra`。

