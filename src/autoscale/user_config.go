package autoscale

import (
	"pkg/do"
	"pkg/doclient"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/manyminds/api2go/jsonapi"

	"golang.org/x/net/context"
)

// UserConfigResponse is a response suitable for returning in a REST API.
type UserConfigResponse struct {
	UserConfig *UserConfig `json:"user_config"`
}

// UserConfig is the DO configuration for a user.
type UserConfig struct {
	ID      string     `json:"id"`
	Regions do.Regions `json:"regions"`
	Sizes   do.Sizes   `json:"sizes"`
	Keys    do.SSHKeys `json:"keys"`
}

var _ jsonapi.MarshalIdentifier = (*UserConfig)(nil)

// NewUserConfig creates an instance of UserConfig.
func NewUserConfig(ctx context.Context, dc *doclient.Client) (*UserConfig, error) {
	logger := logrus.New()
	log := logrus.NewEntry(logger)
	uc := &UserConfig{}
	errs := []error{}

	wg := sync.WaitGroup{}
	wg.Add(4)

	go func() {
		defer wg.Done()
		getRegions(log, dc, uc)
	}()

	go func() {
		defer wg.Done()
		getSizes(log, dc, uc)
	}()

	go func() {
		defer wg.Done()
		getKeys(log, dc, uc)
	}()

	go func() {
		defer wg.Done()
		getID(log, dc, uc)
	}()

	wg.Wait()

	if len(errs) > 0 {
		return nil, errs[0]
	}

	return uc, nil
}

// GetID returns the ID for the UserConfig. Uses the DO account's UUID.
func (uc *UserConfig) GetID() string {
	return uc.ID
}

func getRegions(log *logrus.Entry, dc *doclient.Client, uc *UserConfig) error {
	regions, err := dc.RegionsService.List()
	if err != nil {
		return err
	}

	uc.Regions = regions
	return nil
}

func getSizes(log *logrus.Entry, dc *doclient.Client, uc *UserConfig) error {
	sizes, err := dc.SizesService.List()
	if err != nil {
		return err
	}

	uc.Sizes = sizes
	return nil
}

func getKeys(log *logrus.Entry, dc *doclient.Client, uc *UserConfig) error {
	keys, err := dc.KeysService.List()
	if err != nil {
		return err
	}

	uc.Keys = keys
	return nil
}

func getID(log *logrus.Entry, dc *doclient.Client, uc *UserConfig) error {
	a, err := dc.AccountsService.Get()
	if err != nil {
		return err
	}

	uc.ID = a.UUID
	return nil
}
