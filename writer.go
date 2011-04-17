package xml

import (
	"bufio"
	"io"
)

// NOTE It' s just piece of xml writer
// TODO escaping
type Writer struct {
	*bufio.Writer
	stack
}

func NewWriter(ioWriter io.Writer) *Writer {
	return &Writer{Writer: bufio.NewWriter(ioWriter)}
}

// NOTE ? add these methods to handle errors or not 
func (w *Writer) writeByte(b byte) {
	e := w.Writer.WriteByte(b)
	if e != nil { panic(e) }
}

func (w *Writer) writeString(s string) {
	_, e := w.Writer.WriteString(s)
	if e != nil { panic(e) }
}

func (w *Writer) flush() {
	e := w.Writer.Flush()
	if e != nil { panic(e) }
}

func (w *Writer) StartDocument() *Writer {
	w.writeString("<?xml version='1.0' encoding='utf-8'?>")
	return w
}

func (w *Writer) attributes(args []string) {
	for i := 0; i + 1 < len(args); i += 2 {
		if args[i + 1] == "" { continue }
		w.writeByte(' ')
		w.writeString(args[i])
		w.writeByte('=')
		w.writeByte('\'') // TODO add separator option ?
		w.writeString(args[i + 1])
		w.writeByte('\'')
	}
}

// NOTE provide method with []byte? without check for ""
func (w *Writer) StartElement(name string, args ...string) *Writer {
	w.writeByte('<')
	w.writeString(name)
	w.attributes(args)
	w.writeByte('>')

	w.push(name)

	return w
}

func (w *Writer) EndElement() *Writer {
	name := w.pop()

	w.writeByte('<')
	w.writeByte('/')
	w.writeString(name)
	w.writeByte('>')

	return w
}

func (w *Writer) Element(name string, args ...string) *Writer {
	w.writeByte('<')
	w.writeString(name)

	w.attributes(args)

	if len(args) % 2 != 0 {
		w.writeByte('>')
		w.writeString(args[len(args) - 1])
		w.writeByte('<')
		w.writeByte('/')
		w.writeString(name)
		w.writeByte('>')
	} else {
		w.writeString(" />")
	}

	return w
}

func (w *Writer) EndDocument() {
	for !w.isEmpty() { w.EndElement() }
	w.Flush()
}

// NOTE close all open elements and flush buffer TODO ? improve name
func (w *Writer) End() {
	for !w.isEmpty() { w.EndElement() }
	w.Flush()
}

func (w *Writer) Send() {
	w.stack.clean()
	w.Flush()
}

func (w *Writer) Raw(s string) *Writer {
	w.writeString(s)
	return w
}

type stack []string

func (s *stack) push(v string) { *s = append(*s, v) }

func (s *stack) pop() (v string) {
	v, *s = (*s)[len(*s) - 1], (*s)[0:len(*s) - 1]
	return
}

func (s stack) isEmpty() bool { return len(s) == 0 }
func (s *stack) clean() { *s = (*s)[0:0] }


