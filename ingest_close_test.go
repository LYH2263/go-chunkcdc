package chunkcdc

import (
	"context"
	"errors"
	"io"
	"testing"
)

// TestIngestAfterCloseStableErrClosed is the regression test for the backup-session
// tail-chunk incident. After Close(), every ingest path must return ErrClosed —
// not panic, not return nil, not return a downstream error the retry policy can't
// classify as "closed". Close() must set the closed flag rather than hollowing out
// the index structs to fake a closed state.
func TestIngestAfterCloseStableErrClosed(t *testing.T) {
	t.Run("Close sets closed flag without hollowing struct", func(t *testing.T) {
		s, err := OpenSession(SessionOptions{WindowSize: 4})
		if err != nil {
			t.Fatalf("OpenSession: %v", err)
		}
		// Prime the index via a path that actually populates entries/byFP so we
		// can detect Close() niling them out to fake a closed state. (Ingest
		// alone never writes entries, so it can't reveal hollowing.)
		if err := s.CommitFP(ChunkInfo{FP: "seedfp", Data: []byte("seed")}); err != nil {
			t.Fatalf("seed CommitFP: %v", err)
		}
		s.mu.Lock()
		entriesBefore := len(s.entries)
		s.mu.Unlock()
		if entriesBefore == 0 {
			t.Fatalf("seed did not populate entries")
		}

		if err := s.Close(); err != nil {
			t.Fatalf("Close: %v", err)
		}

		// The closed flag is the authoritative signal — the maps must be left
		// intact (not nilled out to impersonate a closed session).
		s.mu.Lock()
		closed := s.closed
		nilByFP := s.byFP == nil
		nilEntries := s.entries == nil
		s.mu.Unlock()

		if !closed {
			t.Fatalf("Close() did not set s.closed (got false); index was hollowed out instead")
		}
		if nilByFP || nilEntries {
			t.Fatalf("Close() nilled the index (byFP nil=%v, entries nil=%v); must set closed flag instead",
				nilByFP, nilEntries)
		}
	})

	ingestCases := []struct {
		name string
		fn   func(s *Session) error
	}{
		{"Ingest", func(s *Session) error { return s.Ingest([]byte("tail")) }},
		{"IngestPersist", func(s *Session) error { return s.IngestPersist([]byte("tail")) }},
		{"IngestSplit", func(s *Session) error {
			s.SetHasher(AdlerHasher{})
			return s.IngestSplit([]byte("tail"))
		}},
		{"CommitFP", func(s *Session) error {
			return s.CommitFP(ChunkInfo{FP: "fp", Data: []byte("tail")})
		}},
		{"IngestContext", func(s *Session) error {
			return s.IngestContext(context.Background(), nilContextReader([]byte("tail")))
		}},
	}
	for _, c := range ingestCases {
		c := c
		t.Run(c.name+" returns ErrClosed after Close (no panic)", func(t *testing.T) {
			s, err := OpenSession(SessionOptions{WindowSize: 4})
			if err != nil {
				t.Fatalf("OpenSession: %v", err)
			}
			if err := s.Ingest([]byte("seed")); err != nil {
				t.Fatalf("seed Ingest: %v", err)
			}
			if err := s.Close(); err != nil {
				t.Fatalf("Close: %v", err)
			}

			// The old code panicked here on a nil-map write or returned nil/a
			// non-ErrClosed error. It must now return ErrClosed.
			err = c.fn(s)
			if !errors.Is(err, ErrClosed) {
				t.Fatalf("after Close, %s = %v, want ErrClosed", c.name, err)
			}
		})
	}

	t.Run("Close is idempotent", func(t *testing.T) {
		s, err := OpenSession(SessionOptions{WindowSize: 4})
		if err != nil {
			t.Fatalf("OpenSession: %v", err)
		}
		if err := s.Close(); err != nil {
			t.Fatalf("first Close: %v", err)
		}
		if err := s.Close(); err != nil {
			t.Fatalf("second Close: %v", err)
		}
	})
}

// nilContextReader is a minimal io.Reader over an in-memory buffer, used so the
// ingest-context path can be exercised without pulling extra test deps.
func nilContextReader(b []byte) *bytesReader {
	return &bytesReader{b: b}
}

type bytesReader struct {
	b []byte
	i int
}

func (r *bytesReader) Read(p []byte) (int, error) {
	if r.i >= len(r.b) {
		return 0, errEOF
	}
	n := copy(p, r.b[r.i:])
	r.i += n
	return n, nil
}

var errEOF = io.EOF
