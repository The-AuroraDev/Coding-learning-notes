# Bubble Tea 中的命令（Command）（中文版）

这是 Bubble Tea 的第二个教程，讲解用于处理 I/O 的命令（Command）。本教程假设你已经具备可用的 Go 语言知识，并且较好地理解了 [第一个教程][basics]。

本程序“无注释”的版本可以在 [GitHub][source] 找到。

[basics]: https://github.com/charmbracelet/bubbletea/tree/main/tutorials/basics
[source]: https://github.com/charmbracelet/bubbletea/blob/master/tutorials/commands/main.go

## 开始！

在这个教程里，我们要做一个非常简单的程序：它向某个服务器发起一个 HTTP 请求，并报告响应的状态码。

我们先导入几个必要的包，并把要检查的 URL 放进一个 `const` 里。

```go
package main

import (
    "fmt"
    "net/http"
    "os"
    "time"

    tea "charm.land/bubbletea/v2"
)

const url = "https://charm.sh/"
```

## 模型（Model）

接下来定义模型。我们只需要保存两样东西：HTTP 响应的状态码，以及可能的错误。

```go
type model struct {
    status int
    err    error
}
```

## 命令（Cmd）与消息（Msg）

`Cmd` 是一些执行 I/O 然后再返回一个 `Msg` 的函数。查看时间、定时器滴答、从磁盘读取、网络操作——这些都是 I/O，都应该通过命令来运行。这听起来可能有点苛刻，但它能让你的 Bubble Tea 程序保持简洁清晰。

那么，我们来写一个 `Cmd`：它向服务器发请求，并把结果作为 `Msg` 返回。

```go
func checkServer() tea.Msg {

    // 创建一个 HTTP 客户端并发起 GET 请求。
    c := &http.Client{Timeout: 10 * time.Second}
    res, err := c.Get(url)

    if err != nil {
        // 发起请求时出错了。把收到的错误包进一个消息再返回。
        return errMsg{err}
    }
    defer res.Body.Close() // nolint:errcheck

    // 我们收到了服务器的响应。把 HTTP 状态码作为消息返回。
    return statusMsg(res.StatusCode)
}

type statusMsg int

type errMsg struct{ err error }

// 对于包含错误的消息，通常很方便让这个消息也实现 error 接口。
func (e errMsg) Error() string { return e.err.Error() }
```

注意我们定义了两个新的 `Msg` 类型。它们可以是任何类型，甚至可以是空结构体。稍后我们会在 update 函数里用到它们。先来写初始化函数。

## 初始化方法

初始化方法非常简单：我们返回刚才写好的那个 `Cmd`。注意我们并没有调用这个函数；Bubble Tea 运行时会挑合适的时机去调用它。

```go
func (m model) Init() tea.Cmd {
    return checkServer
}
```

## Update 方法

在内部，`Cmd` 会在一个 goroutine 里异步运行。它们返回的 `Msg` 会被收集起来，发送给我们的 update 函数处理。还记得我们写 `checkServer` 命令时定义的那两个消息类型吗？我们就是在这里处理它们。这让处理大量异步操作变得非常轻松。

```go
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {

    case statusMsg:
        // 服务器返回了一个状态消息。把它存进模型。同时告诉 Bubble Tea 运行时我们想要退出，
        // 因为已经没有别的事要做了。我们仍然能够渲染一个带状态码的最终视图。
        m.status = int(msg)
        return m, tea.Quit

    case errMsg:
        // 出错了。记到模型里，并告诉运行时我们做完了、想要退出。
        m.err = msg
        return m, tea.Quit

    case tea.KeyPressMsg:
        // Ctrl+c 退出。即使是运行时间很短的程序，保留一个退出键也是好的，
        // 万一你的逻辑有偏差呢。如果用户退不出去会非常恼火。
		if msg.Mod == tea.ModCtrl && msg.Code == 'c' {
			return m, tea.Quit
		}
    }

    // 如果碰巧收到其他任何消息，什么都不做。
    return m, nil
}
```

## View 函数

我们的 view 非常直白。查看当前模型，据此构建一个 `tea.View`：

```go
func (m model) View() tea.View {
    // 如果有错误，就打印出来，别的什么都不做。
    if m.err != nil {
        return tea.NewView(fmt.Sprintf("\nWe had some trouble: %v\n\n", m.err))
    }

    // 告诉用户我们正在做什么。
    s := fmt.Sprintf("Checking %s ... ", url)

    // 当服务器返回状态码时，把它加到当前这一行上。
    if m.status > 0 {
        s += fmt.Sprintf("%d %s!", m.status, http.StatusText(m.status))
    }

    // 把上面拼出来的内容发送出去渲染。
    return tea.NewView("\n" + s + "\n\n")
}
```

## 运行程序

只剩下运行程序这最后一步，那就开始吧！在这个例子里，我们的初始模型完全不需要任何数据，直接用 `model` 结构体的零值初始化即可。

```go
func main() {
    if _, err := tea.NewProgram(model{}).Run(); err != nil {
        fmt.Printf("Uh oh, there was an error: %v\n", err)
        os.Exit(1)
    }
}
```

这就完成了。不过关于 `Cmd`，还有一件有帮助的事要知道。

## 关于命令的另一件事

在 Bubble Tea 里，`Cmd` 被定义为 `type Cmd func() Msg`。也就是说，它们只是一些“不接受任何参数、返回一个 Msg”的函数，`Msg` 可以是任意类型。如果你需要给命令传参，只要写一个“返回命令的函数”即可。例如：

```go
func cmdWithArg(id int) tea.Cmd {
    return func() tea.Msg {
        return someMsg{id: id}
    }
}
```

一个更贴近真实的例子是这样的：

```go
func checkSomeUrl(url string) tea.Cmd {
    return func() tea.Msg {
        c := &http.Client{Timeout: 10 * time.Second}
        res, err := c.Get(url)
        if err != nil {
            return errMsg{err}
        }
        return statusMsg(res.StatusCode)
    }
}
```

总之，尽量把能做的逻辑都放在最内层的函数里，因为真正异步运行的就是那一个。

## 接下来做什么？

完成这个教程和 [前一个][basics] 之后，你应该已经可以动手写自己的 Bubble Tea 程序了。我们还建议你看看 Bubble Tea 的 [示例程序][examples]，以及 [Bubbles][bubbles]——一个为 Bubble Tea 打造的组件库。

当然，也别忘了查阅 [Go 文档][docs]。

[bubbles]: https://github.com/charmbracelet/bubbles
[docs]: https://pkg.go.dev/charm.land/bubbletea/v2?tab=doc
[examples]: https://github.com/charmbracelet/bubbletea/tree/main/examples

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
