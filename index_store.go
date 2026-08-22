package chunkcdc

func (s *Session) LookupFP(fp string) (ChunkInfo, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	i, ok := s.byFP[fp]
	if !ok {
		return ChunkInfo{}, false
	}
	return s.entries[i], true
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
	cp := info
	if cp.FP == "" {
		cp.FP = Fingerprint(cp.Data)
	}
	if i, ok := s.byFP[cp.FP]; ok {
		old := s.entries[i]
		if old.Hash != cp.Hash || old.Length != cp.Length || string(old.Data) != string(cp.Data) {
			return CollisionFail(cp.FP)
		}
		return nil
	}
	s.entries = append(s.entries, cp)
	s.byFP[cp.FP] = len(s.entries) - 1
	return nil
}
