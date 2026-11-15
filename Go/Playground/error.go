package main

import (
	"fmt"
	"errors"
)

func main() {
	var ERR = errors.New("err")
	fmt.Println(ERR)
}