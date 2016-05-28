package autoscale

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/manyminds/api2go/jsonapi"
)

// Template is a template that will be autoscaled.
type Template struct {
	ID       string  `json:"-" db:"id"`
	Name     string  `json:"name" db:"name"`
	Region   string  `json:"region" db:"region"`
	Size     string  `json:"size" db:"size"`
	Image    string  `json:"image" db:"image"`
	SSHKeys  SSHKeys `json:"ssh-keys" db:"ssh_keys"`
	UserData string  `json:"user-data" db:"user_data"`
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

type SSHKey struct {
	ID          int
	Fingerprint string
}

type SSHKeys []SSHKey

// Value converts SSH keys to JSON to be stored in the databases.
func (s SSHKeys) Value() (driver.Value, error) {
	return json.Marshal(s)
}

// Scan converts a DB value back into a SSHKeys.
func (s *SSHKeys) Scan(src interface{}) error {
	if src == nil {
		return nil
	}

	b := json.RawMessage(src.([]uint8))

	return json.Unmarshal(b, s)
}
