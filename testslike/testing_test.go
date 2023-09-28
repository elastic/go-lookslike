package testslike

import (
	"testing"

	"github.com/elastic/go-lookslike"
)

func TestTest(t *testing.T) {
	validator := lookslike.MustCompile(map[string]interface{}{
		"foo": "bar",
		"a":   123,
	})
	val := map[string]interface{}{
		"foo": "bar",
		"a":   123,
	}
	Test(t, validator, val)
}
