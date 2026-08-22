package chunkcdc_test

import (
	"testing"

	chunkcdc "github.com/LYH2263/go-chunkcdc"
)

func TestBug10_CloseFlushesIndexFirst(t *testing.T) {
	s, err := chunkcdc.OpenSession(chunkcdc.SessionOptions{Hasher: chunkcdc.AdlerHasher{}})
	if err != nil {
		t.Fatal(err)
	}
	data := []byte("flush-me")
	if err := s.CommitFP(chunkcdc.ChunkInfo{
		FP:     chunkcdc.Fingerprint(data),
		Hash:   chunkcdc.WeakHash(data),
		Length: len(data),
		Data:   data,
	}); err != nil {
		t.Fatal(err)
	}
	if s.CloseFlushCount() == 0 {
		t.Fatal("want flush")
	}
}
