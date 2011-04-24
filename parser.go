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
	index	int
	bytes	[]byte
	buffer	[]byte
}

func newParser(reader io.Reader) *parser {
	return &parser{
		reader: reader,
		buffer: make([]byte, 1024),
	}
}

type mark struct { s, e int }

// NOTE eofType instead of f bool ?
func (p *parser) token() (tokenType int, marks []mark, f bool) {
	b, e := p.get()
	if e { return }

	f = true
	if b != '<' { // text
		p.back()
		tokenType, marks = charsType, []mark{p.markUntilEOFOr('<')}
		return
	}

	switch p.mustGet() {
	case '/': // end tag
		tokenType, marks = endType, []mark{p.name()}

		if p.mustGetNonSpace() != '>' {
			p.error("invalid characters between </ and >")
		}
	case '?': // directive
		index := p.index
		var  bb, b byte
		for !(bb == '?' && b == '>') {
			bb, b = b, p.mustGet()
		}

		tokenType, marks = directiveType, []mark{{index, p.index - 2}}
	case '!': // comment or cdata
		switch p.mustGet() {
		case '-':
			if p.mustGet() != '-' {
				p.error("invalid sequence <!- not part of <!--")
			}

			index := p.index
			var b3, b2, b1 byte
			for !(b3 == '-' && b2 == '-' && b1 == '>') {
				b3, b2, b1 = b2, b1, p.mustGet()
			}

			tokenType, marks = commentType, []mark{{index, p.index - 3}}
		case '[':
			for i := 0; i < 6; i++ {
				if p.mustGet() != "CDATA"[i] {
					p.error("invalid <![ sequence")
				}
			}

			index := p.index
			var bb, b byte
			for !(bb == ']' && b == ']') {
				bb, b = b, p.mustGet()
			}

			tokenType, marks = cdataType, []mark{{index, p.index - 2}}
		default:
			// probably a directive <!
		}
	default: // start tag
		p.back()
		tokenType, marks = startType, []mark{p.name()}

		var empty bool
		for {
			b := p.mustGetNonSpace()
			if b == '/' { // empty tag
				empty = true
				if p.mustGet() != '>' {
					p.error("expected /> in element")
				}
				break
			} else if b == '>' {
				break
			}
			p.back()

			marks = append(marks, p.name())

			if p.mustGetNonSpace() != '='  {
				p.error("attribute name without = in element")
			}

			b = p.mustGetNonSpace()
			if b != '\'' && b != '"' {
				p.error("unquoted attribute value")
			}

			marks = append(marks, p.markUntil(b))

			p.mustGet()
		}

		if empty {
			marks = append(marks, mark{})
		}
	}

	return
}

func (p *parser) sliceBytes() (bytes []byte) {
	bytes, p.bytes = p.bytes[:p.index], p.bytes[p.index:]

	p.index = 0
	return
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

func (p *parser) back() {
	p.index--
}

func (p *parser) error(s string) {
	panic(s)
}

func (p *parser) read() bool {
	// TODO more efficiently ? use buffered reader
	l, e := p.reader.Read(p.buffer)

	// XXX when p.reader is tls.Conn Read can return (0, nil)
	if l == 0 && e == nil {
		l, e = p.reader.Read(p.buffer)
	}
	buffer := p.buffer[:l]
	// buf = buf[:length]
	// println("get: bytes =", string(bytes))
	// assert (e == os.EOF && p.l == 0)

	if e != nil {
		if e == os.EOF { return false } else { panic(e) }
	}

	p.bytes = append(p.bytes, buffer...)

	return true
}

func (p *parser) mustRead() {
	if !p.read() {
		p.error("unexpected EOF")
	}
}

func (p *parser) get() (b byte, e bool) {
	if p.index == len(p.bytes) {
		if !p.read() {
			e = true
			return
		}
	}

	b = p.bytes[p.index]
	p.index++
	return
}

func (p *parser) mustGet() byte {
	b, e := p.get()
	if e { p.error("unexpected EOF") }
	return b
}

func isSpace(b byte) bool {
	return b == ' ' || b == '\r' || b == '\n' || b == '\t'
}

func (p *parser) mustGetNonSpace() (b byte) {
	for {
		if p.index == len(p.bytes) {
			p.mustRead()
		}

		b = p.bytes[p.index]
		p.index++

		if !isSpace(b) {
			break
		}
	}
	return
}

func (p *parser) markUntil(b byte) (m mark) {
	m.s = p.index
	for {
		if p.index == len(p.bytes) {
			p.mustRead()
		}

		if p.bytes[p.index] == b {
			break
		}

		p.index++
	}

	m.e = p.index
	return
}

func (p *parser) markUntilEOFOr(b byte) (m mark) {
	m.s = p.index
	for {
		if p.index == len(p.bytes) {
			if !p.read() {
				break
			}
		}

		if p.bytes[p.index] == b {
			break
		}

		p.index++
	}

	m.e = p.index
	return
}

func isNameByte(b byte) bool {
	return 'A' <= b && b <= 'Z' ||
		'a' <= b && b <= 'z' ||
		'0' <= b && b <= '9' ||
		b == '_' || b == ':' || b == '.' || b == '-'
}

func (p *parser) name() (m mark) {
	m.s = p.index
	if b := p.mustGet(); !(b >= utf8.RuneSelf || isNameByte(b)) {
		p.error("invalid tag name first letter")
	}

	for {
		if p.index == len(p.bytes) {
			p.mustRead()
		}

		if b := p.bytes[p.index]; !(b >= utf8.RuneSelf || isNameByte(b)) {
			break
		}

		p.index++

	}

	m.e = p.index
	return

	// TODO check the characters in [i:p.index]
}
