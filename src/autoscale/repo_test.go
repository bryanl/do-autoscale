package autoscale

import (
	"testing"

	"github.com/go-errors/errors"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestCreateTemplate(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	defer db.Close()

	mock.ExpectQuery("INSERT into templates (.+) RETURNING id").
		WithArgs("id").
		WithArgs("a-template", "dev0", "512mb", "ubuntu-14-04-x64", "1,2", "userdata").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("1"))

	repo, err := NewRepository(db)
	assert.NoError(t, err)

	tmpl := &Template{
		Name:       "a-template",
		Region:     "dev0",
		Size:       "512mb",
		Image:      "ubuntu-14-04-x64",
		RawSSHKeys: "1,2",
		UserData:   "userdata",
	}

	id, err := repo.CreateTemplate(tmpl)
	assert.NoError(t, err)

	assert.Equal(t, 1, id)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateTemplate_InvalidName(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	defer db.Close()

	repo, err := NewRepository(db)
	assert.NoError(t, err)

	tmpl := &Template{
		Name:       "!!!",
		Region:     "dev0",
		Size:       "512mb",
		Image:      "ubuntu-14-04-x64",
		RawSSHKeys: "1,2",
		UserData:   "userdata",
	}

	_, err = repo.CreateTemplate(tmpl)

	assert.True(t, errors.Is(err, ValidationErr))

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTemplate(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	defer db.Close()

	columns := []string{"name", "region", "size", "image", "ssh_keys", "user_data"}

	mock.ExpectQuery("SELECT (.+) from templates (.+)").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows(columns).AddRow("a-template", "dev0", "512mb", "ubuntu-14-04-x64", "1,2", "userdata"))

	repo, err := NewRepository(db)
	assert.NoError(t, err)

	ogTmpl := &Template{
		Name:       "a-template",
		Region:     "dev0",
		Size:       "512mb",
		Image:      "ubuntu-14-04-x64",
		RawSSHKeys: "1,2",
		UserData:   "userdata",
	}

	tmpl, err := repo.GetTemplate(1)
	assert.NoError(t, err)
	assert.Equal(t, ogTmpl, tmpl)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestListTemplates(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	defer db.Close()

	columns := []string{"id", "name", "region", "size", "image", "ssh_keys", "user_data"}

	mock.ExpectQuery("SELECT (.+) from templates").
		WillReturnRows(sqlmock.NewRows(columns).
			AddRow(1, "template-1", "dev0", "512mb", "ubuntu-14-04-x64", "1,2", "userdata").
			AddRow(2, "template-2", "dev0", "512mb", "ubuntu-14-04-x64", "3,4", "userdata"))

	repo, err := NewRepository(db)
	assert.NoError(t, err)

	tmpls, err := repo.ListTemplates()
	assert.NoError(t, err)
	assert.Len(t, tmpls, 2)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateGroup(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	defer db.Close()

	mock.ExpectQuery("INSERT into groups (.+) RETURNING id").
		WithArgs("id").
		WithArgs("as", 3, "load", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("abcdefg"))

	repo, err := NewRepository(db)
	assert.NoError(t, err)

	g := &Group{
		BaseName:   "as",
		BaseSize:   3,
		MetricType: "load",
		TemplateID: 1,
	}

	id, err := repo.CreateGroup(g)
	assert.NoError(t, err)

	assert.Equal(t, "abcdefg", id)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetGroup(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	defer db.Close()

	columns := []string{"id", "base_name", "base_size", "metric_type", "template_id"}

	mock.ExpectQuery("SELECT (.+) from groups (.+)").
		WithArgs("abc").
		WillReturnRows(sqlmock.NewRows(columns).AddRow("abc", "as", 3, "load", 1))

	repo, err := NewRepository(db)
	assert.NoError(t, err)

	ogGroup := &Group{
		ID:         "abc",
		BaseName:   "as",
		BaseSize:   3,
		MetricType: "load",
		TemplateID: 1,
	}

	group, err := repo.GetGroup("abc")
	assert.NoError(t, err)
	assert.Equal(t, ogGroup, group)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestListGroups(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	defer db.Close()

	columns := []string{"id", "base_name", "base_size", "metric_type", "template_id"}

	mock.ExpectQuery("SELECT (.+) from groups").
		WillReturnRows(sqlmock.NewRows(columns).
			AddRow("abc", "as", 3, "load", 1).
			AddRow("def", "as2", 3, "load", 2))

	repo, err := NewRepository(db)
	assert.NoError(t, err)

	groups, err := repo.ListGroups()
	assert.NoError(t, err)
	assert.Len(t, groups, 2)

	assert.NoError(t, mock.ExpectationsWereMet())
}
