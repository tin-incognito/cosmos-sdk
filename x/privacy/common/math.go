package common

import "fmt"

func AddUint64(x, y uint64) (uint64, error) {
	temp := x + y
	if temp < x {
		return 0, fmt.Errorf("Out of range uint64")
	}
	return temp, nil
}

func SubUint64(x, y uint64) (uint64, error) {
	temp := x - y
	if temp > x {
		return 0, fmt.Errorf("Out of range uint64")
	}
	return temp, nil
}
