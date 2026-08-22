package chunkcdc

import (
	"sync"

	"github.com/LYH2263/go-chunkcdc/internal/audit"
)

type ChunkInfo struct {
	Offset int
	Length int
	Hash   uint32
	FP     string
	Data   []byte
}

type SessionOptions struct {
	AuditPath   string
	PersistPath string
	WindowSize  int
	Hasher      Hasher
}

type Session struct {
	mu          sync.Mutex
	closed      bool
	win         *Window
	winSize     int
	entries     []ChunkInfo
	byFP        map[string]int
	hasher      Hasher
	audit       *audit.Logger
	persistPath string
}

func OpenSession(opts SessionOptions) (*Session, error) {
	if opts.WindowSize <= 0 {
		opts.WindowSize = 64
	}
	s := &Session{
		winSize:     opts.WindowSize,
		byFP:        make(map[string]int),
		hasher:      opts.Hasher,
		persistPath: opts.PersistPath,
	}
	if opts.AuditPath != "" {
		al, err := audit.Open(opts.AuditPath)
		if err != nil {
			return nil, err
		}
		s.audit = al
	}
	return s, nil
}

func (s *Session) SetHasher(h Hasher) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.hasher = h
}

func (s *Session) VisibleCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.entries)
}

func (s *Session) PeekWindow() []byte {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.win == nil {
		return nil
	}
	return s.win.Bytes()
}

func (s *Session) WindowHash() uint32 {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.win == nil {
		return 0
	}
	return s.win.Hash()
}

func (s *Session) RotateAudit(newPath string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.audit == nil {
		return ErrInvalid
	}
	return s.audit.Rotate(newPath)
}

func (s *Session) AuditPath() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.audit == nil {
		return ""
	}
	return s.audit.Path()
}
