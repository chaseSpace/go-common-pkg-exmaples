package example_1

import (
	"fmt"
)

var EmptySliceErr = fmt.Errorf("must be non-empty slice")

func Sum(slice []int) (int, error) {
	if len(slice) == 0 {
		return 0, EmptySliceErr
	}
	var sum int
	for _, i := range slice {
		sum += i
	}
	return sum, nil
}
