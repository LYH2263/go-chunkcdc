package chunkcdc_test

import (
	"testing"

	chunkcdc "github.com/LYH2263/go-chunkcdc"
)

func TestBug02_ListChunksIsolated(t *testing.T) {
	s, err := chunkcdc.OpenSession(chunkcdc.SessionOptions{Hasher: chunkcdc.AdlerHasher{}})
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	data := []byte("chunk-payload-aaaa")
	info := chunkcdc.ChunkInfo{Offset: 0, Length: len(data), Hash: chunkcdc.WeakHash(data), FP: chunkcdc.Fingerprint(data), Data: data}
	if err := s.CommitFP(info); err != nil {
		t.Fatal(err)
	}
	list := s.ListChunks()
	if len(list) != 1 {
		t.Fatalf("len=%d", len(list))
	}
	want := data[0]
	list[0].Data[0] ^= 0xff
	got, ok := s.LookupFP(info.FP)
	if !ok {
		t.Fatal("missing")
	}
	if got.Data[0] != want {
		t.Fatalf("polluted %v vs %v", got.Data[0], want)
	}
}
