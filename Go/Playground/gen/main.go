package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
)

func newFloat(fan int) float64 {
	return rand.Float64() + float64(rand.Intn(fan))
}

func newInt(fan int) int {
	return rand.Intn(fan)
}

func main() {
	if len(os.Args) != 5 {
		log.Fatalln("用法: ./gen <数量> <类型> <范围> <输出位数>")
	}

	sum, err := strconv.Atoi(os.Args[1])
	if err != nil || sum < 0 {
		log.Fatalln("数量必须是非负整数")
	}

	tttype := os.Args[2]
	if tttype != "int" && tttype != "double" {
		log.Fatalln("类型错误：仅支持 int、double")
	}

	fan, err := strconv.Atoi(os.Args[3])
	if err != nil || fan <= 0 {
		log.Fatalln("范围必须是正整数")
	}

	for i := 0; i < sum; i++ {
		if tttype == "int" {
			fmt.Println(newInt(fan))
		} else {
			fmt.Println(newFloat(fan))
		}
	}
}
