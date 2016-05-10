package autoscale

import (
	"encoding/json"
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

func TestTemplate_SSHKeys(t *testing.T) {
	tmpl := Template{RawSSHKeys: "1,2,3"}
	expected := []string{"1", "2", "3"}
	assert.Equal(t, expected, tmpl.SSHKeys())
}

func TestTemplate_Marshal(t *testing.T) {
	tmpl := Template{RawSSHKeys: "1,2,3"}
	b, err := json.Marshal(&tmpl)
	assert.NoError(t, err)

	var tmpl2 Template
	err = json.Unmarshal(b, &tmpl2)
	assert.NoError(t, err)

	assert.Equal(t, "1,2,3", tmpl2.RawSSHKeys)
}
