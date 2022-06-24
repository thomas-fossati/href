package href

import "net/url"

// query = [*text]
type Query struct {
	Items
}

func (o Query) String() string {
	return o.Items.String("&")
}

func (o Query) StringEscaped() string {
	return o.Items.StringEscaped("&", url.QueryEscape)
}

func (o Query) IsSet() bool {
	return o.Items.IsSet()
}

func (o Query) Get() interface{} {
	return o.Items.Get()
}

func (o Query) GetUnescaped() interface{} {
	return o.Items.GetUnescaped(url.QueryUnescape)
}

func (o *Query) Reset() {
	o.Items.Reset()
}

func (o *Query) Set(v interface{}) error {
	return o.Items.Set(v)
}
