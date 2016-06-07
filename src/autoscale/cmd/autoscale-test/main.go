package main

import (
	"autoscale"
	"encoding/json"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/kelseyhightower/envconfig"

	"golang.org/x/net/context"
)

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
	var s Specification
	err := envconfig.Process("autoscale", &s)
	if err != nil {
		logrus.WithError(err).Fatal("unable to read environment")
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, autoscale.PrometheusURLContextKey, s.PrometheusURL)

	pl, err := autoscale.NewPrometheusLoad(ctx)
	if err != nil {
		panic(err)
	}

	out, err := pl.Values(ctx, "f5a3a208-c7b3-4bd1-a476-1be6e3e32817", autoscale.RangeWeek)
	if err != nil {
		panic(err)
	}

	j, err := json.Marshal(&out)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(j))
}
