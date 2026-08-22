package chunkcdc

import "fmt"

// CollisionFail wraps a fingerprint conflict as ErrHashCollision.
func CollisionFail(fp string) error {
	return fmt.Errorf("%w: %s", ErrHashCollision, fp)
}
