package logTool

import (
	"testing"
)

func BenchmarkInfoAJ(b *testing.B) {
	b.N = 300000
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		InfoAJ("111")
	}
}
