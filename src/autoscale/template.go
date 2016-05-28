package autoscale

import "github.com/manyminds/api2go/jsonapi"

// Template is a template that will be autoscaled.
type Template struct {
	ID       string      `json:"-" db:"id"`
	Name     string      `json:"name" db:"name"`
	Region   string      `json:"region" db:"region"`
	Size     string      `json:"size" db:"size"`
	Image    string      `json:"image" db:"image"`
	SSHKeys  StringSlice `json:"ssh_keys" db:"ssh_keys"`
	UserData string      `json:"user_data" db:"user_data"`
}

var _ jsonapi.MarshalIdentifier = (*Template)(nil)
var _ jsonapi.UnmarshalIdentifier = (*Template)(nil)

// IsValid returns if the template is valid or not.
func (t *Template) IsValid() bool {
	if !nameRe.MatchString(t.Name) {
		return false
	}

	return true
}

// GetID gets the ID for a template.
func (t *Template) GetID() string {
	return t.ID
}

// SetID sets the ID for a template.
func (t *Template) SetID(id string) error {
	t.ID = id
	return nil
}
