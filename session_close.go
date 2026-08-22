package chunkcdc

func (s *Session) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return nil
	}
	// closed is the authoritative signal; do not hollow out the index
	// (byFP/entries/win) to fake a closed state — that leaves mutating
	// paths that miss the check writing to nil maps and panicking.
	s.closed = true
	if s.audit != nil {
		err := s.audit.Close()
		s.audit = nil
		return err
	}
	return nil
}

func (s *Session) CloseFlushCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return 0
	}
	n := s.FlushCount()
	s.ClearEntries()
	s.closed = true
	if s.audit != nil {
		_ = s.audit.Close()
		s.audit = nil
	}
	return n
}
