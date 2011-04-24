package xml

import (
	"testing"
)

func TestBuilder(t *testing.T) {
	f := NewBuilder().
		StartElement("a").
			Element("b", "c").
			End()

	println("testBuilder")
	println(f.String())
	f.inspect()
}

func TestBuilderAppend(t *testing.T) {
	f := NewBuilder().
		StartElement("a").
			Element("b", "c").
			End()

	f = NewBuilder().
		StartElement("r").
			Append(f).
			End()

	println("testBuilderAppend")
	println(f.String())
	f.inspect()
}
