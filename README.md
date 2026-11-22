看了下仓库，README 本身的结构不用动，但有两处新变化值得补进去：根目录多了 `img/`，`.obsidian/` 也已经按白名单追踪了一部分，而且里面确实带了 LaTeX 相关插件（`latex-delimiter-renderer`、`obsidian-latex-suite`），`obsidian-latex-suite` 的 `main.js` 也在版本库里。所以微调集中在**目录结构**和**如何使用**，其余保持不变。

# Coding-learning-notes

> C++ 算法的学习笔记。持续更新，欢迎 PR。

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](CONTRIBUTING.md)

## 这是什么

存放我在学习 **C++ 算法**过程中的笔记与代码片段。以中文为主，注重可运行、可查阅，不追求教材式的完备。笔记用 Obsidian 编辑，仓库里保留了一部分 vault 配置，方便直接在 Obsidian 里打开阅读。

## 如何使用

- 按目录找到主题，直接读代码和要点
- clone 下来修改 `Playground/` 里的示例自行练习
- 用 Obsidian 打开仓库根目录作为 vault：首次打开会处于受限模式，按提示信任作者即可加载已追踪的插件配置；LaTeX 公式渲染是 Obsidian 内置能力，`obsidian-latex-suite` 等插件用于增强输入体验
- 发现错误或缺失，按 [CONTRIBUTING.md](CONTRIBUTING.md) 提 PR

## 文件目录

```
cpp-notes/
├── 00-index/              # 总索引、学习路线、待办
├── 01-language/           # 语言核心
│   ├── basics/            # 类型、表达式、函数、类
│   ├── memory/            # 栈/堆、RAII、智能指针、移动语义
│   ├── templates/         # 模板、泛型、SFINAE、Concepts
│   ├── concurrency/       # thread、mutex、atomic、内存模型
│   └── modern-cpp/        # C++11/14/17/20/23 特性
├── 02-algorithms/         # 算法学习
│ ├── 00-roadmap.md # 算法学习路线、专题顺序
│ ├── 01-basics/ # 复杂度、递归、枚举、模拟
│ ├── 02-data-structures/ # 数组、链表、栈队列、哈希、树、堆、图、并查集
│ ├── 03-techniques/ # 双指针、二分、前缀和、滑动窗口、排序、搜索、DP、贪心、回溯、位运算、字符串、数学
│ ├── 04-templates/ # 算法模板笔记：并查集、线段树、Dijkstra 等
│ ├── 05-problems/ # 题解笔记，按专题或题号组织
│ └── 06-progress.md # 刷题进度、错题清单、待复习
├── 03-stl/                # 容器、迭代器、算法、string、ranges
├── 04-toolchain/          # 编译器、CMake、gdb、sanitizer、clang-format
├── 05-engineering/        # 项目结构、测试、日志、性能、ABI
├── 06-projects/           # 每个小项目一个目录
├── 07-pitfalls/           # 踩坑、编译错误、UB、调试记录
├── wiki/                  # 提到的内容的一些 wiki
├── img/                   # 笔记提到的一些图片附件
└── code/                  # 所有笔记示例代码，保证可编译
```

## 贡献

欢迎修正错别字、指出 bug、补充示例、改进脚本。流程见 [CONTRIBUTING.md](CONTRIBUTING.md)，不确定就先开 Issue。

## 协议

[MIT](LICENSE) © The-AuroraDev
