package main

import (
	"fmt"
	"strings"
	"os"
)

type talker interface {
	talk() string
}

func shout(t talker) {   // 可以把 接口 作为形参
	fmt.Println(strings.ToUpper(t.talk()))
}

type c string
type d int

func (x c) talk() string {
	return string(x)
}

func (x d) talk() string {
	res := ""
	for i := 0; i < int(x); i++ {
		res = res + "hello "
	}
	return res
}

func main() {
	shout(c("Hi,Golang"))
	shout(d(2))
	// HI,GOLANG
	// HELLO HELLO
}



func Println(a ...interface{}) (n int, err error) {
	return fmt.Fprintln(os.Stdout, a...)
}