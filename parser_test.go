package xml

import (
	"strings"
	"testing"
)

const parserXml = "<root k='v'><child ck='cv' /><child>chars</child></root>"

var parserXmlTypes = []int{startType, startType, startType, charsType, endType, endType}
var parserXmlStrings = []string{"root", "k", "v", "child", "ck", "cv", "child", "child", "chars", "child", "root"}

func TestToken(t *testing.T) {
	p := makeParser(strings.NewReader(parserXml))

	i, j := 0, 0
	for {
		tt, m, mm, f := p.token()

		if !f {
			if i < len(parserXmlTypes) { t.Fatal("unexpected eof") }
			if j < len(parserXmlStrings) { t.Fatal("unexpected eof") }
			break
		}

		// println("i =", i, "j =", j, tt, p.string(m[0]), parserXmlTypes[i])

		if a, b := tt, parserXmlTypes[i]; a != b { t.Fatalf("wrong token type, a = %v, b = %v", a, b) }
		if tt == startType {
			if p.string(m) != parserXmlStrings[j] { t.Fatal("wrong start tag name") }
			j++
			for i := 0; i + 1 < len(mm); i, j = i + 2, j + 2 {
				if p.string(mm[i]) != parserXmlStrings[j] { t.Fatal("wrong key") }
				if p.string(mm[i + 1]) != parserXmlStrings[j + 1] { t.Fatal("wrong value") }
			}

			if len(mm) % 2 != 0 {
				if a, b := p.string(m), parserXmlStrings[j]; a != b {
					t.Fatalf("wrong end name, m = %v, j = %v, a = %v, b = %v", m, j, a, b)
				}
				j++
			}
		} else {
			if p.string(m) != parserXmlStrings[j] { t.Fatal("wrong mark string") }
			j++
		}

		i++
	}
}
