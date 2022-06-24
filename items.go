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

func (o Items) StringEscaped(sep string, e escaper) string {
	if !o.IsSet() {
		return ""
	}

	return strings.Join(o.GetEscapedStrings(e), sep)
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

type escaper func(string) string
type unescaper func(string) (string, error)

func (o Items) GetUnescaped(u unescaper) interface{} {
	return o.GetUnescapedStrings(u)
}

func (o Items) GetUnescapedStrings(u unescaper) []string {
	if o.Get() == nil {
		return nil
	}

	unescaped := make([]string, len(o.values))
	copy(unescaped, o.values)

	for i, qe := range unescaped {
		unescaped[i], _ = u(qe)
	}

	return unescaped
}

func (o Items) GetEscapedStrings(e escaper) []string {
	if o.Get() == nil {
		return nil
	}

	escaped := make([]string, len(o.values))
	copy(escaped, o.values)

	for i, qe := range escaped {
		escaped[i] = e(qe)
	}

	return escaped
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
