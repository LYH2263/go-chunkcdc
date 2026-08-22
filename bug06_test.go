package chunkcdc_test

import (
	"os"
	"path/filepath"
	"testing"

	chunkcdc "github.com/LYH2263/go-chunkcdc"
)

func TestBug06_PersistFailRollsBack(t *testing.T) {
	dir := t.TempDir()
	badPath := filepath.Join(dir, "index.json")
	if err := os.Mkdir(badPath, 0o755); err != nil {
		t.Fatal(err)
	}
	s, err := chunkcdc.OpenSession(chunkcdc.SessionOptions{
		Hasher:      chunkcdc.AdlerHasher{},
		PersistPath: badPath,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	before := s.VisibleCount()
	if err := s.IngestPersist([]byte("block-bytes")); err == nil {
		t.Fatal("want persist fail")
	}
	if s.VisibleCount() != before {
		t.Fatalf("visible drifted %d -> %d", before, s.VisibleCount())
	}
}
