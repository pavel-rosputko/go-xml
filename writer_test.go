package xml

import (
	"bytes"
	"testing"
)

func TestWriter(t *testing.T) {
	const s = "<root><child key='value' /><child1>text</child1></root>"
	b := bytes.NewBufferString("")
	w := NewWriter(b)
	w.StartElement("root").
			Element("child", "key", "value").
			Element("child1", "text").
		EndDocument()

	if !w.stack.isEmpty() { t.Fatal("non empty stack") }

	println("r =", b.String())
	if b.String() != s { t.Fatal("wrong") }
}

func TestJustStart(t *testing.T) {
	const s = "<?xml version='1.0' encoding='utf-8'?><root key='value'>"
	b := bytes.NewBufferString("")
	w := NewWriter(b)
	w.StartDocument().StartElement("root", "key", "value").Send()

	if !w.stack.isEmpty() { t.Fatal("non empty stack") }

	println("r =", b.String())
	if b.String() != s { t.Fatal("wrong") }
}
