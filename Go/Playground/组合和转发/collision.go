package main

import "fmt"

type sol int

type Temperature struct {
	high, low float64
}

func (t Temperature) Average() float64 {
	return (t.high + t.low) / 2
}

type location struct {
	lat, long float64
}

type report struct {
	sol
	Temperature   
	location
}
func (s sol) days() int {
	//....
	return 1
}

func (l location) days() int {
	// ....
	return 2
}

func main() {
	chengDuTemperature := Temperature{high: 20, low: 10}
	chengDuLocation := location{30.65, 104.06}
	chengDuReport := report{sol: 15, Temperature: chengDuTemperature, location: chengDuLocation}

	fmt.Println(chengDuReport.days())  // Error 命名冲突(./collision.go:39:28: ambiguous selector chengDuReport.days)
	fmt.Println(chengDuReport.location.days())   // 没有问题
}
