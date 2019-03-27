package isdefs

import (
	"github.com/elastic/lookslike/lookslike"
	"github.com/elastic/lookslike/lookslike/validator"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestOptional(t *testing.T) {
	m := validator.Map{
		"foo": "bar",
	}

	validator := lookslike.MustCompile(validator.Map{
		"non": Optional(IsEqual("foo")),
	})

	require.True(t, validator(m).Valid)
}
