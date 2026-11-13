package main
import "fmt"

func main() {
	type person struct{
		name string
		age int
	}

	var anan = person{name: "AnAn", age: 13}

	fmt.Println(anan)   // {AnAn 13}
}