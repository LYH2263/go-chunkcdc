package chunkcdc

import (
	"errors"
	"testing"
)

// A hasher that records whether it was consulted.
type recordingHasher struct {
	called bool
}

func (r *recordingHasher) Sum32(data []byte) uint32 {
	r.called = true
	return WeakHash(data)
}

// Missing hasher must surface as a decidable ErrNoHasher, never a nil-deref
// panic, and must leave the visible set and dedup index untouched.
func TestIngestSplit_MissingHasher_NoDirtyRegistration(t *testing.T) {
	s, err := OpenSession(SessionOptions{})
	if err != nil {
		t.Fatalf("OpenSession: %v", err)
	}
	defer s.Close()

	data := []byte("the quick brown fox")
	before := s.VisibleCount()

	// No hasher injected — this must not panic.
	err = s.IngestSplit(data)
	if !errors.Is(err, ErrNoHasher) {
		t.Fatalf("want ErrNoHasher, got %v", err)
	}

	if got := s.VisibleCount(); got != before {
		t.Fatalf("VisibleCount changed on failed ingest: before=%d after=%d", before, got)
	}

	if _, ok := s.LookupFP(Fingerprint(data)); ok {
		t.Fatalf("dedup index polluted with half-baked entry for failed ingest")
	}
}

// With a hasher present the block is registered complete (Hash set) and
// reachable through the dedup index.
func TestIngestSplit_WithHasher_RegistersCompleteBlock(t *testing.T) {
	h := &recordingHasher{}
	s, err := OpenSession(SessionOptions{Hasher: h})
	if err != nil {
		t.Fatalf("OpenSession: %v", err)
	}
	defer s.Close()

	data := []byte("lazy dog")
	before := s.VisibleCount()

	if err := s.IngestSplit(data); err != nil {
		t.Fatalf("IngestSplit: %v", err)
	}
	if !h.called {
		t.Fatal("hasher was not consulted")
	}
	if got := s.VisibleCount(); got != before+1 {
		t.Fatalf("VisibleCount = %d, want %d", got, before+1)
	}

	fp := Fingerprint(data)
	info, ok := s.LookupFP(fp)
	if !ok {
		t.Fatal("registered block not found via FP")
	}
	if info.Hash == 0 {
		t.Fatal("registered block has zero Hash — half-baked registration")
	}
	if info.Hash != WeakHash(data) {
		t.Fatalf("Hash = %d, want %d", info.Hash, WeakHash(data))
	}
}
