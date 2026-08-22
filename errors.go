package chunkcdc

import "errors"

var (
	ErrClosed        = errors.New("chunkcdc: closed")
	ErrInvalid       = errors.New("chunkcdc: invalid")
	ErrNotFound      = errors.New("chunkcdc: not found")
	ErrConflict      = errors.New("chunkcdc: conflict")
	ErrHashCollision = errors.New("chunkcdc: hash collision")
	ErrNoHasher      = errors.New("chunkcdc: no hasher")
)
