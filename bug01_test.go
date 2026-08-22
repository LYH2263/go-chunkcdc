package chunkcdc_test

import (
	"bytes"
	"testing"

	chunkcdc "github.com/LYH2263/go-chunkcdc"
)

func TestBug01_IngestIsolatesWindowBuf(t *testing.T) {
	s, err := chunkcdc.OpenSession(chunkcdc.SessionOptions{Hasher: chunkcdc.AdlerHasher{}})
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	buf := []byte("0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	wantHash := chunkcdc.WeakHash(buf[:64])
	if err := s.Ingest(buf); err != nil {
		t.Fatal(err)
	}
	buf[0] = 'Z'
	got := s.PeekWindow()
	if bytes.Equal(got, buf[:len(got)]) {
		t.Fatalf("ingest leaked %q", got)
	}
	if s.WindowHash() != wantHash {
		t.Fatalf("window hash polluted %d vs %d", s.WindowHash(), wantHash)
	}
}
