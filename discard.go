package href

import (
	"errors"
	"fmt"
)

/*
discard     = true / 0..127
*/

type Discard struct {
	val interface{}
}

func isDiscard(e interface{}) bool {
	switch t := e.(type) {
	case bool: // true
		return t
	case uint64: // 0..127
		return t <= 127
	}
	return false
}

func (o *Discard) Set(v interface{}) error {
	switch t := v.(type) {
	case bool: // true
		if !t {
			return errors.New("discard cannot be false")
		}
		o.val = t
	case uint64: // 0..127
		if t > 127 {
			return fmt.Errorf("discard must be in range 0..127, got %d", t)
		}
		o.val = t
	default:
		return fmt.Errorf("unknown discard type: %T", t)
	}

	return nil
}
