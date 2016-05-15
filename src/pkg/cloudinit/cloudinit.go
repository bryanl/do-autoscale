package cloudinit

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"
)

// MIMEType represents a MIME type.
type MIMEType string

const (
	// MIMETypeShellScript is a shell script MIME type.
	MIMETypeShellScript MIMEType = "text/x-shellscript"

	// MIMETypeUnknown is an unknown MIME type.
	MIMETypeUnknown MIMEType = ""
)

// CloudInit represents a cloud init multipart data structure
type CloudInit struct {
	buf      *bytes.Buffer
	envelope *multipart.Writer
	boundary string
}

// New creates an instance of CloudInit
func New() *CloudInit {
	buf := new(bytes.Buffer)
	envelope := multipart.NewWriter(buf)

	return &CloudInit{
		buf:      buf,
		envelope: envelope,
		boundary: envelope.Boundary(),
	}
}

// AddPart adds a part to the cloud init payload.
func (c *CloudInit) AddPart(mimeType MIMEType, name, contents string) error {
	mh := textproto.MIMEHeader{}

	if len(mimeType) > 0 {
		mh.Set("Content-Type", string(mimeType))
		mh.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, name))
		mh.Set("MIME-Version", "1.0")
	}

	part, err := c.envelope.CreatePart(mh)
	if err != nil {
		return err
	}

	if _, err := io.Copy(part, bytes.NewBufferString(contents)); err != nil {
		return err
	}

	return nil
}

// Close closes the envelope.
func (c *CloudInit) Close() error {
	return c.envelope.Close()
}

type proxy struct {
	io.Writer
}

func (c *CloudInit) String() string {
	out := new(bytes.Buffer)
	out.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=\"%s\"\n",
		c.envelope.Boundary()))
	out.WriteString("Mime-Version: 1.0\n\n")

	out.Write(c.buf.Bytes())
	return out.String()
}
