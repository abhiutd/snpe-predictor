package utils

import (
	"fmt"
	"math"
)

// Flop Sizes.
const (
	Flop  = 1
	KFlop = Flop * 1000
	MFlop = KFlop * 1000
	GFlop = MFlop * 1000
	TFlop = GFlop * 1000
	PFlop = TFlop * 1000
	EFlop = PFlop * 1000
)

func logn(n, b float64) float64 {
	return math.Log(n) / math.Log(b)
}

func humanateFlops(s uint64, base float64, sizes []string) string {
	if s < 10 {
		return fmt.Sprintf("%d Flops", s)
	}
	e := math.Floor(logn(float64(s), base))
	suffix := sizes[int(e)]
	val := math.Floor(float64(s)/math.Pow(base, e)*10+0.5) / 10
	f := "%.0f %s"
	if val < 10 {
		f = "%.1f %s"
	}

	return fmt.Sprintf(f, val, suffix)
}

// Flops produces a human readable representation of an SI size.
//
// Flops(82854982) -> 83 MFlops
func Flops(s uint64) string {
	sizes := []string{"Flop", "kFlop", "MFlop", "GFlop", "TFlop", "PFlop", "EFlop"}
	return humanateFlops(s, 1000, sizes)
}
