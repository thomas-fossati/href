package href

import (
	"encoding/json"
	"os"
)

type TestVector struct {
	URI         string `json:"uri"`
	CRI         string `json:"cri"`
	URIFromCRI  string `json:"uri-from-cri"`
	ResolvedCRI string `json:"resolved-cri"`
	ResolvedURI string `json:"resolved-uri"`
}

type Tests struct {
	BaseURI     string       `json:"base-uri"`
	BaseCRI     string       `json:"base-cri"`
	TestVectors []TestVector `json:"test-vectors"`
}

func LoadTests(testsJSONFile string) (Tests, error) {
	j, err := os.ReadFile(testsJSONFile)
	if err != nil {
		return Tests{}, err
	}

	var tests Tests

	err = json.Unmarshal(j, &tests)
	if err != nil {
		return Tests{}, err
	}

	return tests, nil
}
