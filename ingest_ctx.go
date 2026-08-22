package chunkcdc

import (
	"context"
	"io"
)

func runIngestStream(ctx context.Context, s *Session, r io.Reader) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	b, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	return s.Ingest(b)
}

func (s *Session) IngestContext(ctx context.Context, r io.Reader) error {
	return runIngestStream(ctx, s, r)
}
