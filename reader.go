package xml

import (
	// "fmt"
	"io"
)

type marksStack []Mark

func (s *marksStack) push(m Mark) { *s = append(*s, m) }
func (s marksStack) top() Mark { return s[len(s) - 1] }
func (sp *marksStack) pop() (m Mark, f bool) {
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

	tt, m, mm, f := r.token()
	for {
		if !f { r.error("no tokens") }
		if tt == startType { break }
		tt, m, mm, f = r.token()
	}

	fragment.addStart(m, mm, 0)

	fragment.add4(eofType, 0, Mark{})

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
		tokenType, mark, marks, f := r.token()
		if !f { r.error("unexpected eof") }

		switch tokenType {
		case startType:
			// fmt.Println("startType", marks)
			fragment.addStart(mark, marks, depth)

			if len(marks) % 2 == 0 {
				r.push(mark)
				depth++
			} else {
				if r.isEmpty() { done = true }
			}
		case endType:
			depth--
			startMark, f := r.pop()
			if !f { r.error("unexpected end tag") }
			if !r.markEq(mark, startMark) { panic("wrong end tag name") }
			fragment.add4(endType, depth, mark)

			if r.isEmpty() { done = true }
		case charsType:
			// fmt.Println("charsType", marks)

			// skip pre- and after-element chars
			if depth == 0 { continue }
			fragment.add4(charsType, depth, mark)
		case cdataType:
			fragment.add4(cdataType, depth, mark)
		}
	}

	fragment.add4(eofType, 0, Mark{len(r.bytes), len(r.bytes)})

	fragment.bytes = r.sliceBytes()

	return fragment
}

func (r *Reader) error(s string) {
	panic(s)
}

