package href

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/fxamacker/cbor/v2"
)

type CRI struct {
	Discard   Discard
	Scheme    Scheme
	Authority Authority
	Path      Path
	Query     Query
	Fragment  Fragment
}

// Parse ingest a CRI Reference in transfer form into its abstract form
func Parse(rawCRI []byte) (*CRI, error) {
	var (
		pc   *PC
		eof  bool
		err  error
		elem interface{}
		cri  CRI
	)

	if pc, err = NewPC(rawCRI); err != nil {
		return nil, err
	}

	if pc.Empty() {
		// §5.2 If the array is entirely empty, replace it with [0].
		_ = cri.Discard.Set(uint64(0))
		return &cri, nil
	}

	// If we reached here is because pc!=empty, hence we know we have at least
	// one element
	elem, _ = pc.Peek()

	if isScheme(elem) {
		if err := cri.Scheme.Set(elem); err != nil {
			return nil, err
		}

		// since "null" is an acceptable value for authority and trailing null's
		// are suppressed, if we get EOF here, we need to set the authority
		// explicitly and declare success.
		if elem, eof = pc.Next(); eof {
			cri.Authority.SetNull()
			return &cri, nil
		}

		if err := cri.Authority.Set(elem); err != nil {
			return nil, err
		}
	} else if isDiscard(elem) {
		if err := cri.Discard.Set(elem); err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("expecting scheme or discard, got %T", elem)
	}

	if elem, eof = pc.Next(); eof {
		return &cri, nil
	}

	if err := cri.Path.Set(elem); err != nil {
		return nil, err
	}

	if elem, eof = pc.Next(); eof {
		return &cri, nil
	}

	if err := cri.Query.Set(elem); err != nil {
		return nil, err
	}

	if elem, eof = pc.Next(); eof {
		return &cri, nil
	}

	if err := cri.Fragment.Set(elem); err != nil {
		return nil, err
	}

	if _, eof = pc.Next(); !eof {
		return nil, fmt.Errorf("spurious trailing elements")
	}

	return &cri, nil
}

func (o *CRI) ToCBOR() ([]byte, error) {
	var cri []interface{}

	if o.Scheme.IsSet() {
		cri = append(cri, o.Scheme.Get())

		if o.Authority.IsNull {
			cri = append(cri, nil)
		} else if o.Authority.IsTrue {
			cri = append(cri, true)
		} else {
			var authority []interface{}
			authority = append(authority, o.Authority.Host.Get())
			if o.Authority.Port.IsSet() {
				authority = append(authority, o.Authority.Port.Get())
			}
			cri = append(cri, authority)
		}
	} else if o.Discard.IsSet() {
		cri = append(cri, o.Discard.Get())
	} else {
		return nil, fmt.Errorf("neither an absolute CRI nor a relative reference")
	}

	var pathQueryAndFrag []interface{}

	if o.Fragment.IsSet() {
		pathQueryAndFrag = append(pathQueryAndFrag, o.Fragment.Get())
	}

	if o.Query.IsSet() {
		pathQueryAndFrag = prepend(o.Query.GetUnescaped(), pathQueryAndFrag)
	} else if len(pathQueryAndFrag) > 0 {
		pathQueryAndFrag = prepend(nil, pathQueryAndFrag)
	}

	if o.Path.IsSet() {
		pathQueryAndFrag = prepend(o.Path.GetUnescaped(), pathQueryAndFrag)
	} else if len(pathQueryAndFrag) > 0 {
		pathQueryAndFrag = prepend(nil, pathQueryAndFrag)
	}

	if len(pathQueryAndFrag) > 0 {
		cri = append(cri, pathQueryAndFrag...)
	}

	// trailing null(s) suppression
	for i := len(cri) - 1; i >= 0; i-- {
		if cri[i] != nil {
			break
		}
		cri = cri[:i]
	}

	// As a special case, an empty array is sent in place for a remaining [0] (URI "").
	if isDiscardZero(cri) {
		cri = []interface{}{}
	}

	return cbor.Marshal(cri)
}

func isDiscardZero(cri []interface{}) bool {
	if len(cri) != 1 {
		return false
	}

	if v, ok := cri[0].(uint64); !ok || v != 0 {
		return false
	}

	return true
}

func prepend(val interface{}, dst []interface{}) []interface{} {
	var p []interface{}
	return append(append(p, val), dst...)
}

// ToURI convert a CRI reference to a URI reference by determining the
// components of the URI reference according to the steps in §6.1 of href-09 and
// then recomposing the components to a URI reference string as specified in
// §5.3 of RFC3986.
func (o *CRI) ToURI() (*url.URL, error) {
	scheme := o.toURISchemeRules()
	host := o.toURIAuthorityRules()
	path := o.toURIPathRules()
	query := o.toURIQueryRules()
	fragment := o.toURIFragmentRules()

	// From https://pkg.go.dev/net/url#URL
	//
	// A URL represents a parsed URL (technically, a URI reference).
	// The general form represented is:
	// 	[scheme:][//[userinfo@]host][/]path[?query][#fragment]
	// URLs that do not start with a slash after the scheme are interpreted as:
	// 	scheme:opaque[?query][#fragment]

	//	if scheme != "" && host == "" && len(path) > 0 && path[0] != '/' {
	if scheme != "" && host == "" {
		return &url.URL{
			Scheme:   scheme,
			Opaque:   path,
			RawQuery: query,
			Fragment: fragment,
		}, nil
	}

	return &url.URL{
		Scheme:   scheme,
		Host:     host,
		Path:     path,
		RawQuery: query,
		Fragment: fragment,
	}, nil
}

func (o *CRI) toURISchemeRules() string {
	return o.Scheme.String()
}

func (o *CRI) toURIAuthorityRules() string {
	return o.Authority.String()
}

func (o *CRI) toURIPathRules() string {
	var (
		path   string
		rooted bool
	)

	if o.Discard.IsSet() {
		path = o.Discard.ComputePathPrefix()
	} else if o.Authority.IsTrue { // no auth, no slash
		rooted = false
	} else {
		rooted = true
	}

	if o.Path.NumSegments() > 0 {
		if rooted {
			path = "/"
		}

		path += strings.Join(o.Path.GetEscapedSegments(), "/")
	}

	// sanity checks
	//
	// If the authority component is present (not null or true) and the path
	// component does not match the "path-abempty" rule the conversion fails.
	// XXX(tho) should bail out
	if o.Authority.IsSet() {
		if !matchPathAbEmpty(path) {
			fmt.Println("authority is present but path is not path-abempty")
		}
	} else if o.Scheme.IsSet() {
		// If the authority component is not present, but the scheme component
		// is, and the path component does not match the "path-absolute",
		// "path-rootless" (authority == true) or "path-empty" rule the
		// conversion fails.
		// XXX(tho) should bail out
		if !matchPathAbsolute(path) && !matchPathRootless(path) && !matchPathEmpty(path) {
			fmt.Println("authority is not present but scheme is and path is not absolute, rootless or empty")
		}
	} else {
		// If neither the authority component nor the scheme component are
		// present, and the path component does not match the "path-absolute",
		// "path-noscheme" or "path-empty" rule the conversion fails.
		if !matchPathAbsolute(path) && !matchPathNoScheme(path) && !matchPathEmpty(path) {
			fmt.Println("authority and scheme not present and path is not absolute, noscheme or empty")
		}
	}

	return path
}

func (o *CRI) toURIQueryRules() string {
	return o.Query.String()
}

func (o *CRI) toURIFragmentRules() string {
	return o.Fragment.String()
}

// ResolveReference resolves a CRI reference to an absolute CRI from an absolute
// base CRI o, per href-09 Section 5.3. The CRI reference may be relative or
// absolute. ResolveReference always returns a new CRI instance, even if the
// returned CRI is identical to either the base or reference. If ref is an
// absolute CRI, then ResolveReference ignores base and returns a copy of ref.
func (o CRI) ResolveReference(ref *CRI) *CRI { // nolint: gocritic
	var resolvedCRI CRI

	if ref.IsAbs() {
		resolvedCRI = *ref
		return &resolvedCRI
	}

	// 1. Establish the base CRI of the CRI reference and express it in the form
	//    of an abstract absolute CRI reference.
	//
	// XXX(tho): we assume 'o' is that "abstract absolute CRI reference" without
	// checking.  What could we do here to make sure 'o' is fit for purpose?

	// 2. Initialize a buffer with the sections from the base CRI.
	//
	// XXX(tho) shallow copy, we need to deep copy instead.
	resolvedCRI = o

	// 3. If the value of discard is true in the CRI reference, replace the path
	//    in the buffer with the empty array, unset query and fragment, and set
	//    a true authority to null. If the value of discard is an unsigned
	//    number, remove as many elements from the end of the path array; if it
	//    is non-zero, unset query and fragment.
	//
	// NOTE: discard = DISCARD-ALL is implicitly the case when scheme and/or
	//       authority are present in the reference.
	discardAll := ref.Discard.IsTrue() || ref.Scheme.IsSet() || ref.Authority.IsSet()

	if ref.Discard.IsSet() || discardAll {
		if discardAll {
			resolvedCRI.Path.Reset()
			resolvedCRI.Query.Reset()
			resolvedCRI.Fragment.Reset()
			if resolvedCRI.Authority.IsTrue {
				resolvedCRI.Authority.SetNull()
			}
		} else { // unsigned number
			n, ok := ref.Discard.Get().(uint64)
			if !ok {
				panic("discard is not a number")
			}
			resolvedCRI.Path.TrimN(n)
			if n > 0 {
				resolvedCRI.Query.Reset()
				resolvedCRI.Fragment.Reset()
			}
		}
	}

	// (Unconditionally) set discard to true in the buffer.
	_ = resolvedCRI.Discard.Set(true)

	// 4. If the path section is set in the CRI reference, append all elements
	//    from the path array to the array in the path section in the buffer;
	//    unset query and fragment.
	if ref.Path.IsSet() {
		resolvedCRI.Path.Append(ref.Path.GetSegments())
		resolvedCRI.Query.Reset()
		resolvedCRI.Fragment.Reset()
	}

	// 5. Apart from the path and discard, copy all non-null sections from the
	//    CRI reference to the buffer in sequence; unset fragment if query is
	//    non-null and thus copied.
	if ref.Scheme.IsSet() {
		resolvedCRI.Scheme = ref.Scheme
	}

	if ref.Authority.IsSet() {
		// TODO(tho) check that cloning is deep enough
		resolvedCRI.Authority = ref.Authority
	}

	if ref.Query.IsSet() {
		resolvedCRI.Query = ref.Query
		resolvedCRI.Fragment.Reset()
	}

	if ref.Fragment.IsSet() {
		resolvedCRI.Fragment = ref.Fragment
	}

	return &resolvedCRI
}

func (o *CRI) IsAbs() bool {
	// A CRI reference is considered _absolute_ if
	// a) it is well-formed (TODO(tho)), and
	// b) the sequence of sections starts with a non-null "scheme".
	return o.Scheme.IsSet()
}
