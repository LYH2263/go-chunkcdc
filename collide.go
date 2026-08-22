package chunkcdc

import "fmt"

func CollisionFail(fp string) error {
	return fmt.Errorf("collision: %s", fp)
}
