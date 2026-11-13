# Cobra 用户指南（中文核心摘录版）

> 本文件是从英文 `user_guide.md` 中**抽取初学者最常用的核心章节**翻译而成，便于快速上手。
> 完整参考（含 hooks、shell 补全、插件、Active Help 等）请仍然查阅英文 `user_guide.md`。

Cobra 是一个用于创建强大现代 CLI（命令行）应用的 Go 库。本指南覆盖从零搭建一个 Cobra 应用的核心用法。

## 一、推荐的项目结构

虽然你可以自由组织代码，但通常一个基于 Cobra 的应用会采用如下结构：

```console
▾ appName/
  ▾ cmd/
      add.go
      your.go
      commands.go
      here.go
  main.go
```

在 Cobra 应用里，`main.go` 通常非常“瘦”，只做一件事：初始化 Cobra。

```go
package main

import "{pathToYourApp}/cmd"

func main() {
  cmd.Execute()
}
```

## 二、使用 Cobra 生成器（最简方式）

`cobra-cli` 是一个独立程序，可以帮你生成应用骨架并添加任意命令，是引入 Cobra 最省事的方式。
完整用法见 [Cobra-CLI 生成器 README](https://github.com/spf13/cobra-cli/blob/main/README.md)。

## 三、使用 Cobra 库（手写方式）

要手动集成 Cobra，你需要一个瘦 `main.go` 和一个根命令（rootCmd）文件，再按需添加子命令。

### 3.1 创建根命令 rootCmd

Cobra 不需要任何特殊构造函数，直接创建命令即可。推荐放在 `app/cmd/root.go`：

```go
var rootCmd = &cobra.Command{
  Use:   "hugo",
  Short: "Hugo is a very fast static site generator",
  Long: `A Fast and Flexible Static Site Generator built with
                love by spf13 and friends in Go.
                Complete documentation is available at https://gohugo.io/documentation/`,
  Run: func(cmd *cobra.Command, args []string) {
    // 在这里做你的事
  },
}

func Execute() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
  }
}
```

你还会在 `init()` 函数里定义标志（flag）并处理配置。下面是一个结合 [viper](https://github.com/spf13/viper)（配置库）的更完整示例：

```go
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// 供 flag 使用的变量
	cfgFile     string
	userLicense string

	rootCmd = &cobra.Command{
		Use:   "cobra-cli",
		Short: "A generator for Cobra based Applications",
		Long: `Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	}
)

// Execute 执行根命令
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
	rootCmd.PersistentFlags().StringP("author", "a", "YOUR NAME", "author name for copyright attribution")
	rootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "name of license for the project")
	rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")
	viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
	viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
	viper.SetDefault("author", "NAME HERE <EMAIL ADDRESS>")
	viper.SetDefault("license", "apache")

	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(initCmd)
}

func initConfig() {
	if cfgFile != "" {
		// 使用 flag 指定的配置文件
		viper.SetConfigFile(cfgFile)
	} else {
		// 查找 home 目录
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// 在 home 目录下搜索名为 ".cobra"（无扩展名）的配置
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".cobra")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
```

### 3.2 创建 main.go

根命令准备好后，让 `main` 函数去执行它。为了清晰，通常对根命令调用 `Execute`（其实任何命令上都能调用）。

```go
package main

import "{pathToYourApp}/cmd"

func main() {
  cmd.Execute()
}
```

### 3.3 创建额外的命令

额外的命令通常各自放在 `cmd/` 目录下的独立文件里。例如要建一个 `version` 命令，就创建 `cmd/version.go`：

```go
package cmd

import (
  "fmt"

  "github.com/spf13/cobra"
)

func init() {
  rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
  Use:   "version",
  Short: "Print the version number of Hugo",
  Long:  `All software has versions. This is Hugo's`,
  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println("Hugo Static Site Generator v0.9 -- HEAD")
  },
}
```

### 3.4 组织子命令

一个命令可以有子命令，子命令还可以再有子命令，通过 `AddCommand` 实现。在较大的应用里，每个子命令可以放在各自的 Go 包中。

推荐做法：父命令用 `AddCommand` 添加它“最直接”的子命令。例如如下目录结构：

```console
├── cmd
│   ├── root.go
│   └── sub1
│       ├── sub1.go
│       └── sub2
│           ├── leafA.go
│           ├── leafB.go
│           └── sub2.go
└── main.go
```

规则是：

* `root.go` 的 `init` 把 `sub1.go` 里定义的命令加到根命令。
* `sub1.go` 的 `init` 把 `sub2.go` 里定义的命令加到 sub1 命令。
* `sub2.go` 的 `init` 把 `leafA.go`、`leafB.go` 里定义的命令加到 sub2 命令。

这种方式能保证子命令在编译期始终被包含，同时避免循环引用。

## 四、返回与处理错误

如果你想把一个错误返回给命令的调用者，可以使用 `RunE`：

```go
package cmd

import (
  "fmt"

  "github.com/spf13/cobra"
)

func init() {
  rootCmd.AddCommand(tryCmd)
}

var tryCmd = &cobra.Command{
  Use:   "try",
  Short: "Try and possibly fail at something",
  RunE: func(cmd *cobra.Command, args []string) error {
    if err := someFunc(); err != nil {
	return err
    }
    return nil
  },
}
```

返回的错误可以在调用 `Execute` 的地方被捕获到。

## 五、使用 Flags（标志）

Flags 是用来控制命令行为的修饰项。

### 5.1 定义接收 flag 值的变量

由于 flag 的定义和使用处在不同位置，我们需要在外面定义一个作用域合适的变量：

```go
var Verbose bool
var Source string
```

### 5.2 持久标志（Persistent Flags）

“持久”意味着该 flag 对当前命令及其**所有子命令**都可用。全局 flag 应作为持久 flag 挂在根命令上：

```go
rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
```

### 5.3 本地标志（Local Flags）

本地 flag 只作用于该特定命令：

```go
localCmd.Flags().StringVarP(&Source, "source", "s", "", "Source directory to read from")
```

### 5.4 父命令上的本地 flag

默认情况下，Cobra 只解析目标命令上的本地 flag，会忽略父命令的本地 flag。通过设置 `Command.TraverseChildren`，Cobra 会在执行目标命令前先解析每个命令上的本地 flag：

```go
command := cobra.Command{
  Use: "print [OPTIONS] [COMMANDS]",
  TraverseChildren: true,
}
```

### 5.5 用 viper 绑定 flag

你可以把 flag 与 [viper](https://github.com/spf13/viper) 绑定：

```go
var author string

func init() {
  rootCmd.PersistentFlags().StringVar(&author, "author", "YOUR NAME", "Author name for copyright attribution")
  viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
}
```

> 注意：当 `--author` flag 由用户提供时，变量 `author` **不会**被配置文件里的值覆盖。

### 5.6 必填 flag（Required）

flag 默认是可选的。若希望某个 flag 未设置时报错，把它标记为必填：

```go
rootCmd.Flags().StringVarP(&Region, "region", "r", "", "AWS region (required)")
rootCmd.MarkFlagRequired("region")
```

持久 flag 用 `MarkPersistentFlagRequired`：

```go
rootCmd.PersistentFlags().StringVarP(&Region, "region", "r", "", "AWS region (required)")
rootCmd.MarkPersistentFlagRequired("region")
```

### 5.7 标志组（Flag Groups）

**必须同时提供**：例如给了 `--username` 就必须给 `--password`：

```go
rootCmd.Flags().StringVarP(&u, "username", "u", "", "Username (required if password is set)")
rootCmd.Flags().StringVarP(&pw, "password", "p", "", "Password (required if username is set)")
rootCmd.MarkFlagsRequiredTogether("username", "password")
```

**互斥**：例如 `--json` 与 `--yaml` 不能同时出现：

```go
rootCmd.Flags().BoolVar(&ofJson, "json", false, "Output in JSON")
rootCmd.Flags().BoolVar(&ofYaml, "yaml", false, "Output in YAML")
rootCmd.MarkFlagsMutuallyExclusive("json", "yaml")
```

**至少提供一个 + 互斥**（实现“二选一”）：

```go
rootCmd.Flags().BoolVar(&ofJson, "json", false, "Output in JSON")
rootCmd.Flags().BoolVar(&ofYaml, "yaml", false, "Output in YAML")
rootCmd.MarkFlagsOneRequired("json", "yaml")
rootCmd.MarkFlagsMutuallyExclusive("json", "yaml")
```

规则要点：本地和持久 flag 都可用于分组；分组仅在该 flag 被定义的所有命令上生效；一个 flag 可出现在多个组中；组可以包含任意数量的 flag。

### 5.8 重复 flag（进阶简述）

- **计数型（Count）**：实现类似 SSH 的 `-v`/`-vv`/`-vvv`，重复出现累加计数，用 `CountVarP`。
- **切片型（Slice）**：收集多个值（如 `--tag a --tag b`），用 `StringSliceVarP` 等。

## 六、位置参数与未知命令建议（简述）

- **位置参数 / 任意参数**：可通过 `Args` 字段校验传入的位置参数（如 `cobra.MinimumNArgs`、`cobra.ExactArgs`、`cobra.ArbitraryArgs` 等）。
- **拼写建议**：当输入未知命令时，Cobra 会基于子命令名做模糊匹配（Levenshtein 距离）给出“你是不是想输入 XXX”的提示；可用 `command.DisableSuggestions = true` 关闭，或用 `SuggestFor` 手动指定。

## 七、与标准库中文内容的连接

- 本指南大量用到**命令行参数解析**思想，可对照标准库中文书 **第 13.1 章 `flag`**（命令行参数解析）一起学习。
- 用到的 `os`、`fmt`、错误处理等，可对照标准库中文书第 1 章（`fmt`）、第 6 章（`os`）。
- `viper` 是第三方配置库，中文书未覆盖；如需要请查阅其 [官方文档](https://github.com/spf13/viper)。

> 提示：本指南里的代码块均为片段，拼装成可运行项目时请配合 `main.go` + `cmd/` 目录结构。环境上需先 `go get github.com/spf13/cobra`（中国网络建议设置 `GOPROXY=https://goproxy.cn,direct`）。
