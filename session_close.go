package chunkcdc

func (s *Session) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return nil
	}
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
	s.ClearEntries()
	n := s.FlushCount()
	s.closed = true
	if s.audit != nil {
		_ = s.audit.Close()
		s.audit = nil
	}
	return n
}
