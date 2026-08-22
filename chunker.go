package chunkcdc
import ("hash/adler32")
type Chunk struct { Offset int; Length int; Hash uint32 }
type Chunker struct {
        MinSize, MaxSize, AvgSize int
}
func (c Chunker) withDefaults() Chunker {
        if c.MinSize <= 0 { c.MinSize = 64 }
        if c.MaxSize <= 0 { c.MaxSize = 1024 }
        if c.AvgSize <= 0 { c.AvgSize = 256 }
        return c
}
func (c Chunker) Split(data []byte) []Chunk {
        c = c.withDefaults()
        if len(data) == 0 { return nil }
        mask := uint32(c.AvgSize - 1)
        var out []Chunk
        start := 0
        for i := 0; i < len(data); i++ {
            size := i - start + 1
            if size < c.MinSize { continue }
            h := adler32.Checksum(data[start : i+1])
            if (h&mask) == 0 || size >= c.MaxSize {
                    out = append(out, Chunk{Offset: start, Length: size, Hash: h})
                    start = i + 1
            }
        }
        if start < len(data) {
                h := adler32.Checksum(data[start:])
                out = append(out, Chunk{Offset: start, Length: len(data) - start, Hash: h})
        }
        return out
}
