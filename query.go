package href

// query = [*text]
type Query struct {
	Items
}

func (o Query) String() string {
	return o.Items.String("&")
}

func (o Query) IsSet() bool {
	return o.Items.IsSet()
}

func (o Query) Get() interface{} {
	return o.Items.Get()
}

func (o *Query) Reset() {
	o.Items.Reset()
}

func (o *Query) Set(v interface{}) error {
	return o.Items.Set(v)
}
