package fmt

import (
	"testing"
)

func TestT(t *testing.T) {
	for k, v := range Formats() {
		if k != v.Algorithm {
			t.Errorf("%s != %s\n", k, v.Algorithm)
		}
		transformed, err := T(v.Example, v.Pattern)
		if err != nil {
			t.Errorf("Pattern %s for %s is not valid for hash %s\n", v.Pattern, v.Algorithm, v.Example)
		}
		if transformed == "" {
			t.Errorf("Example %s for %s does not comply with %s\n", v.Example, v.Algorithm, v.Pattern)
		}
	}
}
