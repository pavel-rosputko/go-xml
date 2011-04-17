package xml

type Cursor struct {
	iterator
}

func (c *Cursor) ToChild() (f bool) {
	i := c.iterator

	i.next()
	for i.tokenType() != startType {
		i.next()
	}

	if i.tokenType() == startType && i.depth() == c.depth() + 1 {
		c.iterator = i
		f = true
	}

	return
}

func (c *Cursor) MustToChild() {
	if !c.ToChild() { panic("ToChild false") }
}

func (c *Cursor) Name() string {
	return c.string()
}

func (c *Cursor) Attr(name string) (v string, f bool) {
	i := c.iterator

	i.next()
	for i.tokenType() == keyType && !i.equalString(name) {
		i.next(); i.next()
		println("i.s =", i.string())
		println(i.tokenType())
		println(i.equalString(name))
	}

	println(name)
	if i.tokenType() == keyType {
		i.next()
		v, f = i.string(), true
	}

	return
}

func (c *Cursor) MustAttr(name string) string {
	v, f := c.Attr(name)
	if !f { panic("Attr false") }
	return v
}

// NOTE return just first found chars
func (c *Cursor) Chars() (v string, f bool) {
	i := c.iterator

	i.next()
	for i.tokenType() == keyType || i.tokenType() == valueType {
		i.next()
	}

	if i.tokenType() == charsType {
		v, f = i.string(), true
	}

	return
}

func (c *Cursor) MustChars() string {
	v, f := c.Chars()
	if !f { panic("Chars false") }
	return v
}

func (c *Cursor) SetAttr(key, value string) {
	i := c.iterator

	i.next()
	for i.hasNext() && i.tokenType() == keyType && !i.equalString(key) {
		i.next(); i.next()
	}

	if i.tokenType() == keyType {
		i.next()
		i.setString(value)
	} else {
		i.insertAttr(key, value)
	}
}

func (c *Cursor) ChildrenString() string {
	i := c.iterator

	i.next()
	for i.tokenType() != startType {
		i.next()
	}

	startIndex := i.desc.off() - 1 // <
	for i.depth() != c.depth() {
		i.next()
	}
	// has children so should has end tag
	endIndex := i.desc.off() - 2 // </

	return string(c.Fragment.bytes[startIndex:endIndex])
}

// func (c *Cursor) AttrIndex(s string) int
// func (c *Cursor) MatchString(i int, s string) bool
// func (c *Cursor) TextInt() int // TextAsInt
// ToChild, ToFirstChild, ToLastChild, ToSibling, ToNextSibling, toPrevSibling
// func (c *Cursor) ToChildS(s string)


