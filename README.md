## CONTRIBUTING.md（精简版）

# 贡献指南

感谢你对 Coding-learning-notes 的关注。这里主要存放 C++ 算法和 Golang 开发的学习笔记，欢迎通过 PR 修正错误、补充内容。

## 可以贡献什么

- 修正笔误、算法实现中的 bug、过时的语法
- 补充示例、边界情况说明、更清晰的解释
- 改进 `Playground/` 下的脚本
- 完善文档（README、本文件、子目录说明）

不确定是否合适？先开 Issue 聊一下，避免白做。

## 提交流程

1. Fork 仓库并 clone 到本地：

   ```bash
   git clone https://github.com/YOUR_USERNAME/Coding-learning-notes.git
   ```

2. 基于最新 `main` 创建分支，分支名带前缀，例如 `docs/fix-go-example`、`feat/add-cpp-sort-note`。

3. 修改后本地跑通代码示例，自查 diff，不要提交构建产物和编辑器配置。

4. Push 并发起 PR 到 `main`，描述里写清楚**改了什么、为什么改**。一个 PR 只做一件事。

## 提交信息

用 Conventional Commits 格式：

```
docs(cpp): 修正快速排序笔记的边界说明
feat(go): 新增 goroutine 泄漏排查笔记
fix(playground): 补上 -std=c++17
```

## 内容约定

- 正文中文为主，代码标识符和专有名词保留英文
- 代码示例要能实际运行，注明语言版本
- 参考了外部资料请注明出处

## 行为准则

保持友善和建设性，遵守 [Contributor Covenant](https://www.contributor-covenant.org/zh-cn/version/2/1/code_of_conduct/)。

# README.md

# Coding-learning-notes

> C++ 算法与 Golang 开发的学习笔记。持续更新，欢迎 PR。

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](CONTRIBUTING.md)

## 这是什么

存放我在学习 **C++ 算法**和 **Golang 开发**过程中的笔记与代码片段。以中文为主，注重可运行、可查阅，不追求教材式的完备。

## 目录结构

```
.
├── C++/
│   ├── Playground/      # 构建与实验脚本
│   └── 算法笔记/
├── Go/
│   ├── Playground/      # 基础框架与实验代码
│   ├── 程序笔记/         # 核心概念与实战
│   └── 第三方笔记/       # 生态、库与工具
├── CONTRIBUTING.md
├── LICENSE
└── README.md
```

## 如何使用

- 按目录找到主题，直接读代码和要点
- clone 下来修改 `Playground/` 里的示例自行练习
- 发现错误或缺失，按 [CONTRIBUTING.md](CONTRIBUTING.md) 提 PR

## 贡献

欢迎修正错别字、指出 bug、补充示例、改进脚本。流程见 [CONTRIBUTING.md](CONTRIBUTING.md)，不确定就先开 Issue。

## 协议

[MIT](LICENSE) © The-AuroraDev