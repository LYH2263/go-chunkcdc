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
	// Persist before making the chunk visible. The on-disk index is the
	// source of truth: if the write fails, the in-memory visible set
	// (entries / byFP) must stay at its pre-ingest state so ListChunks and
	// the persisted index can never diverge. Committing to memory only
	// after a successful persist guarantees no half-success visibility.
	if err := persistIndex(s.persistPath, append(s.entries, info)); err != nil {
		return err
	}
	s.entries = append(s.entries, info)
	s.byFP[info.FP] = len(s.entries) - 1
	s.win = NewWindow(cp, s.winSize)
	return nil
}
