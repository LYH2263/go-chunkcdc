package chunkcdc

import (
	"bytes"
	"testing"
)

// TestWindow_OwnsBytes regresses the original bug: a Window built from a
// caller buffer must not change when the caller mutates that buffer to fill
// the next packet. Before the fix, NewWindow aliased the caller slice, so
// Peek/Hash drifted and CDC fingerprints stopped matching.
func TestWindow_OwnsBytes(t *testing.T) {
	orig := []byte("the quick brown fox jumps over the lazy dog")
	w := NewWindow(orig, 8)

	want := append([]byte(nil), w.Bytes()...)
	wantHash := w.Hash()

	// Caller reuses the underlying buffer for the next packet.
	for i := range orig {
		orig[i] = 0xAA
	}

	got := w.Bytes()
	if !bytes.Equal(got, want) {
		t.Fatalf("window mutated by caller buffer reuse: got %x, want %x", got, want)
	}
	if h := w.Hash(); h != wantHash {
		t.Fatalf("window hash drifted after caller buffer reuse: got %d, want %d", h, wantHash)
	}
}

// TestSession_Ingest_BufferReuse mirrors the backup-agent scenario end to end:
// Ingest a packet, reuse the same buffer for the next packet, and confirm the
// session's window and fingerprints are stable.
func TestSession_Ingest_BufferReuse(t *testing.T) {
	s, err := OpenSession(SessionOptions{WindowSize: 8})
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	buf := []byte("packet-zero-content-here")

	if err := s.Ingest(buf); err != nil {
		t.Fatal(err)
	}
	peek0 := s.PeekWindow()
	hash0 := s.WindowHash()
	fp0 := Fingerprint(peek0)

	// Backup agent "saves a copy" by reusing buf for the next packet.
	for i := range buf {
		buf[i] = byte('Z')
	}

	if got := s.PeekWindow(); !bytes.Equal(got, peek0) {
		t.Fatalf("PeekWindow drifted after caller reused buffer: got %q, want %q", got, peek0)
	}
	if h := s.WindowHash(); h != hash0 {
		t.Fatalf("WindowHash drifted after caller reused buffer: got %d, want %d", h, hash0)
	}
	if fp := Fingerprint(s.PeekWindow()); fp != fp0 {
		t.Fatalf("window fingerprint drifted after caller reused buffer: got %s, want %s", fp, fp0)
	}
}
