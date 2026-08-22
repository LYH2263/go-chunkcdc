package chunkcdc

// Hasher computes a weak rolling-style checksum for CDC boundaries.
type Hasher interface {
	Sum32(data []byte) uint32
}

// AdlerHasher wraps WeakHash.
type AdlerHasher struct{}

func (AdlerHasher) Sum32(data []byte) uint32 {
	return WeakHash(data)
}
