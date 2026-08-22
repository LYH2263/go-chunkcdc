package chunkcdc

import (
	"context"
	"testing"
)

func TestSplitContext_HonorsCancellation(t *testing.T) {
	// A large input forces many chunk boundaries so a non-cancel-aware loop
	// would grind through the whole slice before returning. With MaxSize=64 the
	// boundary check fires at least every 64 bytes, so the context is consulted
	// frequently and the loop must stop promptly after cancel.
	data := make([]byte, 1<<20) // 1 MiB
	for i := range data {
		data[i] = byte(i)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	chunks, err := SplitContext(ctx, Chunker{MinSize: 64, MaxSize: 64, AvgSize: 64}, data)
	if err != context.Canceled {
		t.Fatalf("expected context.Canceled, got %v (chunks=%d)", err, len(chunks))
	}
	if chunks != nil {
		t.Fatalf("expected nil chunks on cancel, got %d", len(chunks))
	}
}

func TestSplitContext_StillWorksWhenNotCancelled(t *testing.T) {
	data := make([]byte, 4096)
	for i := range data {
		data[i] = byte(i)
	}

	chunks, err := SplitContext(context.Background(), Chunker{}, data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(chunks) == 0 {
		t.Fatalf("expected non-empty chunks, got none")
	}

	// Chunks must cover the full input without gaps or overlaps.
	var covered int
	for _, c := range chunks {
		covered += c.Length
	}
	if covered != len(data) {
		t.Fatalf("chunks cover %d bytes, input is %d", covered, len(data))
	}
}
