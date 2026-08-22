package chunkcdc
type Window struct {
        data []byte
        start int
        size int
}
func NewWindow(data []byte, size int) *Window {
        if size <= 0 { size = 64 }
        return &Window{data: data, size: size}
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
        if end > len(w.data) { end = len(w.data) }
        out := make([]byte, end-w.start)
        copy(out, w.data[w.start:end])
        return out
}
func (w *Window) Hash() uint32 {
        return WeakHash(w.Bytes())
}
