package href

import (
	"errors"

	"github.com/fxamacker/cbor"
)

var (
	ErrEndOfArray = errors.New("end of array")
)

type PC struct {
	comp []interface{}
	curr int
}

func NewPC(cri []byte) (*PC, error) {
	var v []interface{}

	if err := cbor.Unmarshal(cri, &v); err != nil {
		return nil, err
	}

	return &PC{comp: v, curr: 0}, nil
}

func (o PC) Empty() bool {
	return len(o.comp) == 0
}

func (o PC) AtEnd() bool {
	return o.curr == len(o.comp)
}

func (o *PC) Step() bool {
	o.curr += 1

	if o.AtEnd() {
		return true
	}

	return false
}

func (o PC) Curr() (interface{}, bool) {
	if o.AtEnd() {
		return nil, true
	}

	return o.comp[o.curr], false
}

func (o *PC) Next() (interface{}, bool) {
	if o.Step() {
		return nil, true
	}

	return o.Curr()
}

func (o PC) Peek() (interface{}, bool) {
	if o.AtEnd() {
		return nil, true
	}

	return o.comp[o.curr], false
}
