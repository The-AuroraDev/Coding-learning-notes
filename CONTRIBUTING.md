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
