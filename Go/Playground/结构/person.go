package main

import "fmt"

func main() {

	type person struct {
		name string
		age int
		height float64
	}

	var anan person

	anan.name = "AnAn"
	anan.age = 13
	anan.height = 1.60

	fmt.Println("Name:", anan.name)
	fmt.Println("Age:", anan.age)
	fmt.Println("Height:", anan.height)
	fmt.Println(anan)  // {AnAn 13 1.6}
}