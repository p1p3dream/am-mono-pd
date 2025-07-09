package unique

import "testing"

func BenchmarkAlphaNum(b *testing.B) {
	for i := 0; i < b.N; i++ {
		AlphaNum(32)
	}
}
