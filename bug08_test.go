package chunkcdc_test

import (
	"context"
	"testing"

	chunkcdc "github.com/LYH2263/go-chunkcdc"
)

func TestBug08_SplitLoopHonorsCancel(t *testing.T) {
	data := make([]byte, 4096)
	for i := range data {
		data[i] = byte(i)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	c := chunkcdc.Chunker{MinSize: 64, MaxSize: 256, AvgSize: 128}
	if _, err := chunkcdc.SplitContext(ctx, c, data); err == nil {
		t.Fatal("want cancel")
	}
}
