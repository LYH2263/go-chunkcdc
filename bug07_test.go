package chunkcdc_test

import (
	"bytes"
	"context"
	"testing"

	chunkcdc "github.com/LYH2263/go-chunkcdc"
)

func TestBug07_IngestContextHonorsCancel(t *testing.T) {
	s, err := chunkcdc.OpenSession(chunkcdc.SessionOptions{Hasher: chunkcdc.AdlerHasher{}})
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	payload := bytes.NewReader([]byte("0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"))
	if err := s.IngestContext(ctx, payload); err == nil {
		t.Fatal("want cancel")
	}
}
