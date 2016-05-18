package main

import (
	"autoscale"
	"autoscale/api"
	"autoscale/metrics"
	"autoscale/watcher"

	"golang.org/x/net/context"

	"math/rand"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/kelseyhightower/envconfig"
)

// Specification describes our expected environment.
type Specification struct {
	Env                    string `envconfig:"env" default:"development"`
	DBUser                 string `envconfig:"db_user" required:"true"`
	DBPassword             string `envconfig:"db_password" required:"true"`
	DBAddr                 string `envconfig:"db_addr" required:"true"`
	DBName                 string `envconfig:"db_name" required:"true"`
	HTTPAddr               string `envconfig:"http_addr" default:"localhost:8888"`
	AccessToken            string `envconfig:"access_token" required:"true"`
	UseFileStats           bool   `envconfig:"use_file_stats" default:"false"`
	FileStatDir            string `envconfig:"file_stat_dir"`
	PrometheusConfigDir    string `envconfig:"prometheus_config_dir"`
	RegisterOfflineMetrics bool   `envconfig:"register_offline_metrics" default:"false"`
	RegisterDefaultMetrics bool   `envconfig:"register_default_metrics" default:"true"`
	UseMemoryResources     bool   `envconfig:"use_memory_resources" default:"false"`
}

func main() {
	log := logrus.New()
	rand.Seed(time.Now().UnixNano())

	var s Specification
	err := envconfig.Process("autoscale", &s)
	if err != nil {
		log.WithError(err).Fatal("unable to read environment")
	}

	ctx := context.WithValue(context.Background(), "log", log.WithField("env", s.Env))

	if s.UseMemoryResources {
		rm := autoscale.NewLocalResource(ctx)
		log.Info("using memory resources")
		autoscale.ResourceManagerFactory = func(g *autoscale.Group) (autoscale.ResourceManager, error) {
			return rm, nil
		}
	}

	if s.RegisterDefaultMetrics && s.RegisterOfflineMetrics {
		log.Fatal("can't specify offline and default metrics at the same time")
	}

	if s.RegisterDefaultMetrics {
		ctx = context.WithValue(ctx, "prometheusConfigDir", s.PrometheusConfigDir)
		metrics.RegisterDefaultMetrics(ctx)
	}

	if s.RegisterOfflineMetrics {
		metrics.RegisterOfflineMetrics(ctx)
	}

	autoscale.DOAccessToken = func() string {
		return s.AccessToken
	}

	db, err := autoscale.NewDB(s.DBUser, s.DBPassword, s.DBAddr, s.DBName)
	if err != nil {
		log.WithError(err).Fatal("unable to create database connection")
	}

	repo, err := autoscale.NewRepository(db)
	if err != nil {
		log.WithError(err).Fatal("unable to setup data repository")
	}

	watcher := watcher.New(repo)
	go func() {
		if _, err := watcher.Watch(); err != nil {
			log.WithError(err).Fatal("unable to start watcher")
		}
	}()

	a := api.New(repo)
	http.Handle("/", a.Mux)

	log.WithFields(logrus.Fields{
		"http-addr": s.HTTPAddr,
	}).Info("created http server")
	log.Fatal(http.ListenAndServe(s.HTTPAddr, nil))
}
