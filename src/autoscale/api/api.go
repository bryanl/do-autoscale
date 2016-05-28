package api

import (
	"autoscale/gen"
	"encoding/json"
	"net/http"
	"pkg/ctxutil"
	"pkg/echologger"

	"golang.org/x/net/context"

	"autoscale"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
	"github.com/manyminds/api2go/jsonapi"
	"github.com/satori/go.uuid"
)

type errorMsg struct {
	Title string `json:"title"`
}

func writeError(w http.ResponseWriter, msg string, code int) {
	em := errorMsg{
		Title: msg,
	}

	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(&em)
}

// API is the autoscale API.
type API struct {
	Mux  http.Handler
	repo autoscale.Repository
	ctx  context.Context
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
	}

	log := ctxutil.LogFromContext(ctx)

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			reqID := c.Request().Header().Get("X-Request-Id")
			if reqID == "" {
				reqID = uuid.NewV4().String()
				c.Request().Header().Set("X-Request-Id", reqID)

			}

			newCtx := context.WithValue(c, "RequestID", reqID)
			newCtx = context.WithValue(newCtx, "log", log)
			c.SetNetContext(newCtx)

			return next(c)
		}
	})

	logmw := echologger.NewWithNameAndLogger("autoscale", log)
	e.Use(logmw)

	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 1 << 10, // 1 KB
	}))

	g := e.Group("/api")

	g.Get("/templates/:id", a.GetTemplate)
	g.Get("/templates", a.listTemplates)
	g.Post("/templates", a.createTemplate)
	g.Delete("/templates/:id", a.deleteTemplate)
	g.Get("/groups", a.listGroups)
	g.Get("/groups/:id", a.getGroup)
	g.Post("/groups", a.createGroup)
	g.Delete("/groups/:id", a.deleteGroup)
	g.Put("/groups/:id", a.updateGroup)
	g.Get("/user-configs", a.userConfig)

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

	return a
}

func (a *API) listTemplates(c echo.Context) error {
	log := ctxutil.LogFromContext(c)
	tmpls, err := a.repo.ListTemplates(a.ctx)
	if err != nil {
		log.WithError(err).Error("list templates")

		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	r := autoscale.TemplatesResponse{
		Templates: tmpls,
	}

	return c.JSON(http.StatusOK, r)
}

func (a *API) GetTemplate(c echo.Context) error {
	log := ctxutil.LogFromContext(c)

	id := c.Param("id")

	tmpl, err := a.repo.GetTemplate(a.ctx, id)
	if err != nil {
		log.WithError(err).Error("retrieve template")
		return echo.ErrNotFound
	}

	return c.JSON(http.StatusOK, tmpl)
}

func (a *API) createTemplate(c echo.Context) error {
	log := ctxutil.LogFromContext(c)

	var ctr autoscale.CreateTemplateRequest
	if err := c.Bind(&ctr); err != nil {
		return err
	}

	tmpl, err := a.repo.CreateTemplate(a.ctx, ctr)
	if err != nil {
		log.WithError(err).Error("Delete template")

		return echo.NewHTTPError(http.StatusBadRequest)
	}

	return c.JSON(http.StatusCreated, tmpl)
}

func (a *API) deleteTemplate(c echo.Context) error {
	log := ctxutil.LogFromContext(c)

	id := c.Param("id")

	err := a.repo.DeleteTemplate(a.ctx, id)
	if err != nil {
		log.WithError(err).Error("Delete template")
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	return c.NoContent(204)
}

func (a *API) listGroups(c echo.Context) error {
	log := ctxutil.LogFromContext(c)

	groups, err := a.repo.ListGroups(a.ctx)
	if err != nil {
		log.WithError(err).Error("list groups")
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, groups)
}

func (a *API) getGroup(c echo.Context) error {
	log := ctxutil.LogFromContext(c)

	id := c.Param("id")

	group, err := a.repo.GetGroup(a.ctx, id)
	if err != nil {
		log.WithError(err).Error("Get group")
		return echo.ErrNotFound
	}

	return c.JSON(http.StatusOK, group)
}

func (a *API) createGroup(c echo.Context) error {
	log := ctxutil.LogFromContext(c)

	var cgr autoscale.CreateGroupRequest
	if err := c.Bind(&cgr); err != nil {
		return err
	}

	g, err := a.repo.CreateGroup(a.ctx, cgr)
	if err != nil {
		log.WithError(err).Error("create group")
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	return c.JSON(http.StatusCreated, g)
}

func (a *API) deleteGroup(c echo.Context) error {
	log := ctxutil.LogFromContext(c)

	id := c.Param("id")

	err := a.repo.DeleteGroup(a.ctx, id)
	if err != nil {
		log.WithError(err).Error("Delete group")
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	return c.NoContent(204)
}

func (a *API) updateGroup(c echo.Context) error {
	log := ctxutil.LogFromContext(c)

	id := c.Param("id")

	var ugr autoscale.UpdateGroupRequest
	if err := c.Bind(&ugr); err != nil {
		return err
	}

	g, err := a.repo.GetGroup(a.ctx, id)
	if err != nil {
		log.WithError(err).Error("can't retrieve group")
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	err = a.repo.SaveGroup(a.ctx, g)

	if err != nil {
		log.WithError(err).Error("update group")
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	return c.JSON(http.StatusOK, g)
}

func (a *API) userConfig(c echo.Context) error {
	log := ctxutil.LogFromContext(c)

	client := autoscale.DOClientFactory()
	uc, err := autoscale.NewUserConfig(c, client)
	if err != nil {
		log.WithError(err).Error("retrieve user config")
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	j, err := jsonapi.Marshal(uc)
	if err != nil {
		log.WithError(err).Error("unable to marshal user config")
	}

	return c.JSONBlob(http.StatusOK, j)
}
