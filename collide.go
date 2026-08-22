package chunkcdc

import "fmt"

// CollisionFail reports a fingerprint collision (same FP, different content)
// and wraps ErrHashCollision so callers can identify it via errors.Is.
func CollisionFail(fp string) error {
	return fmt.Errorf("%w: %s", ErrHashCollision, fp)
}
