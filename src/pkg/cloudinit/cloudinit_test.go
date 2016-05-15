package cloudinit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCloudInit(t *testing.T) {
	ci := New()
	err := ci.AddPart(MIMETypeShellScript, "hello.txt", "echo hi")
	assert.NoError(t, err)

	err = ci.Close()
	assert.NoError(t, err)

	out := ci.String()
	assert.NotEmpty(t, out)

}
