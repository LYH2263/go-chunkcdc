package chunkcdc

import "github.com/LYH2263/go-chunkcdc/internal/clone"

// ListChunks returns an isolated snapshot of the current chunk index.
//
// The returned slice and each ChunkInfo.Data are deep copies, so callers may
// freely mutate preview data (e.g. truncating Data for display selection)
// without dirtying the internal entries that dedup (CommitFP/LookupFP) and
// index persistence rely on.
func (s *Session) ListChunks() []ChunkInfo {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]ChunkInfo, len(s.entries))
	for i, e := range s.entries {
		out[i] = e
		out[i].Data = clone.Bytes(e.Data)
	}
	return out
}
