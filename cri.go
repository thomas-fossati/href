package href

import (
	"fmt"
)

type CRI struct {
	Discard   Discard
	Scheme    Scheme
	Authority Authority
	Path      Path
	Query     Query
	Fragment  Fragment
}

func (o CRI) String() string {
	if o.Discard != (Discard{}) {
		return "discard"
	}
	return fmt.Sprintf(
		"%s:%s%s%s%s",
		o.Scheme,
		o.Authority,
		o.Path,
		o.Query,
		o.Fragment,
	)
}

func (o *CRI) Parse(data []byte) error {
	var (
		pc   *PC
		err  error
		elem interface{}
	)

	if pc, err = NewPC(data); err != nil || pc.Empty() {
		return err
	}

	// we have checked that pc is not empty, so we have at least one element
	if elem, err = pc.Peek(); err != nil {
		return err
	}

	if isScheme(elem) {
		if err := o.Scheme.Set(elem); err != nil {
			return err
		}

		if elem, err = pc.Next(); err != nil {
			return err
		}

		if err := o.Authority.Set(elem); err != nil {
			return err
		}
	} else if isDiscard(elem) {
		if err := o.Discard.Set(elem); err != nil {
			return err
		}
	}

	if elem, err = pc.Next(); err != nil {
		if err == ErrEndOfArray {
			return nil
		}
		return err
	}

	if err := o.Path.Set(elem); err != nil {
		return err
	}

	if elem, err = pc.Next(); err != nil {
		if err == ErrEndOfArray {
			return nil
		}
		return err
	}

	if err := o.Query.Set(elem); err != nil {
		return err
	}

	if elem, err = pc.Next(); err != nil {
		if err == ErrEndOfArray {
			return nil
		}
		return err
	}

	if err := o.Fragment.Set(elem); err != nil {
		return err
	}

	_, err = pc.Next()
	if err != ErrEndOfArray {
		return fmt.Errorf("spurious elements")
	}

	return nil
}
