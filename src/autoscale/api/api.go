package api

import (
	"autoscale/gen"
	"encoding/json"
	"net/http"
	"pkg/ctxutil"
	"pkg/echologger"
	"strconv"
	"strings"

	"golang.org/x/net/context"

	"autoscale"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
	"github.com/satori/go.uuid"
)

type errorWrapper struct {
	Errors errorMsg `json:"errors"`
}

type errorMsg struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Detail string `json:"detail"`
}

func errorHandler(err error, c echo.Context) {
	log := ctxutil.LogFromContext(c)

	code := http.StatusInternalServerError
	msg := http.StatusText(code)
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		msg = he.Message
	}

	reqID := c.Request().Header().Get("X-Request-Id")

	// TODO only show this if we are in dev mode
	msg = err.Error()

	if !c.Response().Committed() {
		req := c.Request()
		ctype := req.Header().Get(echo.HeaderContentType)

		switch {
		case strings.HasPrefix(ctype, echo.MIMEApplicationJSON):
			apiError := errorWrapper{
				Errors: errorMsg{
					ID:     reqID,
					Status: strconv.Itoa(code),
					Detail: msg,
				},
			}

			if err := c.JSON(code, &apiError); err != nil {
				log.WithError(err).Error("uaable to create api error mesage")
			}

		default:
			c.String(code, msg)
		}
	}

	log.WithError(err).WithField("request_id", reqID).Error("api error")
}

// API is the autoscale API.
type API struct {
	Mux  http.Handler
	repo autoscale.Repository
	ctx  context.Context

	templateResourceFactory    func() Resource
	groupResourceFactory       func() Resource
	userConfigResourceFactory  func() Resource
	groupConfigResourceFactory func() Resource
	timeSeriesResourceFactory  func(groupID string) Resource
}

// New creates an instance of API.
func New(ctx context.Context, repo autoscale.Repository) *API {
	e := echo.New()

	std := standard.WithConfig(engine.Config{})
	std.SetHandler(e)

	a := &API{
		Mux:  std,
		repo: repo,
		ctx:  ctx,

		templateResourceFactory: func() Resource {
			return &templateResource{repo: repo}
		},
		groupResourceFactory: func() Resource {
			return &groupResource{repo: repo}
		},
		userConfigResourceFactory: func() Resource {
			return &userConfigResource{}
		},
		groupConfigResourceFactory: func() Resource {
			return &groupConfigResource{repo: repo}
		},
		timeSeriesResourceFactory: func(groupID string) Resource {
			return &timeSeriesResource{
				groupID: groupID,
				repo:    repo,
			}
		},
	}

	log := ctxutil.LogFromContext(ctx)

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.SetContext(ctx)

			reqID := c.Request().Header().Get("X-Request-Id")
			if reqID == "" {
				reqID = uuid.NewV4().String()
				c.Request().Header().Set("X-Request-Id", reqID)

			}

			// NOTE this causes a stack issue
			// newCtx := context.WithValue(c, "RequestID", reqID)
			// newCtx = context.WithValue(newCtx, "log", log)
			// c.SetNetContext(newCtx)

			return next(c)
		}
	})

	logmw := echologger.NewWithNameAndLogger("autoscale", log)
	e.Use(logmw)

	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 1 << 10, // 1 KB
	}))

	g := e.Group("/api")

	g.Get("/templates/:id", a.getTemplate)
	g.Get("/templates", a.listTemplates)
	g.Post("/templates", a.createTemplate)
	g.Delete("/templates/:id", a.deleteTemplate)
	g.Get("/groups", a.listGroups)
	g.Get("/groups/:id", a.getGroup)
	g.Get("/groups/:id/time_series", a.getGroupTimeSeries)
	g.Post("/groups", a.createGroup)
	g.Delete("/groups/:id", a.deleteGroup)
	g.Put("/groups/:id", a.updateGroup)
	g.Get("/user_config", a.userConfig)
	g.Get("/group_configs", a.groupConfig)

	e.Get("/*", func(c echo.Context) error {
		w := c.Response().(*standard.Response).ResponseWriter
		r := c.Request().(*standard.Request).Request

		http.FileServer(
			&assetfs.AssetFS{
				Asset:     gen.Asset,
				AssetDir:  gen.AssetDir,
				AssetInfo: gen.AssetInfo,
				Prefix:    "static"}).
			ServeHTTP(w, r)
		return nil
	})

	e.Get("/assets/*", func(c echo.Context) error {
		w := c.Response().(*standard.Response).ResponseWriter
		r := c.Request().(*standard.Request).Request

		http.StripPrefix(
			"/assets/",
			http.FileServer(&assetfs.AssetFS{
				Asset:     gen.Asset,
				AssetDir:  gen.AssetDir,
				AssetInfo: gen.AssetInfo,
				Prefix:    "static/assets"})).
			ServeHTTP(w, r)
		return nil
	})

	e.SetHTTPErrorHandler(errorHandler)

	return a
}

func buildResponse(c echo.Context, resp Response) error {
	j, err := json.Marshal(resp.Result())
	if err != nil {
		return err
	}

	return c.JSONBlob(resp.StatusCode(), j)
}

func (a *API) listTemplates(c echo.Context) error {
	resp, err := a.templateResourceFactory().FindAll(c)
	if err != nil {
		return err
	}

	return buildResponse(c, resp)
}

func (a *API) getTemplate(c echo.Context) error {
	id := c.Param("id")
	resp, err := a.templateResourceFactory().FindOne(c, id)
	if err != nil {
		return err
	}

	return buildResponse(c, resp)
}

func (a *API) createTemplate(c echo.Context) error {
	var wrapper templateWrapper
	if err := c.Bind(&wrapper); err != nil {
		return err
	}

	resp, err := a.templateResourceFactory().Create(c, wrapper.Template)
	if err != nil {
		return err
	}

	return buildResponse(c, resp)
}

func (a *API) deleteTemplate(c echo.Context) error {
	id := c.Param("id")
	resp, err := a.templateResourceFactory().Delete(c, id)
	if err != nil {
		return err
	}

	return buildResponse(c, resp)
}

func (a *API) listGroups(c echo.Context) error {
	resp, err := a.groupResourceFactory().FindAll(c)
	if err != nil {
		return err
	}

	return buildResponse(c, resp)
}

func (a *API) getGroup(c echo.Context) error {
	id := c.Param("id")
	resp, err := a.groupResourceFactory().FindOne(c, id)
	if err != nil {
		return err
	}

	return buildResponse(c, resp)
}

func (a *API) createGroup(c echo.Context) error {
	var wrapper groupWrapper
	if err := c.Bind(&wrapper); err != nil {
		return err
	}

	resp, err := a.groupResourceFactory().Create(c, wrapper.Group)

	if err != nil {
		return err
	}

	return buildResponse(c, resp)
}

func (a *API) deleteGroup(c echo.Context) error {
	id := c.Param("id")
	resp, err := a.groupResourceFactory().Delete(c, id)
	if err != nil {
		return err
	}

	return buildResponse(c, resp)

}

func (a *API) updateGroup(c echo.Context) error {
	id := c.Param("id")
	var wrapper groupWrapper
	if err := c.Bind(&wrapper); err != nil {
		return err
	}

	g := wrapper.Group
	g.ID = id

	resp, err := a.groupResourceFactory().Update(c, g)

	if err != nil {
		return err
	}

	return buildResponse(c, resp)
}

func (a *API) userConfig(c echo.Context) error {
	resp, err := a.userConfigResourceFactory().FindAll(c)
	if err != nil {
		return err
	}

	return buildResponse(c, resp)
}

func (a *API) groupConfig(c echo.Context) error {
	resp, err := a.groupConfigResourceFactory().FindAll(c)
	if err != nil {
		return err
	}

	return buildResponse(c, resp)
}

func (a *API) getGroupTimeSeries(c echo.Context) error {
	id := c.Param("id")
	resp, err := a.timeSeriesResourceFactory(id).FindAll(c)
	if err != nil {
		return err
	}

	return buildResponse(c, resp)
}
