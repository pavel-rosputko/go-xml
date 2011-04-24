package xml

import (
	"bytes"
)

type Builder struct {
	*bytes.Buffer
	*Fragment
	stringStack
}

func NewBuilder() *Builder {
	return &Builder{
		Buffer: bytes.NewBuffer([]byte{}),
		Fragment: newFragment()}
}

func (b *Builder) writeMarked(s string) Mark {
	startIndex := b.Buffer.Len()
	b.WriteString(s)
	return Mark{startIndex, b.Buffer.Len()}
}

func (b *Builder) StartElement(name string, args ...string) *Builder {
	b.WriteByte('<')

	mark := b.writeMarked(name)
	b.add4(startType, b.depth(), mark)

	b.attributes(args)

	b.WriteByte('>')

	b.push(name)

	return b
}

func (b *Builder) depth() int {
	return len(b.stringStack)
}

func (b *Builder) attributes(args []string) {
	for i := 0; i + 1 < len(args); i += 2 {
		if args[i + 1] == "" { continue }

		b.WriteByte(' ')

		mark := b.writeMarked(args[i])
		b.add4(keyType, b.depth(), mark)

		b.WriteByte('=')
		b.WriteByte('\'')

		mark = b.writeMarked(args[i + 1])
		b.add4(valueType, b.depth(), mark)

		b.WriteByte('\'')
	}
}

func (b *Builder) EndElement() *Builder {
	name := b.pop()

	b.WriteByte('<')
	b.WriteByte('/')
	mark := b.writeMarked(name)
	b.add4(endType, b.depth(), mark)

	b.WriteByte('>')

	return b
}

func (b *Builder) Element(name string, args ...string) *Builder {
	b.WriteByte('<')

	mark := b.writeMarked(name)
	b.add4(startType, b.depth(), mark)

	b.attributes(args)

	if len(args) % 2 != 0 {
		b.WriteByte('>')

		mark = b.writeMarked(args[len(args) - 1])
		b.add4(charsType, b.depth(), mark)

		b.WriteByte('<')
		b.WriteByte('/')

		mark := b.writeMarked(name)
		b.add4(endType, b.depth() + 1, mark)

		b.WriteByte('>')
	} else {
		// NOTE write end tag or not ?
		b.WriteString(" />")
	}

	return b
}

func (b *Builder) End() *Fragment {
	for !b.isEmpty() { b.EndElement() }

	b.Fragment.bytes = b.Buffer.Bytes()
	return b.Fragment
}

func (b *Builder) Append(f *Fragment) *Builder {
	offset := b.Len()
	depth := b.depth()

	b.Write(f.bytes)

	for _, desc := range f.descs {
		desc.setDepth(depth + desc.depth())
		desc.setOff(offset + desc.off())
		b.add(desc)
	}

	return b
}
