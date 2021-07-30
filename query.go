package href

import (
	"fmt"
	"strings"
)

// query = [*text]
type Query struct {
	Parts  []string
	IsNull bool
}

func (o Query) String() string {
	if o.IsNull || len(o.Parts) == 0 {
		return ""
	}
	return "?" + strings.Join(o.Parts, "&")
}

func (o *Query) Set(v interface{}) error {
	switch t := v.(type) {
	case []interface{}:
		return o.SetQuery(t)
	case nil:
		o.IsNull = true
		return nil
	default:
		return fmt.Errorf("unknown query type: %T", t)
	}
}

func (o *Query) SetQuery(v []interface{}) error {
	for _, e := range v {
		if s, ok := e.(string); ok {
			o.Parts = append(o.Parts, s)
			continue
		}
		return fmt.Errorf("unknow type for query part: %T", e)
	}
	return nil
}
