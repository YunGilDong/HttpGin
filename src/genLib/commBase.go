package genLib

const (
	ED_BIG    = 1
	ED_LITTLE = 2
)

func GetNumber(src []byte, pos int, length int, endian int) int {
	var value int = 0

	switch endian {
	case ED_BIG:
		for idx := 0; idx < length; idx++ {
			value = (value * 256) + int(src[pos+idx])
		}
		break
	case ED_LITTLE:
		for idx := length - 1; idx >= 0; idx-- {
			value = (value * 256) + int(src[pos+idx])
		}
		break
	}

	return value
}

func SetNumber(ptr []*byte, value int, length int, endian int) {
	switch endian {
	case ED_BIG:
		for idx := length - 1; idx >= 0; idx-- {
			*ptr[idx] = byte(value % 256)
			value /= 256
		}
		break
	case ED_LITTLE:
		for idx := 0; idx < length; idx++ {
			*ptr[idx] = byte(value % 256)
			value /= 256
		}
		break
	}
}
