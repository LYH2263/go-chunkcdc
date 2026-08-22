package chunkcdc

import "github.com/LYH2263/go-chunkcdc/internal/clone"

func (s *Session) Ingest(buf []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return ErrClosed
	}
	if s.byFP == nil {
		return ErrClosed
	}
	cp := clone.Bytes(buf)
	s.win = NewWindow(cp, s.winSize)
	if s.audit != nil {
		_ = s.audit.Log("ingest", Fingerprint(cp))
	}
	return nil
}

func (s *Session) IngestPersist(buf []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return ErrClosed
	}
	cp := clone.Bytes(buf)
	info := ChunkInfo{
		Offset: 0,
		Length: len(cp),
		Hash:   WeakHash(cp),
		FP:     Fingerprint(cp),
		Data:   cp,
	}
	// dirty: visible before persist
	s.entries = append(s.entries, info)
	s.byFP[info.FP] = len(s.entries) - 1
	s.win = NewWindow(cp, s.winSize)
	if err := persistIndex(s.persistPath, s.entries); err != nil {
		return err
	}
	return nil
}
