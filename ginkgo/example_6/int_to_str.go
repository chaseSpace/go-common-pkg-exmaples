package example_6

import (
	"fmt"
	"strconv"
)

func WithStrConvItoA(i int) string {
	return strconv.Itoa(i)
}

func WithFmtSprintf(i int) string {
	return fmt.Sprintf("%d", i)
}

func WithFormatInt(i int) string {
	return strconv.FormatInt(int64(i), 10)
}
