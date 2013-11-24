package edifact

type Values []interface{}

type Header [2]string

func (h Header) ComponentDelimiter() byte {
	return h[1][0]
}

func (h Header) DataDelimiter() byte {
	return h[1][1]
}

func (h Header) Decimal() byte {
	return h[1][2]
}

func (h Header) ReleaseIndicator() byte {
	return h[1][3]
}

func (h Header) RepetitionDelimiter() byte {
	return h[1][4]
}

func (h Header) SegmentTerminator() byte {
	return h[1][5]
}
