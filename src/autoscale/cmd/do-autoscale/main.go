package main

import (
	"autoscale"
	"autoscale/api"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/kelseyhightower/envconfig"
)

// Specification describes our expected environment.
type Specification struct {
	DBUser     string `envconfig:"db_user"`
	DBPassword string `envconfig:"db_password"`
	DBAddr     string `envconfig:"db_addr"`
	DBName     string `envconfig:"db_name"`
	HTTPAddr   string `envconfig:"http_addr" default:"localhost:8888"`
}

func main() {
	var s Specification
	err := envconfig.Process("autoscale", &s)
	if err != nil {
		log.WithError(err).Fatal("unable to read environment")
	}

	db, err := autoscale.NewDB(s.DBUser, s.DBPassword, s.DBAddr, s.DBName)
	if err != nil {
		log.WithError(err).Fatal("unable to create database connection")
	}

	repo, err := autoscale.NewRepository(db)
	if err != nil {
		log.WithError(err).Fatal("unable to setup data repository")
	}

	a := api.New(repo)
	http.Handle("/", a.Mux)

	log.WithFields(log.Fields{
		"http-addr": s.HTTPAddr,
	}).Info("created http server")
	log.Fatal(http.ListenAndServe(s.HTTPAddr, nil))
}
