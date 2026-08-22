package chunkcdc

func (s *Session) ListChunks() []ChunkInfo {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.entries
}
