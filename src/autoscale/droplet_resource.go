package autoscale

import (
	"fmt"
	"pkg/doclient"
	"pkg/util/rand"
	"pkg/util/shuffle"
	"strconv"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/digitalocean/godo"
)

var (
	// AccessToken is the access token that will be used to interact with DigitalOcean.
	DOAccessToken func() string
)

// dropletConfig is a configurtion for building a droplet.
type dropletConfig struct {
	doc       *doclient.Client
	wg        *sync.WaitGroup
	log       *logrus.Entry
	groupName string
	template  Template
	tag       string
}

// DropletResource watches Droplets.
type DropletResource struct {
	doClient *doclient.Client
	tag      string
	log      *logrus.Entry
}

var _ Resource = (*DropletResource)(nil)

// NewDropletResource creates an instance of WatchedDropletResource.
func NewDropletResource(doClient *doclient.Client, tag string, log *logrus.Entry) (*DropletResource, error) {
	if err := verifyTag(tag, doClient, log); err != nil {
		return nil, err
	}

	return &DropletResource{
		doClient: doClient,
		tag:      tag,
		log:      log,
	}, nil
}

// Actual returns the amount of actual Droplets.
func (r *DropletResource) Actual() (int, error) {
	droplets, err := r.doClient.DropletsService.ListByTag(r.tag)
	if err != nil {
		return 0, err
	}

	return len(droplets), nil
}

// ScaleUp scales Droplet resources up.
func (r *DropletResource) ScaleUp(g Group, byN int, repo Repository) error {
	r.log.WithField("by-n", byN).Info("scaling up")

	var wg sync.WaitGroup
	wg.Add(byN)

	tmpl, err := repo.GetTemplate(g.TemplateName)
	if err != nil {
		return err
	}

	dc := dropletConfig{
		doc:       r.doClient,
		wg:        &wg,
		log:       r.log,
		groupName: g.BaseName,
		template:  tmpl,
		tag:       r.tag,
	}

	for i := 0; i < byN; i++ {
		go bootDroplet(&dc)
	}

	r.log.Info("waiting for droplets to be created")
	wg.Wait()
	r.log.Info("droplets have been created")

	return nil
}

// ScaleDown scales Droplet resources down.
func (r *DropletResource) ScaleDown(g Group, byN int, repo Repository) error {
	r.log.WithField("by-n", byN).Info("scaling down")
	droplets, err := r.doClient.DropletsService.ListByTag(r.tag)
	if err != nil {
		return err
	}

	ids := []int{}
	for _, d := range droplets {
		ids = append(ids, d.ID)
	}

	shuffle.Int(ids)
	for i := 0; i < byN; i++ {
		id := ids[i]
		r.log.WithField("droplet-id", id).Info("deleting droplet")
		if err := r.doClient.DropletsService.Delete(id); err != nil {
			r.log.WithError(err).WithField("droplet-id", id).Error("could not delete droplet")
			return err
		}
	}

	return nil
}

func bootDroplet(dc *dropletConfig) {
	defer dc.wg.Done()
	id := rand.String(5)

	name := fmt.Sprintf("%s-%s", dc.groupName, id)
	log := dc.log.WithFields(logrus.Fields{
		"droplet-name": name,
	})

	keys := []godo.DropletCreateSSHKey{}
	for _, k := range dc.template.SSHKeys {
		str, _ := strconv.Atoi(k)
		dcs := godo.DropletCreateSSHKey{ID: str}
		keys = append(keys, dcs)
	}

	dcr := godo.DropletCreateRequest{
		Name:    name,
		Region:  dc.template.Region,
		Size:    dc.template.Size,
		Image:   godo.DropletCreateImage{Slug: dc.template.Image},
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
		dc.doc.DropletsService.Delete(droplet.ID)
		logWithTag.WithError(err).Error("deleting droplet because it cannot be tagged")
	}
}

func verifyTag(tag string, doc *doclient.Client, log *logrus.Entry) error {
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
