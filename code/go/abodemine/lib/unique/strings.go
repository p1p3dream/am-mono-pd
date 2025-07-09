package unique

func Alpha(l int) string {
	return string(AlphaB(l))
}

func AlphaNum(l int) string {
	return string(AlphaNumB(l))
}

func Bool() string {
	return string(BoolB())
}

func Digit(l int) string {
	return string(DigitB(l))
}

func Integer(l int) string {
	return string(IntegerB(l))
}

func Lower(l int) string {
	return string(LowerB(l))
}

func LowerNum(l int) string {
	return string(LowerNumB(l))
}

func Upper(l int) string {
	return string(UpperB(l))
}

func UpperNum(l int) string {
	return string(UpperNumB(l))
}

func StringSlice(l1, l2 int, f func(int) string) []string {
	s := make([]string, l1)

	for i := 0; i < l1; i++ {
		s[i] = f(l2)
	}

	return s
}
