package api

import (
	"encoding/json"
	"net/http"

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
	r.HandleFunc("/templates/{id}", a.getTemplate).Methods("GET")
	r.HandleFunc("/templates", a.createTemplate).Methods("POST")
	r.HandleFunc("/templates/{id}", a.deleteTemplate).Methods("DELETE")
	r.HandleFunc("/groups", a.listGroups).Methods("GET")
	r.HandleFunc("/groups/{id}", a.getGroup).Methods("GET")
	r.HandleFunc("/groups", a.createGroup).Methods("POST")
	r.HandleFunc("/groups/{id}", a.deleteGroup).Methods("DELETE")
	r.HandleFunc("/groups/{id}", a.updateGroup).Methods("PUT")

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
	id := vars["id"]

	tmpl, err := a.repo.GetTemplate(id)
	if err != nil {
		log.WithError(err).Error("retrieve template")
		writeError(w, "unable to retrieve template", http.StatusNotFound)
		return
	}

	_ = json.NewEncoder(w).Encode(tmpl)
}

func (a *API) createTemplate(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var ctr autoscale.CreateTemplateRequest
	err := json.NewDecoder(r.Body).Decode(&ctr)
	if err != nil {
		log.WithError(err).Error("create template")
		writeError(w, "invalid create template request", http.StatusBadRequest)
		return
	}

	tmpl, err := a.repo.CreateTemplate(ctr)
	if err != nil {
		writeError(w, "unable to create template", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(&tmpl)
}

func (a *API) deleteTemplate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := a.repo.DeleteTemplate(id)
	if err != nil {
		log.WithError(err).Error("delete template")
		writeError(w, "unable to delete template", http.StatusBadRequest)
		return
	}

	w.WriteHeader(204)
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
		log.WithError(err).Error("get group")
		writeError(w, "unable to retrieve group", http.StatusNotFound)
		return
	}

	_ = json.NewEncoder(w).Encode(group)
}

func (a *API) createGroup(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var cgr autoscale.CreateGroupRequest
	err := json.NewDecoder(r.Body).Decode(&cgr)
	if err != nil {
		writeError(w, "invalid create group request", http.StatusBadRequest)
		return
	}

	g, err := a.repo.CreateGroup(cgr)
	if err != nil {
		log.WithError(err).Error("create group")
		writeError(w, "unable to create group", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(&g)
}

func (a *API) deleteGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := a.repo.DeleteGroup(id)
	if err != nil {
		log.WithError(err).Error("delete group")
		writeError(w, "unable to delete group", http.StatusBadRequest)
		return
	}

	w.WriteHeader(204)
}

func (a *API) updateGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	defer r.Body.Close()

	var ugr autoscale.UpdateGroupRequest
	err := json.NewDecoder(r.Body).Decode(&ugr)
	if err != nil {
		writeError(w, "invalid update group request", http.StatusBadRequest)
		return
	}

	g, err := a.repo.GetGroup(id)
	if err != nil {
		writeError(w, "invalid update group request", http.StatusNotFound)
		return
	}

	g.BaseSize = ugr.BaseSize
	err = a.repo.SaveGroup(g)
	if err != nil {
		log.WithError(err).Error("update group")
		writeError(w, "unable to update group", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(&g)
}
