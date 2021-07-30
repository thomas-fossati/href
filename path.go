package href

import (
	"fmt"
	"strings"
)

// path = [*text]
type Path struct {
	Segments []string
	IsNull   bool
}

func (o Path) String() string {
	if o.IsNull || len(o.Segments) == 0 {
		return ""
	}
	return "/" + strings.Join(o.Segments, "/")
}

func (o *Path) Set(v interface{}) error {
	switch t := v.(type) {
	case []interface{}:
		return o.SetPath(t)
	case nil:
		o.IsNull = true
		return nil
	default:
		return fmt.Errorf("unknown path type: %T", t)
	}
}

func (o *Path) SetPath(v []interface{}) error {
	for _, e := range v {
		if s, ok := e.(string); ok {
			o.Segments = append(o.Segments, s)
			continue
		}
		return fmt.Errorf("unknow type for path segment: %T", e)
	}
	return nil
}
