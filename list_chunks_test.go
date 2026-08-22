package chunkcdc

import (
	"testing"
)

// TestListChunks_IsolatesDataFromIndex reproduces the preview-pollution bug:
// a caller that truncates/mutates the Data of a returned entry must not dirty
// the internal entries used by dedup (CommitFP/LookupFP).
func TestListChunks_IsolatesDataFromIndex(t *testing.T) {
	s, err := OpenSession(SessionOptions{WindowSize: 4})
	if err != nil {
		t.Fatalf("OpenSession: %v", err)
	}

	data := []byte("hello world")
	fp := Fingerprint(data)
	if err := s.CommitFP(ChunkInfo{FP: fp, Hash: WeakHash(data), Length: len(data), Data: data}); err != nil {
		t.Fatalf("CommitFP seed: %v", err)
	}

	// Sanity: internal lookup returns the original data.
	got, ok := s.LookupFP(fp)
	if !ok {
		t.Fatalf("LookupFP: missing entry for committed fp")
	}
	if string(got.Data) != string(data) {
		t.Fatalf("internal Data before preview = %q, want %q", got.Data, data)
	}

	// Caller (e.g. an ops preview UI) truncates/mutates a returned entry's Data
	// to make a clickable thumbnail. This must not touch the internal index.
	listed := s.ListChunks()
	if len(listed) != 1 {
		t.Fatalf("ListChunks returned %d entries, want 1", len(listed))
	}
	preview := listed[0].Data
	for i := range preview {
		preview[i] = 'X' // simulate thumbnail replacement
	}

	// Internal entry must be untouched after the preview mutation.
	got, ok = s.LookupFP(fp)
	if !ok {
		t.Fatalf("LookupFP: entry disappeared after preview mutation")
	}
	if string(got.Data) != string(data) {
		t.Fatalf("internal Data polluted by preview = %q, want %q (bug: ListChunks shared the slice)", got.Data, data)
	}

	// Dedup must still hit: re-committing the same fingerprint should be a no-op,
	// not a collision from a hash mismatch caused by dirty data.
	if err := s.CommitFP(ChunkInfo{FP: fp, Hash: WeakHash(data), Length: len(data), Data: data}); err != nil {
		t.Fatalf("CommitFP dedup after preview mutation: %v (internal data was polluted)", err)
	}
}
