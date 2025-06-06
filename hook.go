package hookintercepter

import (
	"errors"
	"fmt"
	"io"
)

var ErrSeekUnsupported = errors.New("seek is not supported")

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
	return n, nil
}

func (h *Hook) Read(b []byte) (n int, err error) {
	n, err = h.source.Read(b)
	if err != nil && err != io.EOF {
		return
	}
	log := fmt.Sprintf("progress bar has read %d bytes\n.", n)
	if _, werr := h.hook.Write([]byte(log)); werr != nil {
		if werr != io.EOF {
			return n, werr
		}
	}
	return
}
