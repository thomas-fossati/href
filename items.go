package href

import (
	"fmt"
	"math"
	"strings"
)

type Items struct {
	values []string
}

func (o Items) Count() uint64 {
	return uint64(len(o.values))
}

func (o Items) String(sep string) string {
	if !o.IsSet() {
		return ""
	}
	return strings.Join(o.values, sep)
}

func (o Items) IsSet() bool {
	return len(o.values) > 0
}

func (o Items) Get() interface{} {
	if !o.IsSet() {
		return nil
	}
	return o.values
}

func (o Items) GetValues() []string {
	if !o.IsSet() {
		panic("there are no values to get")
	}
	return o.values
}

func (o *Items) Reset() {
	o.values = []string{}
}

func (o *Items) Set(v interface{}) error {
	switch t := v.(type) {
	case []interface{}:
		if err := o.SetValues(t); err != nil {
			return err
		}
		return nil
	case nil:
		return nil
	default:
		return fmt.Errorf("unknown type: %T", t)
	}
}

func (o *Items) SetValues(v []interface{}) error {
	for _, e := range v {
		if s, ok := e.(string); ok {
			o.values = append(o.values, s)
			continue
		}
		return fmt.Errorf("unknow type for item: %T", e)
	}
	return nil
}

func (o *Items) Append(v []string) {
	o.values = append(o.values, v...)
}

func (o *Items) TrimN(n uint64) {
	if o.Count() < n || n > math.MaxInt {
		return
	}

	o.values = o.values[:len(o.values)-int(n)]
}
