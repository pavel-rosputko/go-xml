package xml

import (
	"fmt"
)

// NOTE a: XmlPiece, Xml, IndexedXml, Markup, Map, Fragment, FragmentIndex, ModifiableIndex
// NOTE Index is a bad name, too ambigious
type Fragment struct {
	bytes	[]byte
	descs	[]desc
}

func newFragment() *Fragment {
	return &Fragment{
		bytes:	[]byte{},
		descs:	[]desc{},
	}
}

func (f *Fragment) add4(tokenType, depth int, m Mark) {
	f.add(desc(int64(tokenType) << 60 | int64(depth) << 52 | int64(m.E - m.S) << 32 | int64(m.S)))
}

func (f *Fragment) add(value desc) {
	f.descs = append(f.descs, value)
}

func (f *Fragment) addStart(mark Mark, marks []Mark, depth int) {
	f.add4(startType, depth, mark)

	for i := 0; i + 1 < len(marks); i += 2 {
		f.add4(keyType, depth, marks[i])
		f.add4(valueType, depth, marks[i + 1])
	}

	if len(marks) % 2 != 0 {
		f.add4(endType, depth, mark)
	}
}

func (f *Fragment) inspect() {
	for _, desc := range f.descs {
		fmt.Println(desc.descString())
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
