package href

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func MustHexDecode(t *testing.T, s string) []byte {
	data, err := hex.DecodeString(s)
	require.Nil(t, err)
	return data
}

func TestCRI_ok(t *testing.T) {
	tvs := []struct {
		raw      []byte
		expected string
	}{
		{
			// echo '[]' | diag2cbor.rb | xxd -p
			raw:      MustHexDecode(t, "80"),
			expected: "://",
		},
		{
			// echo '["coap+tcp", ["acme.example", 5683]]' | diag2cbor.rb | xxd -p
			raw:      MustHexDecode(t, "8268636f61702b746370826c61636d652e6578616d706c65191633"),
			expected: "coap+tcp://acme.example:5683",
		},
		{
			// echo '["coap+tcp", ["acme.example", 5683], ["a", "b", "c"]]' | diag2cbor.rb | xxd -p
			raw:      MustHexDecode(t, "8368636f61702b746370826c61636d652e6578616d706c6519163383616161626163"),
			expected: "coap+tcp://acme.example:5683/a/b/c",
		},
		{
			// echo '[-1, ["acme.example"]]' | diag2cbor.rb | xxd -p
			raw:      MustHexDecode(t, "8220816c61636d652e6578616d706c65"),
			expected: "coap://acme.example",
		},
		{
			// echo '[-1, ["acme.example"], null, ["a=b", "c=d"]]' | diag2cbor.rb | xxd -p
			raw:      MustHexDecode(t, "8420816c61636d652e6578616d706c65f68263613d6263633d64"),
			expected: "coap://acme.example?a=b&c=d",
		},
		{
			// echo '[-2, ["acme.example"], null, null, "fragment"]' | diag2cbor.rb | xxd -p
			raw:      MustHexDecode(t, "8521816c61636d652e6578616d706c65f6f668667261676d656e74"),
			expected: "coaps://acme.example#fragment",
		},
		{
			// echo '[-3, ["acme.example"]]' | diag2cbor.rb | xxd -p
			raw:      MustHexDecode(t, "8222816c61636d652e6578616d706c65"),
			expected: "http://acme.example",
		},
		{
			// echo '[-4, ["acme.example"]]' | diag2cbor.rb | xxd -p
			raw:      MustHexDecode(t, "8223816c61636d652e6578616d706c65"),
			expected: "https://acme.example",
		},
		{
			// echo '[-5, ["acme.example"]]' | diag2cbor.rb | xxd -p
			raw:      MustHexDecode(t, "8224816c61636d652e6578616d706c65"),
			expected: "scheme-id(-5)://acme.example",
		},
		{
			// echo '["myscheme", true]' | diag2cbor.rb | xxd -p
			raw:      MustHexDecode(t, "82686d79736368656d65f5"),
			expected: "myscheme:",
		},
		{
			// echo '["urn", null]' | diag2cbor.rb | xxd -p
			raw:      MustHexDecode(t, "826375726ef6"),
			expected: "urn:",
		},
		{
			// echo '[0]' | diag2cbor.rb | xxd -p
			raw:      MustHexDecode(t, "8100"),
			expected: "discard",
		},
		{
			// echo '[127]' | diag2cbor.rb | xxd -p
			raw:      MustHexDecode(t, "81187f"),
			expected: "discard",
		},
		{
			// echo '[true]' | diag2cbor.rb | xxd -p
			raw:      MustHexDecode(t, "81f5"),
			expected: "discard",
		},
	}

	for _, tv := range tvs {
		var actual CRI
		err := actual.Parse(tv.raw)
		assert.Nil(t, err)
		assert.Equal(t, tv.expected, actual.String())
	}
}

func TestCRI_ko(t *testing.T) {
	tvs := []struct {
		raw      []byte
		expected string
	}{
		{
			// echo '["SCHEME"]' | diag2cbor.rb | xxd -p
			raw:      MustHexDecode(t, "8166534348454d45"),
			expected: "scheme-name SCHEME does not match scheme RE ([a-z][a-z0-9+.-]*)",
		},
		{
			// echo '{}' | diag2cbor.rb | xxd -p
			raw:      MustHexDecode(t, "a0"),
			expected: "cbor: cannot unmarshal map into Go value of type []interface {}",
		},
	}

	for _, tv := range tvs {
		var actual CRI
		err := actual.Parse(tv.raw)
		assert.EqualError(t, err, tv.expected)
	}
}
