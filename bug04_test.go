package chunkcdc_test

import (
	"errors"
	"testing"

	chunkcdc "github.com/LYH2263/go-chunkcdc"
)

func TestBug04_NilHasherNoDirtyChunk(t *testing.T) {
	s, err := chunkcdc.OpenSession(chunkcdc.SessionOptions{})
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	var panicked bool
	func() {
		defer func() {
			if recover() != nil {
				panicked = true
			}
		}()
		err := s.IngestSplit([]byte("payload"))
		if panicked {
			return
		}
		if err == nil || !errors.Is(err, chunkcdc.ErrNoHasher) {
			t.Fatalf("%v", err)
		}
	}()
	if panicked {
		t.Fatal("panic")
	}
	if s.VisibleCount() != 0 {
		t.Fatalf("dirty chunks %d", s.VisibleCount())
	}
}
