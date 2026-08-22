package chunkcdc

func (s *Session) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	// wipe index without closed flag
	s.byFP = nil
	s.entries = nil
	s.win = nil
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
