package main
import "fmt"

func main() {
	type person struct {
		name string
		age int
	}

	persons := []person {
		{name: "anan", age: 13},
		{name: "jack", age: 14},
		{name: "tom", age: 10},
	}

	fmt.Println(persons)   // [{anan 13} {jack 14} {tom 10}]
	fmt.Println(persons[1].name)  // jack
}