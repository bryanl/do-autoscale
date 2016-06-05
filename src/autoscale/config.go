package autoscale

import (
	"crypto/md5"
	"fmt"
	"io"
	"pkg/doclient"
	"regexp"
	"time"

	"github.com/Sirupsen/logrus"
)

var (
	ErrActionTimedOut = fmt.Errorf("action timed out")
	ErrDisabledGroup  = fmt.Errorf("group is disabled")

	ScheduleReenqueueTimeout = 10 * time.Second

	SchedulerActionTimeout = 60 * time.Minute

	DefaultGroupCheckTimeout = 5 * time.Second

	// ResourceManagerFactory creates a ResourceManager given a group.
	ResourceManagerFactory ResourceManagerFactoryFn = func(g *Group) (ResourceManager, error) {
		doClient := DOClientFactory()
		tag := fmt.Sprintf("do:as:%s", g.Name)

		h := md5.New()
		io.WriteString(h, tag)
		hash := fmt.Sprintf("%x", h.Sum(nil))

		newTag := hash[0:8]

		log := logrus.WithField("group-name", g.Name)
		return NewDropletResource(doClient, newTag, log)
	}

	DOClientFactory = func() *doclient.Client {
		return doclient.New(DOAccessToken())
	}

	defaultValuePolicy = valuePolicyData{
		MinSize:        1,
		MaxSize:        10,
		ScaleUpValue:   0.8,
		ScaleUpBy:      2,
		ScaleDownValue: 0.2,
		ScaleDownBy:    1,
		WarmUpDuration: 10 * time.Second,
	}

	nameRe = regexp.MustCompile(`^\w[A-Za-z0-9\-]*$`)
)
