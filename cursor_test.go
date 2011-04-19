package xml

import (
	"strings"
	"testing"
)

const cursorXml = "<root k='v' kk='vv'><child>text</child></root>"

var gt *testing.T

func cmp(s1, s2 string) {
	if s1 != s2 { gt.Fatal(s1 + " != " + s2) }
}

func TestCursor(t *testing.T) {
	r := NewReader(strings.NewReader(cursorXml))
	c := r.ReadElement().Cursor()

	gt = t

	cmp(c.Name(), "root")
	cmp(c.MustAttr("k"), "v")
	cmp(c.MustAttr("kk"), "vv")
	c.MustToChild()
	cmp(c.Name(), "child")
	cmp(c.MustChars(), "text")
}

func TestChildrenString(t *testing.T) {
	xml := "<a><b/><c/></a>"
	r := NewReader(strings.NewReader(xml))
	c := r.ReadElement().Cursor()

	if c.ChildrenString() != "<b/><c/>" {
		t.Fatal("")
	}
}

func TestChildrenSlice(t *testing.T) {
	println("TestChildrenSlice")
	xml := "<a><b/><c/></a>"
	r := NewReader(strings.NewReader(xml))
	c := r.ReadElement().Cursor()

	f := c.ChildrenSlice()
	println(f.String())
	if f.String() != "<b/><c/>" {
		t.Fatal("")
	}
}
