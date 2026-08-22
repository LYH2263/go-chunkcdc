package chunkcdc_test

import (
	"errors"
	"testing"

	chunkcdc "github.com/LYH2263/go-chunkcdc"
)

func TestBug03_IngestAfterClose(t *testing.T) {
	s, err := chunkcdc.OpenSession(chunkcdc.SessionOptions{Hasher: chunkcdc.AdlerHasher{}})
	if err != nil {
		t.Fatal(err)
	}
	_ = s.Ingest([]byte("seed-data-0123456789abcdef0123456789abcdef0123456789abcdef01234567"))
	s.Close()
	var panicked bool
	var got error
	func() {
		defer func() {
			if recover() != nil {
				panicked = true
			}
		}()
		got = s.Ingest([]byte("after-close-0123456789abcdef0123456789abcdef0123456789abcdef012345"))
	}()
	if panicked {
		t.Fatal("panic")
	}
	if !errors.Is(got, chunkcdc.ErrClosed) {
		t.Fatalf("%v", got)
	}
}
