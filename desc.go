package xml

import (
	"fmt"
)

type desc int64

func (d desc) tokenType() int {
	return int(d >> 60)
}

func (d desc) depth() int {
	return int(d >> 52) & (1 << 8 - 1)
}

func (d desc) off() int {
	return int(d)
}

func (d desc) descLen() int {
	return int(d >> 32) & (1 << 20 - 1)
}

func (d desc) descString() string {
	return fmt.Sprint(d.tokenType(), d.depth(), d.descLen(), d.off())
}

func (d *desc) setOffLen(offset, length int) {
	v := int64(*d)
	v &^= 1 << 32 - 1
	v |= int64(offset)

	v &^= (1 << 20 - 1) << 32
	v |= int64(length) << 32

	*d = desc(v)
}

func (d *desc) setOff(offset int) {
	v := int64(*d)

	v &^= 1 << 32 - 1
	v |= int64(offset)

	*d = desc(v)
}

func (d *desc) setDepth(depth int) {
	v := int64(*d)

	v &^= (1 << 8 - 1) << 52
	v |= int64(depth) << 52

	*d = desc(v)
}

