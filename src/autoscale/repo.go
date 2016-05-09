package autoscale

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// Repository maps data to an entity models.
type Repository interface {
	SaveTemplate(t *Template) (int, error)
	GetTemplate(id int) (*Template, error)
	ListTemplates() ([]Template, error)
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

func (r *pgRepo) SaveTemplate(t *Template) (int, error) {
	var id int

	err := r.db.Get(&id, sqlSaveTemplate,
		t.Region, t.Size, t.Image, t.RawSSHKeys, t.UserData)
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

var (
	sqlSaveTemplate = `
  INSERT into templates
  (region, size, image, ssh_keys, user_data)
  VALUES ($1, $2, $3, $4, $5)
  RETURNING id`

	sqlGetTemplate = `
  SELECT * from templates where id=$1`

	sqlListTemplates = `
  SELECT * from templates`
)
