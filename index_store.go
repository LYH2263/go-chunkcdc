package chunkcdc

import (
	"fmt"

	"github.com/LYH2263/go-chunkcdc/internal/clone"
)

func (s *Session) LookupFP(fp string) (ChunkInfo, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	i, ok := s.byFP[fp]
	if !ok {
		return ChunkInfo{}, false
	}
	c := s.entries[i]
	return ChunkInfo{
		Offset: c.Offset,
		Length: c.Length,
		Hash:   c.Hash,
		FP:     c.FP,
		Data:   clone.Bytes(c.Data),
	}, true
}

func (s *Session) FlushCount() int {
	return len(s.entries)
}

func (s *Session) ClearEntries() {
	s.entries = nil
	s.byFP = make(map[string]int)
}

func (s *Session) CommitFP(info ChunkInfo) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return ErrClosed
	}
	cp := ChunkInfo{
		Offset: info.Offset,
		Length: info.Length,
		Hash:   info.Hash,
		FP:     info.FP,
		Data:   clone.Bytes(info.Data),
	}
	if cp.FP == "" {
		cp.FP = Fingerprint(cp.Data)
	}
	if i, ok := s.byFP[cp.FP]; ok {
		old := s.entries[i]
		if old.Hash != cp.Hash || old.Length != cp.Length || string(old.Data) != string(cp.Data) {
			// Same fingerprint, different content: a hash collision. Wrap
			// ErrHashCollision so monitoring/callers can match the sentinel
			// with errors.Is instead of grepping message text.
			return fmt.Errorf("%w: fp conflict %s", ErrHashCollision, cp.FP)
		}
		return nil
	}
	s.entries = append(s.entries, cp)
	s.byFP[cp.FP] = len(s.entries) - 1
	return nil
}
