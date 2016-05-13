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
	CreateTemplate(tcr CreateTemplateRequest) (Template, error)
	GetTemplate(name string) (Template, error)
	ListTemplates() ([]Template, error)
	DeleteTemplate(name string) error

	CreateGroup(gcr CreateGroupRequest) (Group, error)
	GetGroup(name string) (Group, error)
	ListGroups() ([]Group, error)
	DeleteGroup(name string) error
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

func (r *pgRepo) CreateTemplate(tcr CreateTemplateRequest) (Template, error) {
	t := Template{
		Name:     tcr.Name,
		Region:   tcr.Region,
		Size:     tcr.Size,
		Image:    tcr.Image,
		SSHKeys:  tcr.SSHKeys,
		UserData: tcr.UserData,
	}

	if !t.IsValid() {
		return Template{}, errors.New(ValidationErr)
	}

	var id string

	tx, err := r.db.Beginx()
	if err != nil {
		return Template{}, err
	}

	err = sqlx.Get(tx, &id, sqlSaveTemplate,
		t.Name, t.Region, t.Size, t.Image, t.SSHKeys, t.UserData)
	if err != nil {
		tx.Rollback()
		return Template{}, err
	}

	err = tx.Commit()
	if err != nil {
		return Template{}, err
	}

	t.ID = id
	return t, nil
}

func (r *pgRepo) GetTemplate(name string) (Template, error) {
	var t Template
	if err := r.db.Get(&t, sqlGetTemplate, name); err != nil {
		return Template{}, err
	}

	return t, nil
}

func (r *pgRepo) ListTemplates() ([]Template, error) {
	ts := []Template{}
	if err := r.db.Select(&ts, sqlListTemplates); err != nil {
		return nil, err
	}

	return ts, nil
}

func (r *pgRepo) DeleteTemplate(name string) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	_, err = tx.Exec(sqlDeleteTemplate, name)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *pgRepo) CreateGroup(gcr CreateGroupRequest) (Group, error) {
	g := Group{
		Name:         gcr.Name,
		BaseName:     gcr.BaseName,
		BaseSize:     gcr.BaseSize,
		MetricType:   gcr.MetricType,
		TemplateName: gcr.TemplateName,
	}

	if !g.IsValid() {
		return Group{}, errors.New(ValidationErr)
	}

	var id string

	tx, err := r.db.Beginx()
	if err != nil {
		return Group{}, err
	}

	err = sqlx.Get(tx, &id, sqlCreateGroup,
		g.Name, g.BaseName, g.BaseSize, g.MetricType, g.TemplateName, g.ScaleGroup)
	if err != nil {
		tx.Rollback()
		return Group{}, err
	}

	err = tx.Commit()
	if err != nil {
		return Group{}, err
	}

	g.ID = id

	return g, nil
}

func (r *pgRepo) GetGroup(name string) (Group, error) {
	var g Group
	if err := r.db.Get(&g, sqlGetGroup, name); err != nil {
		return Group{}, err
	}

	return g, nil
}

func (r *pgRepo) ListGroups() ([]Group, error) {
	ts := []Group{}
	if err := r.db.Select(&ts, sqlListGroups); err != nil {
		return nil, err
	}

	return ts, nil
}

func (r *pgRepo) DeleteGroup(name string) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	_, err = tx.Exec(sqlDeleteGroup, name)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

var (
	sqlSaveTemplate = `
  INSERT into templates
  (name, region, size, image, ssh_keys, user_data)
  VALUES ($1, $2, $3, $4, $5, $6)
  RETURNING id`

	sqlGetTemplate = `
  SELECT * from templates where name=$1`

	sqlListTemplates = `
  SELECT * from templates`

	sqlDeleteTemplate = `
  DELETE from templates WHERE id = $1`

	sqlCreateGroup = `
  INSERT into groups
  (name, base_name, base_size, metric_type, template_name, rules)
  VALUES ($1, $2, $3, $4, $5, $6)
  RETURNING id`

	sqlGetGroup = `
  SELECT * from groups where name=$1`

	sqlListGroups = `
  SELECT * from groups`

	sqlDeleteGroup = `
  DELETE from groups WHERE id = $1`
)
