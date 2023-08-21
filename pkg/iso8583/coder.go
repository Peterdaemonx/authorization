package iso8583

type coder struct {
	format format
	layout layout

	mtiFormat dataEncoding
	lenFormat lengthEncoding
}

func applyOpts(c coder, opts ...Opt) coder {
	for _, opt := range opts {
		c = opt(c)
	}
	return c
}

type Opt func(coder) coder

func FormatEbcdic(c coder) coder {
	c.format = EBCDIC
	return c
}

func FormatAscii(c coder) coder {
	c.format = ASCII
	return c
}

func RdwLayout(c coder) coder {
	c.layout = WithRDW
	return c
}

func NoRdwLayout(c coder) coder {
	c.layout = NoRDW
	return c
}

func MtiBcd4(c coder) coder {
	c.mtiFormat = DataEncodingBcd4
	return c
}

func HexLen(c coder) coder {
	c.lenFormat = LengthEncodingHex
	return c
}
