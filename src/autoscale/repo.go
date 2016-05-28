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
	CreateTemplate(ctx context.Context, tcr CreateTemplateRequest) (*Template, error)
	GetTemplate(ctx context.Context, name string) (*Template, error)
	ListTemplates(ctx context.Context) ([]*Template, error)
	DeleteTemplate(ctx context.Context, name string) error

	CreateGroup(ctx context.Context, gcr CreateGroupRequest) (Group, error)
	GetGroup(ctx context.Context, name string) (Group, error)
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

func (r *pgRepo) CreateTemplate(ctx context.Context, tcr CreateTemplateRequest) (*Template, error) {

	options := tcr.Options

	t := Template{
		Name:     options.Name,
		Region:   options.Region,
		Size:     options.Size,
		Image:    options.Image,
		SSHKeys:  options.SSHKeys,
		UserData: options.UserData,
	}

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

func (r *pgRepo) GetTemplate(ctx context.Context, name string) (*Template, error) {
	var t Template
	if err := r.db.Get(&t, sqlGetTemplate, name); err != nil {
		if err == sql.ErrNoRows {
			return nil, ObjectMissingErr
		}

		return nil, err
	}

	return &t, nil
}

func (r *pgRepo) ListTemplates(ctx context.Context) ([]*Template, error) {
	ts := []*Template{}
	if err := r.db.Select(&ts, sqlListTemplates); err != nil {
		return nil, err
	}

	return ts, nil
}

func (r *pgRepo) DeleteTemplate(ctx context.Context, name string) error {
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

func (r *pgRepo) CreateGroup(ctx context.Context, gcr CreateGroupRequest) (Group, error) {
	g, err := gcr.ConvertToGroup(ctx)
	if err != nil {
		return Group{}, err
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
		g.Name, g.BaseName, g.TemplateName, g.MetricType, g.Metric, g.PolicyType, g.Policy)
	if err != nil {
		tx.Rollback()
		return Group{}, err
	}

	err = tx.Commit()
	if err != nil {
		return Group{}, err
	}

	g.ID = id

	return *g, nil
}

func (r *pgRepo) SaveGroup(ctx context.Context, g Group) error {
	if !g.IsValid() {
		return errors.New(ValidationErr)
	}

	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	_, err = tx.Exec(sqlUpdateGroup, g.Metric, g.Policy, g.Name)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *pgRepo) GetGroup(ctx context.Context, name string) (Group, error) {
	row := r.db.QueryRowx(sqlGetGroup, name)

	var id, baseName, templateName, metricType, policyType string
	var metric, policy interface{}

	if err := row.Scan(&id, &baseName, &templateName, &metricType, &metric, &policyType, &policy); err != nil {
		if err == sql.ErrNoRows {
			return Group{}, ObjectMissingErr
		}

		return Group{}, err
	}

	g := Group{
		ID:           id,
		Name:         name,
		BaseName:     baseName,
		TemplateName: templateName,
		MetricType:   metricType,
		PolicyType:   policyType,
	}

	if err := g.LoadPolicy(policy); err != nil {
		return Group{}, err
	}

	if err := g.LoadMetric(metric); err != nil {
		return Group{}, err
	}

	return g, nil
}

func (r *pgRepo) ListGroups(ctx context.Context) ([]Group, error) {
	groups := []Group{}

	rows, err := r.db.Queryx(sqlListGroups)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var id, name, baseName, templateName, metricType, policyType string
		var metric, policy interface{}

		if err := rows.Scan(&id, &name, &baseName, &templateName, &metricType, &metric, &policyType, &policy); err != nil {
			return nil, err
		}

		g := Group{
			ID:           id,
			Name:         name,
			BaseName:     baseName,
			TemplateName: templateName,
			MetricType:   metricType,
			PolicyType:   policyType,
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

func (r *pgRepo) DeleteGroup(ctx context.Context, name string) error {
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
  SELECT * from templates where name=$1`

	sqlListTemplates = `
  SELECT * from templates`

	sqlDeleteTemplate = `
  DELETE from templates WHERE id = $1`

	sqlCreateGroup = `
  INSERT into groups
  (name, base_name, template_name, metric_type, metric, policy_type, policy)
  VALUES ($1, $2, $3, $4, $5, $6, $7)
  RETURNING id`

	sqlGetGroup = `
  SELECT id, base_name, template_name, metric_type, metric, policy_type, policy from groups where name=$1`

	sqlListGroups = `
  SELECT id, name, base_name, template_name, metric_type, metric, policy_type, policy from groups`

	sqlDeleteGroup = `
  DELETE from groups WHERE name = $1`

	// TODO figure out what we update
	sqlUpdateGroup = `
  UPDATE groups set metrics = $1, policy = $2 WHERE name = $3`
)
