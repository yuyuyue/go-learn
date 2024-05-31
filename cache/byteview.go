package cache

type ByteView struct {
	b []byte // 使用btye可以支持所有类型
}

func (b ByteView) Len() int {
	return len(b.b)
}

func (b ByteView) BtyesSlice() []byte {
	return copyBytes(b.b)
}

func copyBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}

func (b ByteView) String() string {
	return string(b.b)
}
