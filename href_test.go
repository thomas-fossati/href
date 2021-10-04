package href

import (
	"encoding/hex"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func MustHexDecode(s string) []byte {
	data, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return data
}

var BadTestVectors = []struct {
	cri         []byte
	expectedErr string
}{
	{
		// echo '["SCHEME"]' | diag2cbor.rb | xxd -p
		cri:         MustHexDecode("8166534348454d45"),
		expectedErr: "scheme-name SCHEME does not match scheme RE ([a-z][a-z0-9+.-]*)",
	},
	{
		// echo '{}' | diag2cbor.rb | xxd -p
		cri:         MustHexDecode("a0"),
		expectedErr: "cbor: cannot unmarshal map into Go value of type []interface {}",
	},
	{
		// echo -n "[ h'01' ]" | diag2cbor.rb | xxd -p
		cri:         MustHexDecode("814101"),
		expectedErr: "expecting scheme or discard, got []uint8",
	},
}

// TODO(tho) since the "official" test vectors are supposed to cover all nominal
// cases, we should repurpose this to explore corner cases.
var GoodTestVectors = []struct {
	cri []byte
	uri string
	// set .criOut if round-tripping is supposed to produce a different value than .cri
	criOut []byte
}{
	{
		// echo '[]' | diag2cbor.rb | xxd -p
		cri: MustHexDecode("80"),
		uri: "",
	},
	{
		// echo '[0]' | diag2cbor.rb | xxd -p
		cri: MustHexDecode("8100"),
		uri: "",
		// echo '[]' | diag2cbor.rb | xxd -p
		criOut: MustHexDecode("80"),
	},
	{
		cri: MustHexDecode("82218144c0a80061"),
		uri: "coaps://192.168.0.97",
	},
	{
		// echo '["coap+tcp", ["acme.example", 5683]]' | diag2cbor.rb | xxd -p
		cri: MustHexDecode("8268636f61702b746370826c61636d652e6578616d706c65191633"),
		uri: "coap+tcp://acme.example:5683",
	},
	{
		// echo '["coap+tcp", ["acme.example", 5683], ["a", "b", "c"]]' | diag2cbor.rb | xxd -p
		cri: MustHexDecode("8368636f61702b746370826c61636d652e6578616d706c6519163383616161626163"),
		uri: "coap+tcp://acme.example:5683/a/b/c",
	},
	{
		// echo '[-1, ["acme.example"]]' | diag2cbor.rb | xxd -p
		cri: MustHexDecode("8220816c61636d652e6578616d706c65"),
		uri: "coap://acme.example",
	},
	{
		// echo '[-1, ["acme.example"], null, ["a=b", "c=d"]]' | diag2cbor.rb | xxd -p
		cri: MustHexDecode("8420816c61636d652e6578616d706c65f68263613d6263633d64"),
		uri: "coap://acme.example?a=b&c=d",
	},
	{
		// echo '[-2, ["acme.example"], null, null, "fragment"]' | diag2cbor.rb | xxd -p
		cri: MustHexDecode("8521816c61636d652e6578616d706c65f6f668667261676d656e74"),
		uri: "coaps://acme.example#fragment",
	},
	{
		// echo '[-3, ["acme.example"]]' | diag2cbor.rb | xxd -p
		cri: MustHexDecode("8222816c61636d652e6578616d706c65"),
		uri: "http://acme.example",
	},
	{
		// echo '[-4, ["acme.example"]]' | diag2cbor.rb | xxd -p
		cri: MustHexDecode("8223816c61636d652e6578616d706c65"),
		uri: "https://acme.example",
	},
	{
		// echo '[-5, ["acme.example"]]' | diag2cbor.rb | xxd -p
		cri: MustHexDecode("8224816c61636d652e6578616d706c65"),
		uri: "urn://acme.example",
	},
	{
		// echo '["myscheme", true]' | diag2cbor.rb | xxd -p
		cri: MustHexDecode("82686d79736368656d65f5"),
		uri: "myscheme:",
	},
	{
		// echo '[3]' | diag2cbor.rb | xxd -p
		cri: MustHexDecode("8103"),
		uri: "../../",
	},
	{
		// echo '[true]' | diag2cbor.rb | xxd -p
		cri: MustHexDecode("81f5"),
		uri: "/",
	},
}

func TestCRI_full_circle(t *testing.T) {
	for i, tv := range GoodTestVectors {
		decoded, err := Parse(tv.cri)
		assert.NoError(t, err, "test case at index %d failed decoding", i)

		uri, err := decoded.ToURI()
		assert.NoError(t, err, "test case at index %d failed translation to URI", i)
		assert.Equal(t, tv.uri, uri.String())

		encoded, err := decoded.ToCBOR()
		assert.NoError(t, err, "test case at index %d failed encoding", i)
		if tv.criOut != nil {
			assert.Equal(t, tv.criOut, encoded)
		} else {
			assert.Equal(t, tv.cri, encoded)
		}
	}
}

func TestCRI_ko(t *testing.T) {
	for _, tv := range BadTestVectors {
		_, err := Parse(tv.cri)
		assert.EqualError(t, err, tv.expectedErr)
	}
}

func Test_repo(t *testing.T) {
	tests, err := LoadTests("./tests.json")
	require.NoError(t, err)

	// XXX(tho) workaround for missing deep copy in reference resolution
	// baseCRI, err := Parse(MustHexDecode(tests.BaseCRI))
	// NoError(t, err)

	baseURI, err := url.Parse(tests.BaseURI)
	require.NoError(t, err)

	for i, tv := range tests.TestVectors {
		baseCRI, err := Parse(MustHexDecode(tests.BaseCRI))
		require.NoError(t, err)

		// CBOR serialization (§5.1 of href-09)
		c, err := Parse(MustHexDecode(tv.CRI))
		require.NoError(t, err, "TC[%d] failed: parsing CRI", i)

		// CRI to URL (§6.1 of href-09)
		u, err := c.ToURI()
		assert.NoError(t, err, "TC[%d] failed: mapping CRI to URI", i)
		assert.Equal(t, tv.URIFromCRI, u.String(), "TC[%d] %v", i, tv)

		// resolve URI reference
		resolvedURI := baseURI.ResolveReference(u)
		assert.Equal(t, tv.ResolvedURI, resolvedURI.String())

		// (§5.3 of href-09)
		expected := MustHexDecode(tv.ResolvedCRI)
		resolvedCRI := baseCRI.ResolveReference(c)
		got, err := resolvedCRI.ToCBOR()
		assert.NoError(t, err, "TC[%d] failed: resolving CRI reference", i)
		assert.Equal(t, expected, got, "TC[%d] want: %x, got %x", i, expected, got)
	}
}
