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
		eof  bool
		err  error
		elem interface{}
	)

	if pc, err = NewPC(data); err != nil || pc.Empty() {
		return err
	}

	// we have checked that pc is not empty, so we have at least one element
	elem, _ = pc.Peek()

	if isScheme(elem) {
		if err := o.Scheme.Set(elem); err != nil {
			return err
		}

		if elem, eof = pc.Next(); eof {
			return fmt.Errorf("expecting authority, found EOF")
		}

		if err := o.Authority.Set(elem); err != nil {
			return err
		}
	} else if isDiscard(elem) {
		if err := o.Discard.Set(elem); err != nil {
			return err
		}
	}

	if elem, eof = pc.Next(); eof {
		return nil
	}

	if err := o.Path.Set(elem); err != nil {
		return err
	}

	if elem, eof = pc.Next(); eof {
		return nil
	}

	if err := o.Query.Set(elem); err != nil {
		return err
	}

	if elem, eof = pc.Next(); eof {
		return nil
	}

	if err := o.Fragment.Set(elem); err != nil {
		return err
	}

	if _, eof = pc.Next(); !eof {
		return fmt.Errorf("spurious trailing elements")
	}

	return nil
}
