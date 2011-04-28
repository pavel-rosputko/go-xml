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

const (
	noneState = iota
	startState
	keyState
)

// NOTE a: Indexer
// NOTE: use bytes.Buffer instead of bytes ?
type parser struct {
	reader	io.Reader
	index	int
	bytes	[]byte
	buffer	[]byte
}

func makeParser(reader io.Reader) parser {
	return parser{
		reader: reader,
		buffer: make([]byte, 1024),
	}
}

type Mark struct { S, E int }

// NOTE eofType instead of f bool ?
func (p *parser) token() (tokenType int, mark Mark, marks []Mark, f bool) {
	b, e := p.get()
	if e { return }

	f = true
	if b != '<' {
		p.back()
		tokenType, mark = charsType, p.markUntilEOFOr('<')
		return
	}

	switch p.mustGet() {
	case '/':
		tokenType, mark = endType, p.name()

		if p.mustGetNonSpace() != '>' {
			p.error("invalid characters between </ and >")
		}
	case '?':
		index := p.index
		var  bb, b byte
		for !(bb == '?' && b == '>') {
			bb, b = b, p.mustGet()
		}

		tokenType, mark = directiveType, Mark{index, p.index - 2}
	case '!':
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

			tokenType, mark = commentType, Mark{index, p.index - 3}
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

			tokenType, mark = cdataType, Mark{index, p.index - 2}
		default:
		}
	default:
		p.back()
		tokenType, mark = startType, p.name()

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
			marks = append(marks, Mark{})
		}
	}

	return
}

func (p *parser) sliceBytes() (bytes []byte) {
	bytes, p.bytes = p.bytes[:p.index], p.bytes[p.index:]
	p.index = 0
	return
}

func (p *parser) markEq(m1, m2 Mark) bool {
	if m1.E - m1.S != m2.E - m2.S { return false }

	i, j := m1.S, m2.S
	for i < m1.E {
		if p.bytes[i] != p.bytes[j] { return false }
		i++; j++
	}

	return true
}

func (p *parser) string(m Mark) string {
	return string(p.bytes[m.S:m.E])
}

func (p *parser) back() {
	p.index--
}

type ParserError string

func (p *parser) error(s string) {
	panic(ParserError(s))
}

// TODO more efficiently ? use bufio.Reader
func (p *parser) read() bool {
	l, e := p.reader.Read(p.buffer)

	// XXX when p.reader is tls.Conn Read can return (0, nil)
	if l == 0 && e == nil {
		l, e = p.reader.Read(p.buffer)
	}
	// buf = buf[:length]
	// println("get: bytes =", string(bytes))
	// assert (e == os.EOF && p.l == 0)

	if e != nil {
		if e == os.EOF {
			return false
		} else {
			panic(e)
		}
	}

	p.bytes = append(p.bytes, p.buffer[:l]...)

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

func (p *parser) markUntil(b byte) (m Mark) {
	m.S = p.index
	for {
		if p.index == len(p.bytes) {
			p.mustRead()
		}

		if p.bytes[p.index] == b {
			break
		}

		p.index++
	}

	m.E = p.index
	return
}

func (p *parser) markUntilEOFOr(b byte) (m Mark) {
	m.S = p.index
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

	m.E = p.index
	return
}

func isNameByte(b byte) bool {
	return 'A' <= b && b <= 'Z' ||
		'a' <= b && b <= 'z' ||
		'0' <= b && b <= '9' ||
		b == '_' || b == ':' || b == '.' || b == '-'
}

func (p *parser) name() (m Mark) {
	m.S = p.index
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

	m.E = p.index
	return

	// TODO check the characters in [i:p.index]
}
