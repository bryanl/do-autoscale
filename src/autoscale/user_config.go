package autoscale

import (
	"pkg/do"
	"pkg/doclient"

	"gopkg.in/tomb.v2"
)

// UserConfig is the DO configuration for a user.
type UserConfig struct {
	Regions do.Regions `json:"regions"`
	Sizes   do.Sizes   `json:"sizes"`
	Keys    do.SSHKeys `json:"keys"`
}

// NewUserConfig creates an instance of UserConfig.
func NewUserConfig() (*UserConfig, error) {
	dc := doclient.New(DOAccessToken())
	uc := &UserConfig{}

	t := tomb.Tomb{}
	t.Go(func() error {
		return getRegions(dc, uc)
	})

	t.Go(func() error {
		return getSizes(dc, uc)
	})

	t.Go(func() error {
		return getKeys(dc, uc)
	})

	if err := t.Wait(); err != nil {
		return nil, err
	}

	return uc, nil
}

func getRegions(dc *doclient.Client, uc *UserConfig) error {
	regions, err := dc.RegionsService.List()
	if err != nil {
		return err
	}

	uc.Regions = regions
	return nil
}

func getSizes(dc *doclient.Client, uc *UserConfig) error {
	sizes, err := dc.SizesService.List()
	if err != nil {
		return err
	}

	uc.Sizes = sizes
	return nil
}

func getKeys(dc *doclient.Client, uc *UserConfig) error {
	keys, err := dc.KeysService.List()
	if err != nil {
		return err
	}

	uc.Keys = keys
	return nil
}
