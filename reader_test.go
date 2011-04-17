package xml

import (
	"fmt"
	"strings"
	"testing"
)

const readerXml = "<root k='v' kk='vv'><child>text</child></root>"

func TestReadSE(t *testing.T) {
	fmt.Println("testReadSE")
	r := NewReader(strings.NewReader(readerXml))

	f := r.ReadStartElement()
	f.inspect()
}

func TestReadE(t *testing.T) {
	fmt.Println("test-read-e")

	r := NewReader(strings.NewReader(readerXml))
	f := r.ReadElement()
	f.inspect()
}
