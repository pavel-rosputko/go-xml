package xml

import (
	"io"
	"utf8"
	"os"
)

// NOTE separate desc types from ?
const (
	startType = iota
	keyType
	valueType
	endType
	charsType
	cdataType
	directiveType
	commentType
	eofType
)

// NOTE a: Indexer
// NOTE: use bytes.Buffer instead of bytes ?
type parser struct {
	reader	io.Reader
	next	byte
	index	int
	length	int
	eof	bool
	bytes	[]byte
}

func newParser(reader io.Reader) *parser {
	return &parser{reader: reader, index: -1, length: 0}
}

// NOTE use slice ?
type mark struct { s, e int }

// NOTE eofType instead of f bool ?
func (p *parser) token() (tokenType int, marks []mark, f bool) {
	println("token")
	p.get()
	if p.eof { return }

	f = true
	if p.next != '<' { // text
		index := p.index
		for !p.eof && p.next != '<' { p.get() }

		tokenType, marks = charsType, []mark{{index, p.index}}
		p.back()
		return
	}

	p.mustGet()
	switch p.next {
	case '/': // end tag
		p.mustGet()
		index := p.index

		p.name()
		tokenType, marks = endType, []mark{{index, p.index}}

		p.space()

		if p.next != '>' { p.error("invalid characters between </ and >") }
	case '?': // directive
		p.mustGet()
		index := p.index

		var b byte
		for !(b == '?' && p.next == '>') {
			b = p.next
			p.mustGet()
		}

		tokenType, marks = directiveType, []mark{{index, p.index - 1}}
	case '!': // comment or cdata
		p.mustGet()
		switch p.next {
		case '-':
			p.mustGet()
			if p.next != '-' { p.error("invalid sequence <!- not part of <!--") }

			p.mustGet()
			index := p.index

			var b0, b1 byte
			for !(b0 == '-' && b1 == '-' && p.next == '>') {
				b0, b1 = b1, p.next
				p.mustGet()
			}

			tokenType, marks = commentType, []mark{{index, p.index - 2}}
		case '[':
			for i := 0; i < 6; i++ {
				p.mustGet()
				if p.next != "CDATA"[i] { p.error("invalid <![ sequence") }
			}

			p.mustGet()
			index := p.index

			var b byte
			for !(b == ']' && p.next == ']') {
				b = p.next
				p.mustGet()
			}

			tokenType, marks = cdataType, []mark{{index, p.index - 1}}
		default:
			// probably a directive <!
		}
	default: // start tag
		index := p.index
		p.name()
		tokenType, marks = startType, []mark{{index, p.index}}

		var empty bool
		for {
			p.space()

			if p.next == '/' { // empty tag
				empty = true
				p.mustGet()
				if p.next != '>' { p.error("expected /> in element") }
				break
			}

			if p.next == '>' { break }

			index := p.index
			p.name()
			marks = append(marks, mark{index, p.index})

			p.space()

			if p.next != '='  { p.error("attribute name without = in element") }
			p.mustGet()

			p.space()

			if p.next != '\'' && p.next != '"' { p.error("unquoted attribute value") }
			delim := p.next

			p.mustGet()
			index = p.index

			for p.next != delim { p.mustGet() }

			marks = append(marks, mark{index, p.index})

			p.mustGet()
		}

		if empty { marks = append(marks, mark{}) }
	}

	return
}

func (p *parser) clean() {
	p.index = -1
	p.length = 0
	p.bytes = []byte{}
}

func (p *parser) markEq(m1, m2 mark) bool {
	if m1.e - m1.s != m2.e - m2.s { return false }

	i, j := m1.s, m2.s
	for i < m1.e {
		if p.bytes[i] != p.bytes[j] { return false }
		i++; j++
	}

	return true
}

func (p *parser) string(m mark) string {
	return string(p.bytes[m.s:m.e])
}

func isNameByte(b byte) bool {
	return 'A' <= b && b <= 'Z' ||
		'a' <= b && b <= 'z' ||
		'0' <= b && b <= '9' ||
		b == '_' || b == ':' || b == '.' || b == '-'
}

func (p *parser) name() {
	if !(p.next >= utf8.RuneSelf || isNameByte(p.next)) { p.error("invalid tag name first letter") }

	p.mustGet()

	for p.next >= utf8.RuneSelf || isNameByte(p.next) { p.mustGet() }

	// TODO check the characters in [i:p.index]
}

func (p *parser) space() {
	for p.next == ' ' || p.next == '\r' || p.next == '\n' || p.next == '\t' { p.mustGet() }
}


func (p *parser) back() {
	p.index--
	p.next = p.bytes[p.index]
}

func (p *parser) error(s string) {
	panic(s)
}

func (p *parser) get() {
	if p.eof { return }

	p.index++
	if p.index >= p.length {
		// TODO more efficiently ? use buffered reader
		bytes := make([]byte, 1024)
		length, error := p.reader.Read(bytes)

		// XXX when p.reader is tls.Conn Read can return (0, nil)
		if length == 0 && error == nil {
			length, error = p.reader.Read(bytes)
		}
		bytes = bytes[:length]
		println("get: bytes =", string(bytes))
		// assert (e == os.EOF && p.l == 0)

		if error != nil {
			if error == os.EOF {
				p.eof = true
				return
			} else {
				panic(error)
			}
		}

		p.bytes = append(p.bytes, bytes...)
		p.length += length
	}

	p.next = p.bytes[p.index]
}

func (p *parser) mustGet() {
	p.get()
	if p.eof { p.error("unexpected EOF") }
}

