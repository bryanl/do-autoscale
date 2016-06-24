package autoscale

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultTagName(t *testing.T) {
	gn := "groupname"
	tag := defaultTagName(gn)
	assert.Equal(t, "as:7c13c7fc", tag)
}
