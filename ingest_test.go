package chunkcdc

import (
	"path/filepath"
	"testing"
)

// TestIngestPersist_FailureLeavesVisibleSetUntouched reproduces bug6: when the
// index persist path is unwritable, IngestPersist returns an error AND the
// in-memory visible set (entries / byFP) must remain at its pre-ingest state.
//
// Previously the chunk was appended to memory before persisting, so a failed
// persist left ListChunks / VisibleCount reporting the just-ingested chunk
// while the on-disk index lagged behind -- diverging memory vs disk, breaking
// backup verification. The on-disk index is the source of truth: persist
// first, commit to memory only on success.
func TestIngestPersist_FailureLeavesVisibleSetUntouched(t *testing.T) {
	// A path whose parent directory cannot be created: the drive root has no
	// write permission for an unwritable subdir, so MkdirAll / WriteFile fails.
	badPath := filepath.Join(string(filepath.Separator), "cannot", "exist", "index.json")

	s, err := OpenSession(SessionOptions{PersistPath: badPath})
	if err != nil {
		t.Fatalf("OpenSession: %v", err)
	}
	defer s.Close()

	before := s.VisibleCount()

	if err := s.IngestPersist([]byte("hello")); err == nil {
		t.Fatalf("IngestPersist: expected error on unwritable persist path, got nil")
	}

	if got := s.VisibleCount(); got != before {
		t.Fatalf("VisibleCount after failed persist: got %d, want %d (visible set must not advance on persist failure)", got, before)
	}
	if got := s.FlushCount(); got != before {
		t.Fatalf("FlushCount after failed persist: got %d, want %d", got, before)
	}

	if chunks := s.ListChunks(); len(chunks) != before {
		t.Fatalf("ListChunks after failed persist: got %d chunks, want %d", len(chunks), before)
	}

	if _, ok := s.LookupFP(Fingerprint([]byte("hello"))); ok {
		t.Fatalf("LookupFP: failed ingest chunk is visible via byFP (must be rolled back)")
	}
}

// TestIngestPersist_EmptyPathReturnsErrorWithoutSideEffect covers the empty
// persist path branch: persistIndex returns "empty path" and, like any other
// persist failure, must not advance the visible set.
func TestIngestPersist_EmptyPathReturnsErrorWithoutSideEffect(t *testing.T) {
	s, err := OpenSession(SessionOptions{})
	if err != nil {
		t.Fatalf("OpenSession: %v", err)
	}
	defer s.Close()

	before := s.VisibleCount()

	if err := s.IngestPersist([]byte("hello")); err == nil {
		t.Fatalf("IngestPersist: expected error on empty persist path, got nil")
	}

	if got := s.VisibleCount(); got != before {
		t.Fatalf("VisibleCount after failed persist (empty path): got %d, want %d", got, before)
	}
}

// TestIngestPersist_SuccessCommitsToMemory ensures the fix did not break the
// happy path: a successful persist must make the chunk visible in memory and
// on disk, and a second ingest of the same bytes is a no-op via byFP.
func TestIngestPersist_SuccessCommitsToMemory(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "index.json")

	s, err := OpenSession(SessionOptions{PersistPath: path})
	if err != nil {
		t.Fatalf("OpenSession: %v", err)
	}
	defer s.Close()

	if err := s.IngestPersist([]byte("hello")); err != nil {
		t.Fatalf("IngestPersist: %v", err)
	}

	if got := s.VisibleCount(); got != 1 {
		t.Fatalf("VisibleCount after success: got %d, want 1", got)
	}

	chunks := s.ListChunks()
	if len(chunks) != 1 || string(chunks[0].Data) != "hello" {
		t.Fatalf("ListChunks: got %+v, want one chunk with data \"hello\"", chunks)
	}

	if _, ok := s.LookupFP(Fingerprint([]byte("hello"))); !ok {
		t.Fatalf("LookupFP: successful ingest chunk not visible via byFP")
	}
}

// TestIngestPersist_FailureDoesNotCorruptSubsequentSuccess confirms that after
// a failed ingest, a subsequent successful ingest on a valid path sees exactly
// one chunk -- no leftover from the rolled-back attempt.
func TestIngestPersist_FailureDoesNotCorruptSubsequentSuccess(t *testing.T) {
	badPath := filepath.Join(string(filepath.Separator), "cannot", "exist", "index.json")
	dir := t.TempDir()
	goodPath := filepath.Join(dir, "index.json")

	// Start on the bad path so the first ingest fails and is rolled back.
	s, err := OpenSession(SessionOptions{PersistPath: badPath})
	if err != nil {
		t.Fatalf("OpenSession: %v", err)
	}
	defer s.Close()

	if err := s.IngestPersist([]byte("dropped")); err == nil {
		t.Fatalf("first IngestPersist: expected error, got nil")
	}
	if got := s.VisibleCount(); got != 0 {
		t.Fatalf("after failed ingest: VisibleCount got %d, want 0", got)
	}

	// Switch to a writable path. The rolled-back chunk must not resurface.
	s.persistPath = goodPath

	if err := s.IngestPersist([]byte("kept")); err != nil {
		t.Fatalf("second IngestPersist: %v", err)
	}
	if got := s.VisibleCount(); got != 1 {
		t.Fatalf("after successful ingest: VisibleCount got %d, want 1", got)
	}

	chunks := s.ListChunks()
	if len(chunks) != 1 || string(chunks[0].Data) != "kept" {
		t.Fatalf("ListChunks after recovery: got %+v, want single chunk \"kept\"", chunks)
	}
}
