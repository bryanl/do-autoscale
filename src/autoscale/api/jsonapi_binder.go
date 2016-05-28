package api

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/manyminds/api2go/jsonapi"
)

const (
	// MIMEApplicationJSONAPI is the MIME type for JSONAPI.
	MIMEApplicationJSONAPI = "application/vnd.api+json"
)

type jsonAPIBinder struct{}

var _ echo.Binder = (*jsonAPIBinder)(nil)

func (b *jsonAPIBinder) Bind(i interface{}, c echo.Context) error {
	req := c.Request()
	ctype := req.Header().Get(echo.HeaderContentType)

	if req.Body() == nil {
		err := echo.NewHTTPError(http.StatusBadRequest, "request body can't be empty")
		return err
	}

	switch {
	case strings.HasPrefix(ctype, MIMEApplicationJSONAPI):
		b, err := ioutil.ReadAll(req.Body())
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		if err := jsonapi.Unmarshal(b, i); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

	default:
		return echo.ErrUnsupportedMediaType
	}

	return nil
}
