package autoscale

import (
	"database/sql"

	"github.com/go-errors/errors"
	"github.com/jmoiron/sqlx"
)

var (
	// ValidationErr is returned when the model isn't valid.
	ValidationErr = errors.Errorf("is not valid")
)

// Repository maps data to an entity models.
type Repository interface {
	CreateTemplate(t *Template) (int, error)
	GetTemplate(id int) (*Template, error)
	ListTemplates() ([]Template, error)

	CreateGroup(t *Group) (string, error)
	GetGroup(id string) (*Group, error)
	ListGroups() ([]Group, error)
}

type pgRepo struct {
	db *sqlx.DB
}

var _ Repository = (*pgRepo)(nil)

// NewRepository creates a Repository backed with postgresql.
func NewRepository(db *sql.DB) (Repository, error) {
	repoDB := sqlx.NewDb(db, "postgres")
	return &pgRepo{
		db: repoDB,
	}, nil
}

func (r *pgRepo) CreateTemplate(t *Template) (int, error) {
	if !t.IsValid() {
		return 0, errors.New(ValidationErr)
	}

	var id int

	err := r.db.Get(&id, sqlSaveTemplate,
		t.Name, t.Region, t.Size, t.Image, t.RawSSHKeys, t.UserData)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *pgRepo) GetTemplate(id int) (*Template, error) {
	var t Template
	if err := r.db.Get(&t, sqlGetTemplate, id); err != nil {
		return nil, err
	}

	return &t, nil
}

func (r *pgRepo) ListTemplates() ([]Template, error) {
	ts := []Template{}
	if err := r.db.Select(&ts, sqlListTemplates); err != nil {
		return nil, err
	}

	return ts, nil
}

func (r *pgRepo) CreateGroup(g *Group) (string, error) {
	var id string

	err := r.db.Get(&id, sqlCreateGroup,
		g.BaseName, g.BaseSize, g.MetricType, g.TemplateID)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *pgRepo) GetGroup(id string) (*Group, error) {
	var g Group
	if err := r.db.Get(&g, sqlGetGroup, id); err != nil {
		return nil, err
	}

	return &g, nil
}

func (r *pgRepo) ListGroups() ([]Group, error) {
	ts := []Group{}
	if err := r.db.Select(&ts, sqlListGroups); err != nil {
		return nil, err
	}

	return ts, nil
}

var (
	sqlSaveTemplate = `
  INSERT into templates
  (name, region, size, image, ssh_keys, user_data)
  VALUES ($1, $2, $3, $4, $5, $6)
  RETURNING id`

	sqlGetTemplate = `
  SELECT * from templates where id=$1`

	sqlListTemplates = `
  SELECT * from templates`

	sqlCreateGroup = `
  INSERT into groups
  (base_name, base_size, metric_type, template_id)
  VALUES ($1, $2, $3, $4)
  RETURNING id`

	sqlGetGroup = `
  SELECT * from groups where id=$1`

	sqlListGroups = `
  SELECT * from groups`
)
