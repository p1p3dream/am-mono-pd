package unique

var (
	digitB = []byte("0123456789")
	lowerB = []byte("abcdefghijklmnopqrstuvwxyz")
	upperB = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")

	digitLowerB = append(digitB, lowerB...)
	digitUpperB = append(digitB, upperB...)
	lowerUpperB = append(lowerB, upperB...)

	digitLowerUpperB = append(digitB, lowerUpperB...)
)

func AlphaB(l int) []byte {
	if l < 1 {
		return nil
	}
	return gen(l, lowerUpperB)
}

func AlphaNumB(l int) []byte {
	if l < 1 {
		return nil
	}
	return gen(l, digitLowerUpperB)
}

func BoolB() []byte {
	if gen(1, digitB[:2])[0] == '1' {
		return []byte("true")
	}
	return []byte("false")
}

func DigitB(l int) []byte {
	if l < 1 {
		return nil
	}
	return gen(l, digitB)
}

func IntegerB(l int) []byte {
	if l < 1 {
		return nil
	}
	if l == 1 {
		return gen(1, digitB)
	}
	return append(gen(1, digitB[1:]), gen(l-1, digitB)...)
}

func LowerB(l int) []byte {
	if l < 1 {
		return nil
	}
	return gen(l, lowerB)
}

func LowerNumB(l int) []byte {
	if l < 1 {
		return nil
	}
	return gen(l, digitLowerB)
}

func UpperB(l int) []byte {
	if l < 1 {
		return nil
	}
	return gen(l, upperB)
}

func UpperNumB(l int) []byte {
	if l < 1 {
		return nil
	}
	return gen(l, digitUpperB)
}

func BytesSlice(l1, l2 int, f func(int) []byte) [][]byte {
	s := make([][]byte, l1)

	for i := 0; i < l1; i++ {
		s[i] = f(l2)
	}

	return s
}
