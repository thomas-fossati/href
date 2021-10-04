package href

import (
	"fmt"
	"net"
)

type (
	Authority struct {
		Host   Host
		Port   Port
		IsNull bool // no authority, leading slash
		IsTrue bool // no authority, no slash
	}
	Host struct {
		val interface{}
	}
	Port struct {
		val *uint64
	}
)

func (o *Port) Set(v interface{}) error {
	var (
		ok bool
		p  uint64
	)

	if p, ok = v.(uint64); !ok {
		return fmt.Errorf("unexpected port type: %T", v)
	}

	if p > 65535 {
		return fmt.Errorf("port number must be in range 0..65535: got %d", p)
	}

	o.val = &p

	return nil
}

func (o Port) IsSet() bool {
	return o.val != nil
}

func (o Port) Get() uint64 {
	if o.IsSet() {
		return *o.val
	}
	return ^uint64(0)
}

func (o Port) String() string {
	if !o.IsSet() {
		return ""
	}
	return fmt.Sprintf(":%d", *o.val)
}

func (o *Host) Set(v interface{}) error {
	switch t := v.(type) {
	case string:
		// host-name
		o.val = t
		return nil
	case []byte:
		// host-ip
		l := len(t)
		if l != 4 && l != 16 {
			return fmt.Errorf("host-ip must be 4 or 16 bytes, got %d", l)
		}
		o.val = net.IP(t)
	default:
		return fmt.Errorf("unknown host type: %T", t)
	}

	return nil
}

func (o Host) IsSet() bool {
	return o.val != nil
}

func (o Host) Get() interface{} {
	return o.val
}

func (o Host) String() string {
	switch t := o.val.(type) {
	case string:
		// TODO(tho) escape hostname
		// href-06, §6.1:
		// The "host-name" is turned into a single string by joining the
		// elements separated by dots (".").  Any character in the value of a
		// "host-name" item that is not in the set of unreserved characters
		// (Section 2.3 of [RFC3986]) or "sub-delims" (Section 2.2 of
		// [RFC3986]) MUST be percent-encoded.
		return t
	case net.IP:
		return t.String()
	default:
		return ""
	}
}

func (o Authority) String() string {
	if o.IsNull || o.IsTrue || o.Host.String() == "" {
		return ""
	}

	return fmt.Sprintf("%s%s", o.Host, o.Port)
}

func (o *Authority) Set(val interface{}) error {
	switch t := val.(type) {
	case []interface{}:
		if err := o.SetHostPort(t); err != nil {
			return err
		}
	case nil:
		o.SetNull()
	case bool:
		if !t {
			return fmt.Errorf("unexpected authority type: false")
		}
		o.SetTrue()
	default:
		return fmt.Errorf("unexpected authority type: %T", val)
	}

	return nil
}

func (o *Authority) SetNull() {
	o.IsNull = true

	o.IsTrue = false
	o.Host.val = nil
	o.Port.val = nil
}

func (o *Authority) SetTrue() {
	o.IsTrue = true

	o.IsNull = false
	o.Host.val = nil
	o.Port.val = nil
}

func (o *Authority) SetHostPort(val []interface{}) error {
	switch len(val) {
	case 2:
		// host + port
		if err := o.Port.Set(val[1]); err != nil {
			return err
		}
		fallthrough
	case 1:
		// host
		if err := o.Host.Set(val[0]); err != nil {
			return err
		}
		o.IsTrue = false
		o.IsNull = false
	default:
		return fmt.Errorf("wrong number of elements in authority: %d", len(val))
	}

	return nil
}

func (o *Authority) IsSet() bool {
	return !o.IsNull && !o.IsTrue && o.Host.IsSet() // port is optional
}
