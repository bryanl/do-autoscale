package autoscale

import (
	"database/sql"

	"github.com/go-errors/errors"
	"github.com/jmoiron/sqlx"

	"golang.org/x/net/context"
)

var (
	// ValidationErr is returned when the model isn't valid.
	ValidationErr = errors.Errorf("is not valid")

	// ObjectMissingErr is returned with the requested object does not exist.
	ObjectMissingErr = errors.Errorf("object does not exist")
)

// Repository maps data to an entity models.
type Repository interface {
	CreateTemplate(ctx context.Context, t Template) (*Template, error)
	GetTemplate(ctx context.Context, name string) (*Template, error)
	ListTemplates(ctx context.Context) ([]Template, error)
	DeleteTemplate(ctx context.Context, name string) error

	CreateGroup(ctx context.Context, g Group) (*Group, error)
	GetGroup(ctx context.Context, name string) (*Group, error)
	ListGroups(ctx context.Context) ([]Group, error)
	DeleteGroup(ctx context.Context, name string) error
	SaveGroup(ctx context.Context, group Group) error

	Close() error
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

func (r *pgRepo) CreateTemplate(ctx context.Context, t Template) (*Template, error) {
	if !t.IsValid() {
		return nil, errors.New(ValidationErr)
	}

	var id string

	tx, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}

	err = sqlx.Get(tx, &id, sqlSaveTemplate,
		t.Name, t.Region, t.Size, t.Image, t.SSHKeys, t.UserData)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	t.ID = id
	return &t, nil
}

func (r *pgRepo) GetTemplate(ctx context.Context, id string) (*Template, error) {
	var t Template
	if err := r.db.Get(&t, sqlGetTemplate, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, ObjectMissingErr
		}

		return nil, err
	}

	return &t, nil
}

func (r *pgRepo) ListTemplates(ctx context.Context) ([]Template, error) {
	ts := []Template{}
	if err := r.db.Select(&ts, sqlListTemplates); err != nil {
		return nil, err
	}

	return ts, nil
}

func (r *pgRepo) DeleteTemplate(ctx context.Context, id string) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	_, err = tx.Exec(sqlDeleteTemplate, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *pgRepo) CreateGroup(ctx context.Context, g Group) (*Group, error) {
	if !g.IsValid() {
		return nil, errors.New(ValidationErr)
	}

	var id string

	tx, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}

	err = sqlx.Get(tx, &id, sqlCreateGroup,
		g.Name, g.BaseName, g.TemplateID, g.MetricType, g.Metric, g.PolicyType, g.Policy)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	g.ID = id

	return &g, nil
}

func (r *pgRepo) SaveGroup(ctx context.Context, g Group) error {
	if !g.IsValid() {
		return errors.New(ValidationErr)
	}

	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	_, err = tx.Exec(sqlUpdateGroup, g.Metric, g.Policy, g.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *pgRepo) GetGroup(ctx context.Context, id string) (*Group, error) {
	row := r.db.QueryRowx(sqlGetGroup, id)

	var name, baseName, templateName, metricType, policyType string
	var metric, policy interface{}

	if err := row.Scan(&name, &baseName, &templateName, &metricType, &metric, &policyType, &policy); err != nil {
		if err == sql.ErrNoRows {
			return nil, ObjectMissingErr
		}

		return nil, err
	}

	g := Group{
		ID:         id,
		Name:       name,
		BaseName:   baseName,
		TemplateID: templateName,
		MetricType: metricType,
		PolicyType: policyType,
	}

	if err := g.LoadPolicy(policy); err != nil {
		return nil, err
	}

	if err := g.LoadMetric(metric); err != nil {
		return nil, err
	}

	return &g, nil
}

func (r *pgRepo) ListGroups(ctx context.Context) ([]Group, error) {
	groups := []Group{}

	rows, err := r.db.Queryx(sqlListGroups)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var id, name, baseName, templateID, metricType, policyType string
		var metric, policy interface{}

		if err := rows.Scan(&id, &name, &baseName, &templateID, &metricType, &metric, &policyType, &policy); err != nil {
			return nil, err
		}

		g := Group{
			ID:         id,
			Name:       name,
			BaseName:   baseName,
			TemplateID: templateID,
			MetricType: metricType,
			PolicyType: policyType,
		}

		if err := g.LoadPolicy(policy); err != nil {
			return nil, err
		}

		if err := g.LoadMetric(metric); err != nil {
			return nil, err
		}

		groups = append(groups, g)
	}

	return groups, nil
}

func (r *pgRepo) DeleteGroup(ctx context.Context, id string) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	_, err = tx.Exec(sqlDeleteGroup, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *pgRepo) Close() error {
	return r.db.Close()
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

	sqlDeleteTemplate = `
  DELETE from templates WHERE id = $1`

	sqlCreateGroup = `
  INSERT into groups
  (name, base_name, template_id, metric_type, metric, policy_type, policy)
  VALUES ($1, $2, $3, $4, $5, $6, $7)
  RETURNING id`

	sqlGetGroup = `
  SELECT name, base_name, template_id, metric_type, metric, policy_type, policy from groups where id=$1`

	sqlListGroups = `
  SELECT id, name, base_name, template_id, metric_type, metric, policy_type, policy from groups`

	sqlDeleteGroup = `
  DELETE from groups WHERE id = $1`

	// TODO figure out what we update
	sqlUpdateGroup = `
  UPDATE groups set metric = $1, policy = $2 WHERE id = $3`
)
