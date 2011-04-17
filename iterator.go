package xml

type iterator struct {
	*Fragment
	index	int
	desc
}

func (i *iterator) hasNext() bool {
	return i.nexts[i.index] != -1
}

func (i *iterator) next() {
	println("next: i.index =", i.index, i.string())
	i.index = i.nexts[i.index]
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

func (i *iterator) setString(s string) {
	i.Fragment.SetString(i.index, s)
}

func (i *iterator) insertAttr(key, value string) {
	i.Fragment.insertAttr(i.index, key, value)
}

