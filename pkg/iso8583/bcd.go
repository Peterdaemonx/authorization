package iso8583

import (
	"io"

	"github.com/yerden/go-util/bcd"
)

type bcdEncoder struct {
	w io.Writer
	r io.Reader
}

func (be bcdEncoder) Write(p []byte) (n int, err error) {
	if len(p)%2 != 0 {
		p = append([]byte("0"), p...)
	}

	enc := bcd.NewEncoder(bcd.Standard)
	dst := make([]byte, bcd.EncodedLen(len(p)))
	_, err = enc.Encode(dst, p)
	if err != nil {
		return 0, err
	}

	return be.w.Write(dst)
}

func (be bcdEncoder) Read(p []byte) (n int, err error) {
	// how many bytes we will read
	bcdLen := bcd.EncodedLen(len(p))
	bcdBuf := make([]byte, bcdLen)
	pBuf := make([]byte, bcdLen*2)

	bcdN, err := be.r.Read(bcdBuf)
	if err != nil {
		return bcdN, err
	}

	dec := bcd.NewDecoder(bcd.Standard)
	// Length of the fields needs to be read binary
	// Data of the fields needs to be read with bcd4 decoding.
	// We need to split these two.
	n, err = dec.Decode(pBuf, bcdBuf)
	if err != nil {
		return n, err
	}

	// because BCD is right aligned, we skip first bytes and
	// read only what we need
	// e.g. 0643 => 643
	copy(p, pBuf[len(pBuf)-len(p):])
	return n, nil
}
