# Bubble Tea 基础教程（中文版）

Bubble Tea 基于 [Elm 架构][elm] 的函数式设计范式，而它恰好能和 Go 很好地配合。用它来构建终端应用是一种令人愉快的体验。

本教程假设你已经具备可用的 Go 语言知识。

顺便一提，本程序的“无注释”源码可以在 [GitHub][tut-source] 上找到。

[elm]: https://guide.elm-lang.org/architecture/
[tut-source]: https://github.com/charmbracelet/bubbletea/tree/main/tutorials/basics

## 够了！开始吧

在这个教程里，我们要做一个购物清单。

首先，我们定义包并导入一些库。我们唯一用到的外部包就是 Bubble Tea 库，我们简称它为 `tea`。

```go
package main

import (
    "fmt"
    "os"

    tea "charm.land/bubbletea/v2"
)
```

Bubble Tea 程序由一个描述应用状态的 **model（模型）**，以及该模型上的三个简单方法组成：

- **Init**：一个函数，返回程序启动时运行的初始命令（Cmd）。
- **Update**：一个函数，处理发生的事件并更新模型。
- **View**：一个函数，根据模型中的数据渲染 UI。

## 模型（Model）

先来定义我们的模型，它会保存应用的状态。它可以是任何类型，但通常用 `struct` 最合适。

```go
type model struct {
    choices  []string           // 待办清单上的项目
    cursor   int                // 光标指向的是第几个待办项目
    selected map[int]struct{}   // 哪些待办项目被选中了
}
```

## 初始化（Initialization）

接下来，定义应用的初始状态。这里我们写一个函数返回初始模型；当然，你也可以直接在其他地方把一个初始模型定义成变量。

```go
func initialModel() model {
	return model{
		// 我们的待办清单是一张杂货购物单
		choices:  []string{"Buy carrots", "Buy celery", "Buy kohlrabi"},

		// 一个 map，用来表示哪些项目被选中了。我们把它当成数学中的“集合”来用。
		// 键对应的是上面 choices 切片的索引。
		selected: make(map[int]struct{}),
	}
}
```

之后，我们定义 `Init` 方法。`Init` 可以返回一个执行某些初始 I/O 的 `Cmd`。目前我们不需要做任何 I/O，所以就返回一个 `nil`，意思是“没有命令”。

```go
func (m model) Init() tea.Cmd {
    return nil
}
```

## Update 方法

接下来是 update 方法。当“有事情发生”时就会调用 update 函数。它的职责是查看发生了什么，并据此返回一个更新后的模型。它也可以返回一个 `Cmd` 去触发更多事情，不过现在先不用管这部分。

在我们的例子里，当用户按下“下箭头”时，`Update` 的工作就是察觉到“下箭头被按下了”，并相应地移动光标（或者不作移动）。

“发生的事情”以 `Msg`（消息）的形式传来，它可以是任何类型。消息是某些已发生的 I/O 的结果，比如按键、定时器触发，或服务器的响应。

我们通常用类型 switch 来判断收到了哪种 `Msg`，当然你也可以用类型断言。

目前，我们只处理 `tea.KeyPressMsg` 消息，它会在按键时自动发送给 update 函数。

```go
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {

    // 是按键吗？
    case tea.KeyPressMsg:

        // 好，到底按了哪个键？
        switch msg.String() {

        // 这些键应当退出程序
        case "ctrl+c", "q":
            return m, tea.Quit

        // "up" 和 "k" 把光标向上移动
        case "up", "k":
            if m.cursor > 0 {
                m.cursor--
            }

        // "down" 和 "j" 把光标向下移动
        case "down", "j":
            if m.cursor < len(m.choices)-1 {
                m.cursor++
            }

        // enter 键和空格键切换光标所指项目的选中状态
        case "enter", "space":
            _, ok := m.selected[m.cursor]
            if ok {
                delete(m.selected, m.cursor)
            } else {
                m.selected[m.cursor] = struct{}{}
            }
        }
    }

    // 把更新后的模型交还给 Bubble Tea 运行时去处理。
    // 注意：这里我们没有返回任何命令。
    return m, nil
}
```

你可能注意到了，上面的 <kbd>ctrl+c</kbd> 和 <kbd>q</kbd> 在返回模型的同时返回了一个 `tea.Quit` 命令。那是一个特殊命令，用来指示 Bubble Tea 运行时退出、结束程序。

## View 方法

终于，到了渲染 UI 的时候。在几个方法里，view 是最简单的。我们查看模型当前的状态，据此返回一个 `tea.View`。view 声明了应用的 UI 内容，以及可选的终端特性，比如全窗口模式（即 altscreen 模式）、鼠标跟踪、光标位置等等。

因为 view 描述的就是你应用的整个 UI，所以你完全不用操心“重绘逻辑”之类的事情，Bubble Tea 会帮你搞定。

```go
func (m model) View() tea.View {
    // 头部
    s := "What should we buy at the market?\n\n"

    // 遍历我们的选项
    for i, choice := range m.choices {

        // 光标是否指向这个选项？
        cursor := " " // 没有光标
        if m.cursor == i {
            cursor = ">" // 光标！
        }

        // 这个选项被选中了吗？
        checked := " " // 未选中
        if _, ok := m.selected[i]; ok {
            checked = "x" // 选中！
        }

        // 渲染这一行
        s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
    }

    // 底部
    s += "\nPress q to quit.\n"

    // 发送给运行时进行渲染
    v := tea.NewView(s)
    v.WindowTitle = "Grocery List"

    return v
}
```

> 说明：上面的 `v.WindowTitle` 设置了终端窗口的标题（可选）。原始英文教程的示例代码里没有这一行，这里已和同目录的 `main.go` 对齐保持一致。

## 合在一起

最后一步就是运行我们的程序。把初始模型传给 `tea.NewProgram`，然后让它跑起来：

```go
func main() {
    p := tea.NewProgram(initialModel())
    if _, err := p.Run(); err != nil {
        fmt.Printf("Alas, there's been an error: %v", err)
        os.Exit(1)
    }
}
```

## 下一步

本教程覆盖了构建交互式终端 UI 的基础，但在真实世界里你通常还需要做 I/O。要了解这部分，请看 [命令（Command）教程][cmd]，它相当简单。

此外还有若干 [Bubble Tea 示例][examples] 可供参考，当然也少不了 [Go 文档][docs]。

[cmd]: https://github.com/charmbracelet/bubbletea/tree/main/tutorials/commands/
[examples]: https://github.com/charmbracelet/bubbletea/tree/main/examples
[docs]: https://pkg.go.dev/charm.land/bubbletea/v2?tab=doc

## 额外资源

- [我们和 Bubble Tea 一起使用的库](https://github.com/charmbracelet/bubbletea/#libraries-we-use-with-bubble-tea)
- [野生中的 Bubble Tea](https://github.com/charmbracelet/bubbletea/#bubble-tea-in-the-wild)

### 反馈

我们很想听听你对这个教程的想法。欢迎随时给我们留言！

- [Twitter](https://twitter.com/charmcli)
- [Fediverse](https://mastodon.social/@charmcli)
- [Discord](https://charm.sh/chat)

---

属于 [Charm](https://charm.sh) 项目的一部分。

<a href="https://charm.sh/"><img alt="Charm 的 logo" src="https://stuff.charm.sh/charm-badge.jpg" width="400"></a>

Charm 热爱开源 • Charm loves open source
