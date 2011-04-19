package xml

import (
	"testing"
)

const fragmentXml = "<root k='v' kk='vv'><child>text</child></root>"

func BenchmarkType(b *testing.B) {
	d := desc(0)
	for i := 0; i < b.N; i++ {
		d.depth()
	}
}
