package main

import "fmt"

// DMS 格式坐标（度、分、秒）
	type coordinate struct {
		d, m, s float64
		h rune
		// rune is an alias for int32 and is equivalent to int32 in all ways. It is
		// used, by convention, to distinguish character values from integer values.
	} 

func main() {
	// ChengDu: 30°39′34.2″N 104°03′56.52″E
	lat  := coordinate{30, 39, 34.2, 'N'}   // 30°39′34.2″N
	long := coordinate{104, 3, 56.52, 'E'}  // 104°03′56.52″E

	fmt.Println(lat.decimal(), long.decimal())  // 30.659499999999998 104.06569999999999
}
// DMS -> DD
func (c coordinate) decimal() float64{
	sign := 1.0 // float
	switch c.h {
		case 'S', 'W', 's', 'w':
			sign = -1
	}
		return sign * (c.d + c.m / 60 + c.s / 3600)
}