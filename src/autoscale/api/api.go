package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"autoscale"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
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

type createTemplateRequest struct {
	Name     string   `json:"name"`
	Region   string   `json:"region"`
	Size     string   `json:"size"`
	Image    string   `json:"image"`
	SSHKeys  []string `json:"ssh_keys"`
	UserData string   `json:"user_data"`
}

type createGroupRequest struct {
	BaseName   string `json:"base_name"`
	BaseSize   int    `json:"base_size"`
	MetricType string `json:"metric_type"`
	TemplateID int    `json:"template_id"`
}

// API is the autoscale API.
type API struct {
	Mux  *mux.Router
	repo autoscale.Repository
}

// New creates an instance of API.
func New(repo autoscale.Repository) *API {
	r := mux.NewRouter()

	a := &API{
		Mux:  r,
		repo: repo,
	}

	r.HandleFunc("/templates", a.listTemplates).Methods("GET")
	r.HandleFunc("/templates/{id:[0-9]+}", a.getTemplate).Methods("GET")
	r.HandleFunc("/templates", a.createTemplate).Methods("POST")
	r.HandleFunc("/groups", a.listGroups).Methods("GET")
	r.HandleFunc("/groups/{id}", a.getGroup).Methods("GET")
	r.HandleFunc("/groups", a.createGroup).Methods("POST")

	return a
}

func (a *API) listTemplates(w http.ResponseWriter, r *http.Request) {
	tmpls, err := a.repo.ListTemplates()
	if err != nil {
		log.WithError(err).Error("list templates")
		writeError(w, "unable to list templates", http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(tmpls)
}

func (a *API) getTemplate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		writeError(w, "invalid id", http.StatusBadRequest)
		return
	}

	tmpl, err := a.repo.GetTemplate(id)
	if err != nil {
		writeError(w, "unable to retrieve template", http.StatusNotFound)
		return
	}

	_ = json.NewEncoder(w).Encode(tmpl)
}

func (a *API) createTemplate(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var ctr createTemplateRequest
	err := json.NewDecoder(r.Body).Decode(&ctr)
	if err != nil {
		writeError(w, "invalid create template request", http.StatusBadRequest)
		return
	}

	tmpl := autoscale.Template{
		Name:       ctr.Name,
		Region:     ctr.Region,
		Size:       ctr.Size,
		Image:      ctr.Image,
		RawSSHKeys: strings.Join(ctr.SSHKeys, ","),
		UserData:   ctr.UserData,
	}

	id, err := a.repo.CreateTemplate(&tmpl)
	if err != nil {
		writeError(w, "unable to create template", http.StatusBadRequest)
		return
	}

	tmpl.ID = id

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(tmpl)
}

func (a *API) listGroups(w http.ResponseWriter, r *http.Request) {
	groups, err := a.repo.ListGroups()
	if err != nil {
		log.WithError(err).Error("list groups")
		writeError(w, "unable to list groups", http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(groups)
}

func (a *API) getGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	group, err := a.repo.GetGroup(id)
	if err != nil {
		writeError(w, "unable to retrieve group", http.StatusNotFound)
		return
	}

	_ = json.NewEncoder(w).Encode(group)
}

func (a *API) createGroup(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var cgr createGroupRequest
	err := json.NewDecoder(r.Body).Decode(&cgr)
	if err != nil {
		writeError(w, "invalid create group request", http.StatusBadRequest)
		return
	}

	group := autoscale.Group{
		BaseName:   cgr.BaseName,
		BaseSize:   cgr.BaseSize,
		MetricType: cgr.MetricType,
		TemplateID: cgr.TemplateID,
	}

	id, err := a.repo.CreateGroup(&group)
	if err != nil {
		writeError(w, "unable to create group", http.StatusBadRequest)
		return
	}

	group.ID = id

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(group)
}
