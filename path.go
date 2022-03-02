package href

// path = [*text]
type Path struct {
	Items
}

func (o *Path) Set(v interface{}) error {
	return o.Items.Set(v)
}

func (o Path) IsSet() bool {
	return o.Items.IsSet()
}

func (o Path) Get() interface{} {
	return o.Items.Get()
}

func (o Path) GetSegments() []string {
	return o.Items.GetValues()
}

func (o *Path) Reset() {
	o.Items.Reset()
}

func (o Path) NumSegments() uint64 {
	return o.Items.Count()
}

func (o *Path) TrimN(n uint64) {
	o.Items.TrimN(n)
}

func (o *Path) Append(v []string) {
	o.Items.Append(v)
}
