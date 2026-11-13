package main

import "fmt"

func main() {

	var person struct {
		name string
		age  int
	}

	person.name = "AnAn"
	person.age = 13

	fmt.Println("Name:", person.name)
	fmt.Println("Age:", person.age)
}