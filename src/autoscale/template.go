package autoscale

import (
	"database/sql/driver"
	"encoding/json"
	"strconv"
)

// Template is a template that will be autoscaled.
type Template struct {
	ID       string  `json:"id" db:"id"`
	Name     string  `json:"name" db:"name"`
	Region   string  `json:"region" db:"region"`
	Size     string  `json:"size" db:"size"`
	Image    string  `json:"image" db:"image"`
	SSHKeys  SSHKeys `json:"sshKeys" db:"ssh_keys"`
	UserData string  `json:"userData" db:"user_data"`
}

// IsValid returns if the template is valid or not.
func (t *Template) IsValid() bool {
	if !nameRe.MatchString(t.Name) {
		return false
	}

	return true
}

func (t *Template) GetName() string {
	return "template"
}

// SSHKey is a DO ssh key.
type SSHKey struct {
	ID          int    `json:"id"`
	Fingerprint string `json:"fingerprint"`
}

func (s *SSHKey) GetID() string {
	return strconv.Itoa(s.ID)
}

func (s *SSHKey) SetID(id string) {
	i, _ := strconv.Atoi(id)
	s.ID = i
}

// SSHKeys is a slice of DO ssh keys.
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
