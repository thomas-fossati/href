package href

import (
	"fmt"
	"regexp"
)

type Scheme struct {
	val interface{}
}

const (
	schemeREString = `[a-z][a-z0-9+.-]*`
)

var (
	schemeRE = regexp.MustCompile(schemeREString)

	schemeIDtoString = map[int64]string{
		-1: "coap",
		-2: "coaps",
		-3: "http",
		-4: "https",
		-5: "urn",
		-6: "did",
	}
)

func (o Scheme) IsSet() bool {
	return o.val != nil
}

func (o Scheme) Get() interface{} {
	return o.val
}

func (o *Scheme) Set(v interface{}) error {
	switch t := v.(type) {
	case string:
		// scheme-name
		if !schemeRE.MatchString(t) {
			return fmt.Errorf("scheme-name %s does not match scheme RE (%s)", t, schemeREString)
		}
		o.val = t
	case int64:
		// scheme-id
		if t >= 0 {
			return fmt.Errorf("scheme-id must be nint, got %d", t)
		}
		o.val = t
	case nil:
		// no scheme
		o.val = nil
	default:
		return fmt.Errorf("unknown scheme type: %T", t)
	}

	return nil
}

func SchemeIDtoString(schemeID int64) string {
	s, ok := schemeIDtoString[schemeID]
	if !ok {
		return fmt.Sprintf("scheme-id(%d)", schemeID)
	}
	return s
}

func (o Scheme) String() string {
	switch t := o.val.(type) {
	case string:
		return t
	case int64:
		return SchemeIDtoString(t)
	case nil:
		return ""
	}

	return ""
}

func isScheme(e interface{}) bool {
	switch t := e.(type) {
	case string: // scheme-name
		return true
	case int64: // scheme-id is nint
		return t <= 0
	case nil: // no scheme
		return true
	}
	return false
}
