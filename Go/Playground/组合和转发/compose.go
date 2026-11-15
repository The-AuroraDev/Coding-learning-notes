package main

import "fmt"

// 类型必须声明在包级别，方法才能绑定
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
	sol         int
	temperature Temperature
	location    location
}

func main() {
	chengDuTemperature := Temperature{high: 20, low: 10}
	chengDuLocation := location{30.65, 104.06}
	chengDuReport := report{sol: 15, temperature: chengDuTemperature, location: chengDuLocation}

	fmt.Printf("%+v\n", chengDuReport)   // {sol:15 temperature:{high:20 low:10} location:{lat:30.65 long:104.06}}
	fmt.Println(chengDuTemperature.Average())  // 15
}
