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
	// ErrActionTimedOut is a timeout error.
	ErrActionTimedOut = fmt.Errorf("action timed out")

	// ErrDisabledGroup is return if the group is disabled.
	ErrDisabledGroup = fmt.Errorf("group is disabled")

	// ScheduleReenqueueTimeout is how long to wait when reenqueuing a check.
	ScheduleReenqueueTimeout = 10 * time.Second

	// SchedulerActionTimeout is time it takes a schedule action to timeout.
	SchedulerActionTimeout = 60 * time.Minute

	// DefaultGroupCheckTimeout is how often new groups are checked.
	DefaultGroupCheckTimeout = 5 * time.Second

	tagNameFn = defaultTagName

	// ResourceManagerFactory creates a ResourceManager given a group.
	ResourceManagerFactory ResourceManagerFactoryFn = func(g *Group) (ResourceManager, error) {
		doClient := DOClientFactory()
		tagName := tagNameFn(g.Name)
		log := logrus.WithField("group-id", g.ID)
		return NewDropletResource(doClient, tagName, log)
	}

	// DOClientFactory createsa  do client.
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

	prometheusAgentPort = 9100

	// BaseTag is the tag applied to all droplets created by autoscale.
	BaseTag = "autoscale"
)

func defaultTagName(groupName string) string {
	h := md5.New()
	io.WriteString(h, groupName)
	hash := fmt.Sprintf("%x", h.Sum(nil))

	tag := fmt.Sprintf("as:%s", hash[0:8])
	return tag
}
