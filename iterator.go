package xml

type iterator struct {
	*Fragment
	index	int
	desc
}

func (i *iterator) hasNext() bool {
	return i.index + 1 < len(i.descs)
}

func (i *iterator) next() {
	i.index++
	i.desc = i.descs[i.index]
}

func (i *iterator) string() string {
	return i.Fragment.AtString(i.index)
}

func (i *iterator) bytes() []byte {
	return i.Fragment.AtBytes(i.index)
}

func (i *iterator) equalString(s string) bool {
	return i.Fragment.equalString(i.index, s)
}
