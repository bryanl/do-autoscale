package autoscale

import (
	"pkg/doclient"

	"github.com/manyminds/api2go/jsonapi"
	"golang.org/x/net/context"
)

// GroupConfig are the options available for creating a group.
type GroupConfig struct {
	ID        string     `json:"-"`
	Policies  []string   `json:"policies"`
	Metrics   []string   `json:"metrics"`
	Templates []Template `json:"-"`
}

var _ jsonapi.MarshalIncludedRelations = (*GroupConfig)(nil)

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

// GetID returns the user id these policies are valid for.
func (gc *GroupConfig) GetID() string {
	return gc.ID
}

func (gc *GroupConfig) GetReferences() []jsonapi.Reference {
	return []jsonapi.Reference{
		{
			Type: "templates",
			Name: "templates",
		},
	}
}

func (gc *GroupConfig) GetReferencedIDs() []jsonapi.ReferenceID {
	result := []jsonapi.ReferenceID{}
	for _, t := range gc.Templates {
		result = append(result, jsonapi.ReferenceID{
			ID:   t.ID,
			Type: "templates",
			Name: "templates",
		})
	}

	return result
}

func (gc *GroupConfig) GetReferencedStructs() []jsonapi.MarshalIdentifier {
	result := []jsonapi.MarshalIdentifier{}
	for key := range gc.Templates {
		result = append(result, &gc.Templates[key])
	}

	return result
}
