package main
import "fmt"

var t interface {
	returnInt() int
}

type a int
type b string

func (a a) returnInt() int {
	return int(a) * 10
}

func (b b) returnInt() int {
	return len(string(b))
}

func main() {
	t = a(100)
	fmt.Println(t.returnInt())  // 1000

	t = b("Hello")
	fmt.Println(t.returnInt())  // 5
}