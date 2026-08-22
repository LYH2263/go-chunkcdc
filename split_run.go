package chunkcdc

import "github.com/LYH2263/go-chunkcdc/internal/clone"

func (s *Session) IngestSplit(data []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return ErrClosed
	}
	if s.hasher == nil {
		return ErrNoHasher
	}
	cp := clone.Bytes(data)
	h := s.hasher.Sum32(cp)
	info := ChunkInfo{
		Offset: 0,
		Length: len(cp),
		Hash:   h,
		FP:     Fingerprint(cp),
		Data:   cp,
	}
	s.entries = append(s.entries, info)
	s.byFP[info.FP] = len(s.entries) - 1
	s.win = NewWindow(data, s.winSize)
	return nil
}
