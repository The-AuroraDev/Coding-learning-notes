package main

import "fmt"

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
	Temperature   // 转发
	location
}

func main() {
	chengDuTemperature := Temperature{high: 20, low: 10}
	chengDuLocation := location{30.65, 104.06}
	chengDuReport := report{sol: 15, Temperature: chengDuTemperature, location: chengDuLocation}

	fmt.Printf("%+v\n", chengDuReport)   // {sol:15 temperature:{high:20 low:10} location:{lat:30.65 long:104.06}}
	fmt.Println(chengDuReport.high)  // 20
	fmt.Println(chengDuReport.Temperature.high) // 20  与与上一条等价（使用了转发）
	fmt.Println(chengDuTemperature.Average())  // 15
}
