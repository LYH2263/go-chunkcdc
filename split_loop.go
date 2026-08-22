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
	mask := uint32(c.AvgSize - 1)
	var out []Chunk
	start := 0
	for i := 0; i < len(data); i++ {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		size := i - start + 1
		if size < c.MinSize {
			continue
		}
		h := adler32.Checksum(data[start : i+1])
		if (h&mask) == 0 || size >= c.MaxSize {
			out = append(out, Chunk{Offset: start, Length: size, Hash: h})
			start = i + 1
		}
	}
	if start < len(data) {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		h := adler32.Checksum(data[start:])
		out = append(out, Chunk{Offset: start, Length: len(data) - start, Hash: h})
	}
	return out, nil
}
