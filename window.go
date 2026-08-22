package chunkcdc

import "github.com/LYH2263/go-chunkcdc/internal/clone"

type Window struct {
	data  []byte
	start int
	size  int
}

func NewWindow(data []byte, size int) *Window {
	if size <= 0 {
		size = 64
	}
	// Take a private copy so the window owns its bytes and is decoupled
	// from the caller's buffer. Reusing the caller's slice to fill the
	// next packet must not dirty the current window, or Peek/Hash drift
	// and CDC fingerprints stop matching (bad blocks on restore).
	return &Window{data: clone.Bytes(data), size: size}
}

func (w *Window) Advance() bool {
	if w.start+w.size >= len(w.data) {
		return false
	}
	w.start++
	return true
}

func (w *Window) Bytes() []byte {
	end := w.start + w.size
	if end > len(w.data) {
		end = len(w.data)
	}
	out := make([]byte, end-w.start)
	copy(out, w.data[w.start:end])
	return out
}

func (w *Window) Hash() uint32 {
	return WeakHash(w.Bytes())
}
