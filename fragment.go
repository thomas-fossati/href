package href

import "fmt"

type Fragment struct {
	val *string
}

func (o Fragment) String() string {
	return o.Get()
}

func (o Fragment) IsSet() bool {
	return o.val != nil
}

func (o Fragment) Get() string {
	if o.IsSet() {
		return *o.val
	}
	return ""
}

func (o *Fragment) Reset() {
	o.val = nil
}

func (o *Fragment) Set(v interface{}) error {
	if s, ok := v.(string); ok {
		o.val = &s
		return nil
	}

	return fmt.Errorf("unknown type for fragment: %T", v)
}
