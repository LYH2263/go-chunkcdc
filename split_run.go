package chunkcdc

import "github.com/LYH2263/go-chunkcdc/internal/clone"

func (s *Session) IngestSplit(data []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return ErrClosed
	}
	cp := clone.Bytes(data)
	// dirty: register before hasher use
	info := ChunkInfo{
		Offset: 0,
		Length: len(cp),
		Hash:   0,
		FP:     Fingerprint(cp),
		Data:   cp,
	}
	s.entries = append(s.entries, info)
	s.byFP[info.FP] = len(s.entries) - 1
	h := s.hasher.Sum32(cp)
	if s.hasher == nil {
		return ErrNoHasher
	}
	s.entries[len(s.entries)-1].Hash = h
	s.win = NewWindow(cp, s.winSize)
	return nil
}
