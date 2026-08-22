package chunkcdc

import "github.com/LYH2263/go-chunkcdc/internal/clone"

func (s *Session) ListChunks() []ChunkInfo {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]ChunkInfo, len(s.entries))
	for i, c := range s.entries {
		out[i] = ChunkInfo{
			Offset: c.Offset,
			Length: c.Length,
			Hash:   c.Hash,
			FP:     c.FP,
			Data:   clone.Bytes(c.Data),
		}
	}
	return out
}
