package iso8583

import (
	"io"

	"golang.org/x/text/transform"
)

// Reader wraps another io.Reader by transforming the bytes read.
type reader struct {
	r io.Reader
	t transform.Transformer
}

// Read implements the io.Reader interface.
//
//nolint:wrapcheck
func (r *reader) Read(p []byte) (int, error) {
	buf := make([]byte, len(p))

	n, err := r.r.Read(buf)
	if err != nil {
		return n, err
	}

	_, _, err = r.t.Transform(p, buf, true)

	return n, err
}

// Writer wraps another io.Writer by transforming the bytes written.
type writer struct {
	w io.Writer
	t transform.Transformer
}

// Write implements the io.Writer interface.
//
//nolint:wrapcheck
func (w *writer) Write(p []byte) (int, error) {
	buf := make([]byte, len(p))

	_, _, err := w.t.Transform(buf, p, true)
	if err != nil {
		return 0, err
	}

	return w.w.Write(buf)
}
