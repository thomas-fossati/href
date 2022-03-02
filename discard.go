package href

import (
	"errors"
	"fmt"
	"strings"
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

func (o Discard) IsSet() bool {
	return o.val != nil
}

func (o Discard) Get() interface{} {
	return o.val
}

func (o Discard) IsTrue() bool {
	if !o.IsSet() {
		return false
	}

	v := o.Get()

	if t, ok := v.(bool); ok {
		if !t {
			panic("boolean discard cannot be false")
		}
		return true
	}

	return false
}

// TODO(tho) handle failure condition "discard==0 && path item present"
func (o Discard) ComputePathPrefix() string {
	var pathPrefix string

	v := o.Get()

	switch t := v.(type) {
	case bool:
		// If the CRI reference contains a discard item of value true, the path
		// component is prefixed by a slash ("/") character.
		if !t {
			panic("boolean discard cannot be false")
		}
		pathPrefix = "/"
	case uint64:
		// If it contains a discard item of value 0 and the path item is
		// present, the conversion fails.
		//   XXX(tho) we are missing the path context to decide
		//   XXX(tho) whether to fail or not.
		// If it contains a positive discard item, the path component is
		// prefixed by as many "../" components as the discard value minus one
		// indicates.
		if t > 0 {
			var (
				i uint64
				a []string
			)
			for i = 0; i < t-1; i++ {
				a = append(a, "../")
			}
			pathPrefix = strings.Join(a, "")
		}
	default:
		panic(fmt.Sprintf("unknown type %T for discard", t))
	}

	return pathPrefix
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
