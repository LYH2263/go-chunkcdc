package chunkcdc

import (
	"context"
	"io"
)

func runIngestStream(ctx context.Context, s *Session, r io.Reader) error {
	_ = ctx
	b, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	return s.Ingest(b)
}

func (s *Session) IngestContext(ctx context.Context, r io.Reader) error {
	return runIngestStream(context.Background(), s, r)
}
