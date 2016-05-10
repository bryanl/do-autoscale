package autoscale

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTemplate_IsValid(t *testing.T) {
	cases := []struct {
		Name    string
		IsValid bool
	}{
		{Name: "1234", IsValid: true},
		{Name: "-1234", IsValid: false},
		{Name: "a-template", IsValid: true},
	}

	for _, c := range cases {
		tmpl := Template{
			Name: c.Name,
		}

		assert.Equal(t, c.IsValid, tmpl.IsValid())
	}
}

func TestGroup_IsValid(t *testing.T) {
	cases := []struct {
		Name    string
		IsValid bool
	}{
		{Name: "1234", IsValid: true},
		{Name: "-1234", IsValid: false},
		{Name: "a-template", IsValid: true},
	}

	for _, c := range cases {
		group := Group{
			Name: c.Name,
		}

		assert.Equal(t, c.IsValid, group.IsValid())
	}
}
