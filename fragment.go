package xml

import (
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
}

func newFragment() *Fragment {
	return &Fragment{}
}

func (f *Fragment) add4(tokenType, depth int, m mark) {
	f.add(desc(int64(tokenType) << 60 | int64(depth) << 52 | int64(m.e - m.s) << 32 | int64(m.s)))
}

func (f *Fragment) add(value desc) {
	f.descs = append(f.descs, value)
}

func (f *Fragment) addStart(marks []mark, depth int) {
	f.add4(startType, depth, marks[0])

	for i := 1; i + 1 < len(marks); i += 2 {
		f.add4(keyType, depth, marks[i])
		f.add4(valueType, depth, marks[i + 1])
	}

	if len(marks) % 2 == 0 {
		f.add4(endType, depth, marks[0])
	}
}

func (f *Fragment) inspect() {
	for desc := range f.descs {
		fmt.Println(desc)
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
func (f *Fragment) String() string {
	return string(f.bytes)
}

func (f *Fragment) iterator() iterator {
	return iterator{f, 0, f.descs[0]}
}
