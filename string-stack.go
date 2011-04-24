package xml

type stringStack []string

func (s *stringStack) push(v string) {
	*s = append(*s, v)
}

func (s *stringStack) pop() (v string) {
	v, *s = (*s)[len(*s) - 1], (*s)[:len(*s) - 1]
	return
}

func (s stringStack) isEmpty() bool {
	return len(s) == 0
}

func (s *stringStack) clean() {
	*s = (*s)[0:0]
}
