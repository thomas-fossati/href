package href

import "fmt"

type Fragment struct {
	val string
}

func (o Fragment) String() string {
	if o.val == "" {
		return ""
	}
	return "#" + o.val
}

func (o *Fragment) Set(v interface{}) error {
	if s, ok := v.(string); ok {
		o.val = s
		return nil
	}

	return fmt.Errorf("unknown type for fragment: %T", v)
}
