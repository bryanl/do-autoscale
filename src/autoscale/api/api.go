package api

import (
	"encoding/json"
	"net/http"
	"pkg/ctxutil"

	"golang.org/x/net/context"

	"autoscale"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine"
	"github.com/labstack/echo/engine/standard"
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

	e.GET("/templates/:id", a.getTemplate)
	e.GET("/templates", a.listTemplates)
	e.POST("/templates", a.createTemplate)
	e.DELETE("/templates/:id", a.deleteTemplate)
	e.GET("/groups", a.listGroups)
	e.GET("/groups/:id", a.getGroup)
	e.POST("/groups", a.createGroup)
	e.DELETE("/groups/:id", a.deleteGroup)
	e.PUT("/groups/:id", a.updateGroup)

	return a
}

func (a *API) listTemplates(c echo.Context) error {
	log := ctxutil.LogFromContext(a.ctx)
	tmpls, err := a.repo.ListTemplates(a.ctx)
	if err != nil {
		log.WithError(err).Error("list templates")

		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, tmpls)
}

func (a *API) getTemplate(c echo.Context) error {
	log := ctxutil.LogFromContext(a.ctx)

	id := c.Param("id")

	tmpl, err := a.repo.GetTemplate(a.ctx, id)
	if err != nil {
		log.WithError(err).Error("retrieve template")
		return echo.ErrNotFound
	}

	return c.JSON(http.StatusOK, tmpl)
}

func (a *API) createTemplate(c echo.Context) error {
	log := ctxutil.LogFromContext(a.ctx)

	var ctr autoscale.CreateTemplateRequest
	if err := c.Bind(&ctr); err != nil {
		return err
	}

	tmpl, err := a.repo.CreateTemplate(a.ctx, ctr)
	if err != nil {
		log.WithError(err).Error("delete template")

		return echo.NewHTTPError(http.StatusBadRequest)
	}

	return c.JSON(http.StatusCreated, tmpl)
}

func (a *API) deleteTemplate(c echo.Context) error {
	log := ctxutil.LogFromContext(a.ctx)

	id := c.Param("id")

	err := a.repo.DeleteTemplate(a.ctx, id)
	if err != nil {
		log.WithError(err).Error("delete template")
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	return c.NoContent(204)
}

func (a *API) listGroups(c echo.Context) error {
	log := ctxutil.LogFromContext(a.ctx)

	groups, err := a.repo.ListGroups(a.ctx)
	if err != nil {
		log.WithError(err).Error("list groups")
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, groups)
}

func (a *API) getGroup(c echo.Context) error {
	log := ctxutil.LogFromContext(a.ctx)

	id := c.Param("id")

	group, err := a.repo.GetGroup(a.ctx, id)
	if err != nil {
		log.WithError(err).Error("get group")
		return echo.ErrNotFound
	}

	return c.JSON(http.StatusOK, group)
}

func (a *API) createGroup(c echo.Context) error {
	log := ctxutil.LogFromContext(a.ctx)

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
	log := ctxutil.LogFromContext(a.ctx)

	id := c.Param("id")

	err := a.repo.DeleteGroup(a.ctx, id)
	if err != nil {
		log.WithError(err).Error("delete group")
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	return c.NoContent(204)
}

func (a *API) updateGroup(c echo.Context) error {
	log := ctxutil.LogFromContext(a.ctx)

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
