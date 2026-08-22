package chunkcdc

import (
	"errors"
	"testing"
)

// TestCommitFP_Collision wraps ErrHashCollision so monitoring can identify a
// collision via the sentinel rather than grepping message text.
func TestCommitFP_Collision(t *testing.T) {
	s, err := OpenSession(SessionOptions{WindowSize: 8})
	if err != nil {
		t.Fatalf("OpenSession: %v", err)
	}
	defer s.Close()

	first := ChunkInfo{Offset: 0, Length: 3, Hash: 1, FP: "fp-same", Data: []byte("aaa")}
	if err := s.CommitFP(first); err != nil {
		t.Fatalf("first CommitFP: %v", err)
	}

	// Same FP, different content -> collision.
	second := ChunkInfo{Offset: 0, Length: 3, Hash: 2, FP: "fp-same", Data: []byte("bbb")}
	err = s.CommitFP(second)
	if err == nil {
		t.Fatal("expected collision error, got nil")
	}
	if !errors.Is(err, ErrHashCollision) {
		t.Errorf("CommitFP collision: errors.Is(err, ErrHashCollision) = false; err=%q", err)
	}

	// Same FP AND same content -> dedup, not a collision.
	dup := ChunkInfo{Offset: 0, Length: 3, Hash: 1, FP: "fp-same", Data: []byte("aaa")}
	if err := s.CommitFP(dup); err != nil {
		t.Errorf("dedup CommitFP: want nil, got %v", err)
	}
}

// TestCollisionFail_wrapsErrHashCollision ensures CollisionFail is identifiable
// via the sentinel rather than relying on message text.
func TestCollisionFail_wrapsErrHashCollision(t *testing.T) {
	err := CollisionFail("deadbeef")
	if !errors.Is(err, ErrHashCollision) {
		t.Errorf("errors.Is(CollisionFail, ErrHashCollision) = false; err=%q", err)
	}
}
