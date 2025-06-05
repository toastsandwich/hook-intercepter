package hookintercepter

import (
	"io"
)

type Hook struct {
	source io.Reader
	hook   io.Writer
}

func NewHook(r io.Reader, w io.Writer) *Hook {
	return &Hook{
		source: r,
		hook:   w,
	}
}

func (h *Hook) Seek(offset int64, whence int) (n int64, err error) {
	// Verify for source has embedded Seeker, use it.
	sourceSeeker, ok := h.source.(io.Seeker)
	if ok {
		return sourceSeeker.Seek(offset, whence)
	}
	// Verify if hook has embedded Seeker, use it.
	hookSeeker, ok := h.hook.(io.Seeker)
	if ok {
		return hookSeeker.Seek(offset, whence)
	}
	return n, nil
}

func (h *Hook) Read(b []byte) (n int, err error) {
	n, err = h.source.Read(b)
	if err != nil && err != io.EOF {
		return
	}
	if _, werr := h.hook.Write(b[:n]); werr != nil {
		if werr != io.EOF {
			return n, werr
		}
	}
	return
}
