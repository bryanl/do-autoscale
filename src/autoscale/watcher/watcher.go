package watcher

import (
	"autoscale"
	"fmt"
	"pkg/do"
	"pkg/util/rand"
	"strconv"
	"sync"

	"golang.org/x/oauth2"

	"github.com/Sirupsen/logrus"
	"github.com/digitalocean/godo"
)

// doClient is our interface to digitalocean
type doClient struct {
	TagsService     do.TagsService
	DropletsService do.DropletsService
}

// dropletConfig is a configurtion for building a droplet.
type dropletConfig struct {
	doc   *doClient
	wg    *sync.WaitGroup
	log   *logrus.Entry
	group autoscale.Group
	tmpl  autoscale.Template
	tag   string
}

// Watcher watches groups.
type Watcher struct {
	repo     autoscale.Repository
	log      *logrus.Entry
	doClient *doClient
}

// New creates an instance of Watcher.
func New(pat string, repo autoscale.Repository) *Watcher {
	godoClient := createClient(pat)

	dc := &doClient{
		DropletsService: do.NewDropletsService(godoClient),
		TagsService:     do.NewTagsService(godoClient),
	}

	return &Watcher{
		repo:     repo,
		log:      logrus.WithField("action", "watcher"),
		doClient: dc,
	}
}

// Watch starts the watching process.
func (w *Watcher) Watch() {
	log := w.log

	groups, err := w.repo.ListGroups()
	if err != nil {
		log.WithError(err).Error("list groups")
	}

	for _, group := range groups {
		log := log.WithField("group-name", group.Name)
		log.Info("watching group")

		w.Check(group, log)
	}
}

// Check group to make sure it is at capacity.
func (w *Watcher) Check(g autoscale.Group, log *logrus.Entry) {
	checkMinStatus(g, log, w.doClient, w.repo)
}

func checkMinStatus(g autoscale.Group, log *logrus.Entry, doc *doClient, repo autoscale.Repository) error {

	tag := fmt.Sprintf("do:as:%s", g.ID)
	if err := verifyTag(tag, doc, log); err != nil {
		return err
	}

	droplets, err := doc.DropletsService.ListByTag(tag)
	if err != nil {
		return err
	}

	dCount := len(droplets)
	if dCount < g.BaseSize {
		log.WithFields(logrus.Fields{
			"wanted-droplets": g.BaseSize,
			"actual-droplets": dCount,
		}).Info("group is missing resources")

		var wg sync.WaitGroup
		n := g.BaseSize - dCount
		wg.Add(n)

		tmpl, err := repo.GetTemplate(g.TemplateName)
		if err != nil {
			return err
		}

		dc := dropletConfig{
			doc:   doc,
			wg:    &wg,
			log:   log,
			group: g,
			tmpl:  tmpl,
			tag:   tag,
		}

		for i := 0; i < n; i++ {
			go bootDroplet(&dc)
		}

		log.Info("waiting droplets to be created")
		wg.Wait()
		log.Info("droplets have been created")
	} else {
		log.Info("autoscale group exists")
	}

	return nil
}

func bootDroplet(dc *dropletConfig) {
	defer dc.wg.Done()
	id := rand.String(5)

	name := fmt.Sprintf("%s-%s", dc.group.BaseName, id)
	log := dc.log.WithFields(logrus.Fields{
		"droplet-name": name,
	})

	keys := []godo.DropletCreateSSHKey{}
	for _, k := range dc.tmpl.SSHKeys {
		str, _ := strconv.Atoi(k)
		dcs := godo.DropletCreateSSHKey{ID: str}
		keys = append(keys, dcs)
	}

	dcr := godo.DropletCreateRequest{
		Name:    name,
		Region:  dc.tmpl.Region,
		Size:    dc.tmpl.Size,
		Image:   godo.DropletCreateImage{Slug: dc.tmpl.Image},
		SSHKeys: keys,
	}

	log.Info("creating droplet")

	droplet, err := dc.doc.DropletsService.Create(&dcr, true)
	if err != nil {
		log.WithError(err).Error("unable to create droplet")
	}

	log.Info("created droplet")

	trr := &godo.TagResourcesRequest{
		Resources: []godo.Resource{
			{ID: strconv.Itoa(droplet.ID), Type: godo.DropletResourceType},
		},
	}

	logWithTag := log.WithFields(logrus.Fields{
		"tag-name":   dc.tag,
		"droplet-id": droplet.ID,
	})

	logWithTag.Info("tagging droplet")
	if err := dc.doc.TagsService.TagResources(dc.tag, trr); err != nil {
		logWithTag.WithError(err).Error("unable to tag droplet")
	}
}

func verifyTag(tag string, doc *doClient, log *logrus.Entry) error {
	log.WithField("tag", tag).Info("checking if tag exists")
	tags, err := doc.TagsService.List()
	if err != nil {
		log.WithError(err).Error("unable to list tags")
		return err
	}

	tagFound := false
	for _, t := range tags {
		if t.Name == tag {
			tagFound = true
			break
		}
	}

	if !tagFound {
		log.WithField("tag", tag).Info("creating tag")
		tcr := godo.TagCreateRequest{Name: tag}
		if _, err := doc.TagsService.Create(&tcr); err != nil {
			log.WithError(err).WithField("tag", tag).Error("could not create tag")
			return err
		}
	}

	return nil
}

type tokenSource struct {
	AccessToken string
}

// Token creates a token
func (t *tokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

func createClient(pat string) *godo.Client {
	ts := &tokenSource{
		AccessToken: pat,
	}

	oc := oauth2.NewClient(oauth2.NoContext, ts)
	return godo.NewClient(oc)
}
