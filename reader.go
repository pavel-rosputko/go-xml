package xml

import (
	// "fmt"
	"io"
)

type marksStack []mark

func (s *marksStack) push(m mark) { *s = append(*s, m) }
func (s marksStack) top() mark { return s[len(s) - 1] }
func (sp *marksStack) pop() (m mark, f bool) {
	s := *sp
	if len(s) != 0 {
		m, f = s[len(s) - 1], true
		*sp = s[:len(s) - 1]
	}

	return
}

func (s marksStack) isEmpty() bool { return len(s) == 0 }

// NOTE parser with * or not ?
type Reader struct {
	*parser
	marksStack
}

func NewReader(ioReader io.Reader) *Reader {
	return &Reader{parser: newParser(ioReader)}
}

// TODO return Fragment or mark ?
func (r *Reader) ReadStartElement() *Fragment {
	fragment := newFragment()

	tt, m, f := r.token()
	for {
		if !f { r.error("no tokens") }
		if tt == startType { break }
		tt, m, f = r.token()
	}

	fragment.addStart(m, 0)

	fragment.add4(eofType, 0, mark{})

	fragment.bytes = r.sliceBytes()

	// TODO add satellite end !?

	return fragment
}

// FIXME when chars exists before first start-element they are ignore in descs but exist in bytes, remove them ?
func (r *Reader) ReadElement() *Fragment {
	fragment := newFragment()

	depth := 0
	done := false
	for !done {
		tokenType, marks, f := r.token()
		if !f { r.error("unexpected eof") }

		switch tokenType {
		case startType:
			// fmt.Println("startType", marks)
			fragment.addStart(marks, depth)

			if len(marks) % 2 != 0 {
				r.push(marks[0])
				depth++
			} else {
				if r.isEmpty() { done = true }
			}
		case endType:
			depth--
			mark, f := r.pop()
			if !f { r.error("unexpected end tag") }
			if !r.markEq(marks[0], mark) { panic("wrong end tag name") }
			fragment.add4(endType, depth, marks[0])

			if r.isEmpty() { done = true }
		case charsType:
			// fmt.Println("charsType", marks)

			// skip pre- and after-element chars
			if depth == 0 { continue }
			fragment.add4(charsType, depth, marks[0])
		case cdataType:
			fragment.add4(cdataType, depth, marks[0])
		}
	}

	fragment.add4(eofType, 0, mark{len(r.bytes), len(r.bytes)})

	fragment.bytes = r.sliceBytes()

	return fragment
}

func (r *Reader) error(s string) {
	panic(s)
}

