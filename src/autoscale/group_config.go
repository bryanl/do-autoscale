package autoscale

import (
	"pkg/doclient"

	"golang.org/x/net/context"
)

// GroupConfig are the options available for creating a group.
type GroupConfig struct {
	ID        string     `json:"id"`
	Policies  []string   `json:"policies"`
	Metrics   []string   `json:"metrics"`
	Templates []Template `json:"templates"`
}

// NewGroupConfig creates an instance of GroupConfig.
func NewGroupConfig(ctx context.Context, dc *doclient.Client, repo Repository) (*GroupConfig, error) {
	a, err := dc.AccountsService.Get()
	if err != nil {
		return nil, err
	}

	tmpls, err := repo.ListTemplates(ctx)
	if err != nil {
		return nil, err
	}

	return &GroupConfig{
		ID:        a.UUID,
		Policies:  []string{"value"},
		Metrics:   []string{"load"},
		Templates: tmpls,
	}, nil
}
