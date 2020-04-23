package example_6

import "testing"

func BenchmarkWithStrConvItoAn(b *testing.B) {
	for i := 0; i < b.N; i++ {
		WithStrConvItoA(10000)
	}
}

func BenchmarkWithFmtSprintf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		WithStrConvItoA(10000)
	}
}

func BenchmarkWithFormatInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		WithStrConvItoA(10000)
	}
}
