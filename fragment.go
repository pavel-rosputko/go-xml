package xml

import (
	"bytes"
	"fmt"
)

type desc int64

func (d desc) tokenType() int { return int(d >> 60) }
func (d desc) depth() int { return int(d >> 52) & (1 << 8 - 1) }
func (d desc) off() int { return int(d) }
func (d desc) descLen() int { return int(d >> 32) & (1 << 20 - 1) }
func (d desc) descString() string { return fmt.Sprint(d.tokenType(), d.depth(), d.descLen(), d.off()) }

func (d *desc) setOffLen(offset, length int) {
	v := int64(*d)
	v &^= 1 << 32 - 1
	v |= int64(offset)

	v &^= (1 << 20 - 1) << 32
	v |= int64(length) << 32

	*d = desc(v)
}

// NOTE a: XmlPiece, Xml, IndexedXml, Markup, Map, Fragment, FragmentIndex, ModifiableIndex
// NOTE Index is a bad name, too ambigious
type Fragment struct {
	bytes	[]byte
	descs	[]desc
	nexts	[]int
	free	int
	root	int
	last	int
	mod	bool
}

func newFragment() *Fragment {
	return &Fragment{free: -1, root: -1, last: -1}
}

func (f *Fragment) add4(tokenType, depth int, m mark) {
	f.add(desc(int64(tokenType) << 60 | int64(depth) << 52 | int64(m.e - m.s) << 32 | int64(m.s)))
}

func (f *Fragment) getFree() int {
	if f.free == -1 {
		f.free = len(f.descs)
		f.descs = append(f.descs, 0)
		f.nexts = append(f.nexts, -1)
	}

	index := f.free
	f.free = f.nexts[f.free]

	return index
}

func (f *Fragment) add(value desc) {
	index := f.getFree()

	f.descs[index] = value
	f.nexts[index] = -1
	if f.last != -1 { f.nexts[f.last] = index }
	f.last = index
}

func (f *Fragment) addStart(marks []mark, depth int) {
	f.add4(startType, depth, marks[0])
	if depth == 0 { f.root = f.last } // NOTE move out of loop?

	for i := 1; i + 1 < len(marks); i += 2 {
		f.add4(keyType, depth, marks[i])
		f.add4(valueType, depth, marks[i + 1])
	}

	if len(marks) % 2 == 0 {
		f.add4(endType, depth, marks[0])
	}
}

func (f *Fragment) inspect() {
	i := f.root
	for i != -1 {
		fmt.Println(f.descs[i])
		i = f.nexts[i]
	}
}

func (f *Fragment) equalString(descIndex int, s string) bool {
	d := f.descs[descIndex]

	if d.descLen() != len(s) { return false }

	for i, j := d.off(), 0; j < len(s); i, j = i + 1, j + 1 {
		if f.bytes[i] != s[j] { return false }
	}

	return true
}

func (f *Fragment) AtString(i int) string {
	d := f.descs[i]
	return string(f.bytes[d.off() : d.off() + d.descLen()])
}

func (f *Fragment) AtBytes(i int) []byte {
	d := f.descs[i]
	return f.bytes[d.off() : d.off() + d.descLen()]
}

func (f *Fragment) Cursor() *Cursor {
	return &Cursor{f.iterator()}
}

func (f *Fragment) SetString(descIndex int, s string) {
	index := len(f.bytes)
	f.bytes = append(f.bytes, []byte(s)...)

	f.mod = true
	f.descs[descIndex].setOffLen(index, len(f.bytes) - index)
}

func (f *Fragment) String() string {
	if !f.mod { return string(f.bytes) }

	b := bytes.NewBuffer(nil)

	i := f.iterator()

	for {
		switch i.tokenType() {
		case startType:
			b.WriteByte('<')
			b.Write(i.bytes())

			i.next()
			for i.tokenType() == keyType || i.tokenType() == valueType {
				if i.tokenType() == keyType {
					b.WriteByte(' ')
					b.Write(i.bytes())
				} else {
					b.WriteByte('=')
					b.WriteByte('\'')
					b.Write(i.bytes())
					b.WriteByte('\'')
				}

				i.next()
			}

			if i.tokenType() == endType {
				b.WriteByte('/')
				b.WriteByte('>')

				i.next()
			} else {
				b.WriteByte('>')
			}
		case endType:
			b.WriteByte('<')
			b.Write(i.bytes())
			b.WriteByte('>')

			i.next()
		case charsType:
			b.Write(i.bytes())
			i.next()
		default:
			i.next()
		}

		if !i.hasNext() { break }
	}

	return b.String()
}

// NOTE insert at. other variant is to insert after descIndex
// but it required to maintain f.last variable as it may be changed
func (f *Fragment) insert(index int, value desc) {
	newIndex := f.getFree()

	f.descs[newIndex] = f.descs[index]
	f.nexts[newIndex] = f.nexts[index]

	f.descs[index] = value
	f.nexts[index] = newIndex
}

func (f *Fragment) insert4(index int, tokenType, depth int, m mark) {
	f.insert(index, desc(int64(tokenType) << 60 | int64(depth) << 52 | int64(m.e - m.s) << 32 | int64(m.s)))
}

func (f *Fragment) appendBytes(bytes []byte) mark {
	index := len(f.bytes)
	f.bytes = append(f.bytes, bytes...)
	return mark{index, len(f.bytes)}
}

func (f *Fragment) insertAttr(descIndex int, key, value string) {
	mark := f.appendBytes([]byte(value))
	f.insert4(descIndex, valueType, 0, mark)
	mark = f.appendBytes([]byte(key))
	f.insert4(descIndex, keyType, 0, mark)
	f.mod = true
}

func (f *Fragment) iterator() iterator {
	return iterator{f, f.root, f.descs[f.root]}
}
