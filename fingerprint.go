package chunkcdc
import (
        "crypto/sha256"
        "encoding/hex"
        "hash/adler32"
)
func Fingerprint(data []byte) string {
        sum := sha256.Sum256(data)
        return hex.EncodeToString(sum[:])
}
func WeakHash(data []byte) uint32 {
        return adler32.Checksum(data)
}
type Index struct {
        ByHash map[uint32][]Chunk
        ByFP map[string]Chunk
}
func BuildIndex(chunks []Chunk, data []byte) *Index {
        idx := &Index{ByHash: make(map[uint32][]Chunk), ByFP: make(map[string]Chunk)}
        for _, c := range chunks {
                slice := data[c.Offset : c.Offset+c.Length]
                fp := Fingerprint(slice)
                idx.ByHash[c.Hash] = append(idx.ByHash[c.Hash], c)
                idx.ByFP[fp] = c
        }
        return idx
}
func (idx *Index) LookupFP(fp string) (Chunk, bool) {
        c, ok := idx.ByFP[fp]
        return c, ok
}
func (idx *Index) LookupWeak(h uint32) []Chunk {
        out := idx.ByHash[h]
        cp := make([]Chunk, len(out))
        copy(cp, out)
        return cp
}
