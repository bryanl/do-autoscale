package main

import (
	"autoscale"
	"autoscale/api"
	"pkg/ctxutil"

	"golang.org/x/net/context"

	"math/rand"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/tylerb/graceful.v1"
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
	PrometheusURL          string `envconfig:"prometheus_url" required:"true"`
	RegisterOfflineMetrics bool   `envconfig:"register_offline_metrics" default:"false"`
	RegisterDefaultMetrics bool   `envconfig:"register_default_metrics" default:"true"`
	UseMemoryResources     bool   `envconfig:"use_memory_resources" default:"false"`
}

func main() {
	rand.Seed(time.Now().UnixNano())

	ctx, s, log := initContext()

	if s.RegisterDefaultMetrics {
		ctx = context.WithValue(ctx, autoscale.PrometheusConfigDirContextKey, s.PrometheusConfigDir)
		autoscale.RegisterDefaultMetrics(ctx)
	}

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

	if s.RegisterOfflineMetrics {
		autoscale.RegisterOfflineMetrics(ctx)
	}

	autoscale.DOAccessToken = func() string {
		return s.AccessToken
	}

	repo, err := initRepository(ctx, s, log)
	if err != nil {
		log.WithError(err).Fatal("unable to initialize repository")
	}

	watcher, err := initWatcher(ctx, repo, log)
	if err != nil {
		log.WithError(err).Fatal("unable to initialize watcher")
	}

	a := api.New(ctx, repo)

	log.WithFields(logrus.Fields{
		"http-addr": s.HTTPAddr,
	}).Info("starting http server")

	if err := graceful.RunWithErr(s.HTTPAddr, 5*time.Second, a.Mux); err != nil {
		log.WithError(err).Error("http server did not exit successfully")
	}

	log.Info("shutting down")
	watcher.Stop()
	if err := repo.Close(); err != nil {
		log.WithError(err).Error("repository did not close successfully")
	}
}

func initContext() (context.Context, Specification, *logrus.Entry) {
	logger := logrus.New()
	var s Specification
	err := envconfig.Process("autoscale", &s)
	if err != nil {
		logger.WithError(err).Fatal("unable to read environment")
	}

	log := logger.WithField("env", s.Env)
	ctx := context.WithValue(context.Background(), ctxutil.KeyLog, log)
	ctx = context.WithValue(ctx, autoscale.PrometheusURLContextKey, s.PrometheusURL)
	ctx = context.WithValue(ctx, ctxutil.KeyEnv, s.Env)
	ctx = context.WithValue(ctx, ctxutil.KeyDOToken, s.AccessToken)

	return ctx, s, log
}

func initRepository(ctx context.Context, s Specification, log *logrus.Entry) (autoscale.Repository, error) {
	db, err := autoscale.NewDB(ctx, s.DBUser, s.DBPassword, s.DBAddr, s.DBName)
	if err != nil {
		log.WithError(err).Error("unable to create database connection")
		return nil, err
	}

	repo, err := autoscale.NewRepository(db)
	if err != nil {
		log.WithError(err).Error("unable to setup data repository")
		return nil, err
	}

	return repo, nil
}

func initWatcher(ctx context.Context, repo autoscale.Repository, log *logrus.Entry) (*autoscale.Watcher, error) {
	watcher, err := autoscale.NewWatcher(ctx, repo)
	if err != nil {
		log.WithError(err).Error("unable to setup watcher")
		return nil, err
	}

	_, err = watcher.Watch()
	if err != nil {
		log.WithError(err).Error("unable to start watcher")
		return nil, err
	}

	return watcher, nil
}
