看了下仓库，README 本身的结构不用动，但有两处新变化值得补进去：根目录多了 `img/`，`.obsidian/` 也已经按白名单追踪了一部分，而且里面确实带了 LaTeX 相关插件（`latex-delimiter-renderer`、`obsidian-latex-suite`），`obsidian-latex-suite` 的 `main.js` 也在版本库里。所以微调集中在**目录结构**和**如何使用**，其余保持不变。

# Coding-learning-notes

> C++ 算法与 Golang 开发的学习笔记。持续更新，欢迎 PR。

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](CONTRIBUTING.md)

## 这是什么

存放我在学习 **C++ 算法**和 **Golang 开发**过程中的笔记与代码片段。以中文为主，注重可运行、可查阅，不追求教材式的完备。笔记用 Obsidian 编辑，仓库里保留了一部分 vault 配置，方便直接在 Obsidian 里打开阅读。

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
├── img/                 # 笔记引用的图片
├── .obsidian/           # Obsidian vault 配置（白名单追踪）
├── CONTRIBUTING.md
├── LICENSE
└── README.md
```

## 如何使用

- 按目录找到主题，直接读代码和要点
- clone 下来修改 `Playground/` 里的示例自行练习
- 用 Obsidian 打开仓库根目录作为 vault：首次打开会处于受限模式，按提示信任作者即可加载已追踪的插件配置；LaTeX 公式渲染是 Obsidian 内置能力，`obsidian-latex-suite` 等插件用于增强输入体验
- 发现错误或缺失，按 [CONTRIBUTING.md](CONTRIBUTING.md) 提 PR

## 贡献

欢迎修正错别字、指出 bug、补充示例、改进脚本。流程见 [CONTRIBUTING.md](CONTRIBUTING.md)，不确定就先开 Issue。

## 协议

[MIT](LICENSE) © The-AuroraDev
