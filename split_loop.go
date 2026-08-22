package chunkcdc

import (
	"context"
	"hash/adler32"
)

func SplitContext(ctx context.Context, c Chunker, data []byte) ([]Chunk, error) {
	c = c.withDefaults()
	if len(data) == 0 {
		return nil, nil
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	mask := uint32(c.AvgSize - 1)
	var out []Chunk
	start := 0
	for i := 0; i < len(data); i++ {
		size := i - start + 1
		if size < c.MinSize {
			continue
		}
		h := adler32.Checksum(data[start : i+1])
		if (h&mask) == 0 || size >= c.MaxSize {
			out = append(out, Chunk{Offset: start, Length: size, Hash: h})
			start = i + 1
			// Honor cancellation at chunk boundaries (guaranteed at least every
			// MaxSize bytes) so a cancelled split returns promptly instead of
			// grinding through the whole file. Return nil so partial chunks do
			// not keep accumulating downstream.
			if err := ctx.Err(); err != nil {
				return nil, err
			}
		}
	}
	if start < len(data) {
		h := adler32.Checksum(data[start:])
		out = append(out, Chunk{Offset: start, Length: len(data) - start, Hash: h})
	}
	return out, nil
}
