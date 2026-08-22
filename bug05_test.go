package chunkcdc_test

import (
	"errors"
	"testing"

	chunkcdc "github.com/LYH2263/go-chunkcdc"
)

func TestBug05_HashCollisionWrapped(t *testing.T) {
	s, err := chunkcdc.OpenSession(chunkcdc.SessionOptions{Hasher: chunkcdc.AdlerHasher{}})
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	a := []byte("same-fp-body-1")
	fp := chunkcdc.Fingerprint(a)
	if err := s.CommitFP(chunkcdc.ChunkInfo{FP: fp, Hash: 1, Length: len(a), Data: a}); err != nil {
		t.Fatal(err)
	}
	b := []byte("same-fp-body-2")
	err = s.CommitFP(chunkcdc.ChunkInfo{FP: fp, Hash: 2, Length: len(b), Data: b})
	if err == nil || !errors.Is(err, chunkcdc.ErrHashCollision) {
		t.Fatalf("%v", err)
	}
	err2 := chunkcdc.CollisionFail(fp)
	if !errors.Is(err2, chunkcdc.ErrHashCollision) {
		t.Fatalf("collisionfail %v", err2)
	}
}
