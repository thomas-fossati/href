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

func (o *PC) Step() error {
	o.curr += 1

	if o.AtEnd() {
		return ErrEndOfArray
	}

	return nil
}

func (o PC) Curr() (interface{}, error) {
	if o.AtEnd() {
		return nil, ErrEndOfArray
	}

	return o.comp[o.curr], nil
}

func (o *PC) Next() (interface{}, error) {
	if err := o.Step(); err != nil {
		return nil, err
	}

	return o.Curr()
}

func (o PC) Peek() (interface{}, error) {
	if o.AtEnd() {
		return nil, ErrEndOfArray
	}

	return o.comp[o.curr], nil
}
