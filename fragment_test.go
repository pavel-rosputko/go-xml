package xml

import (
	"fmt"
	"strings"
	"testing"
)

const fragmentXml = "<root k='v' kk='vv'><child>text</child></root>"

func BenchmarkType(b *testing.B) {
	d := desc(0)
	for i := 0; i < b.N; i++ {
		d.depth()
	}
}

func TestSetValue(t *testing.T) {
	fmt.Println("test-read-e")

	r := NewReader(strings.NewReader(fragmentXml))
	f := r.ReadElement()
	f.inspect()

	fmt.Println(f.AtString(0))
	fmt.Println(f.String())
	f.SetString(0, "new-name")
	fmt.Println(f.AtString(0))

	fmt.Println(f.String())
}
