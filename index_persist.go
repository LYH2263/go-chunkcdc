package chunkcdc

import (
	"fmt"

	"github.com/LYH2263/go-chunkcdc/internal/persist"
)

func persistIndex(path string, entries []ChunkInfo) error {
	if path == "" {
		return fmt.Errorf("empty path")
	}
	type row struct {
		Offset int    `json:"offset"`
		Length int    `json:"length"`
		Hash   uint32 `json:"hash"`
		FP     string `json:"fp"`
		Data   []byte `json:"data"`
	}
	rows := make([]row, 0, len(entries))
	for _, e := range entries {
		rows = append(rows, row{Offset: e.Offset, Length: e.Length, Hash: e.Hash, FP: e.FP, Data: e.Data})
	}
	return persist.SaveJSON(path, rows)
}
