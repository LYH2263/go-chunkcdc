package chunkcdc_test

import (
	"os"
	"path/filepath"
	"testing"

	chunkcdc "github.com/LYH2263/go-chunkcdc"
)

func TestBug09_RotateClosesOldAudit(t *testing.T) {
	dir := t.TempDir()
	p1 := filepath.Join(dir, "a.log")
	p2 := filepath.Join(dir, "b.log")
	s, err := chunkcdc.OpenSession(chunkcdc.SessionOptions{
		Hasher:    chunkcdc.AdlerHasher{},
		AuditPath: p1,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	if err := s.Ingest([]byte("0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")); err != nil {
		t.Fatal(err)
	}
	if err := s.RotateAudit(p2); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(p1); err != nil {
		t.Fatalf("locked %v", err)
	}
}
