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
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("id"))

	repo, err := NewRepository(db)
	assert.NoError(t, err)

	ctr := CreateTemplateRequest{
		Name:     "a-template",
		Region:   "dev0",
		Size:     "512mb",
		Image:    "ubuntu-14-04-x64",
		SSHKeys:  []string{"1", "2"},
		UserData: "userdata",
	}

	expected := Template{
		ID:       "id",
		Name:     "a-template",
		Region:   "dev0",
		Size:     "512mb",
		Image:    "ubuntu-14-04-x64",
		SSHKeys:  []string{"1", "2"},
		UserData: "userdata",
	}

	tmpl, err := repo.CreateTemplate(ctr)
	assert.NoError(t, err)

	assert.Equal(t, expected, tmpl)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateTemplate_InvalidName(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	defer db.Close()

	repo, err := NewRepository(db)
	assert.NoError(t, err)

	ctr := CreateTemplateRequest{
		Name:     "!!!",
		Region:   "dev0",
		Size:     "512mb",
		Image:    "ubuntu-14-04-x64",
		SSHKeys:  []string{"1", "2"},
		UserData: "userdata",
	}

	_, err = repo.CreateTemplate(ctr)

	assert.True(t, errors.Is(err, ValidationErr))

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTemplate(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	defer db.Close()

	columns := []string{"id", "name", "region", "size", "image", "ssh_keys", "user_data"}

	mock.ExpectQuery("SELECT (.+) from templates (.+)").
		WithArgs("a-template").
		WillReturnRows(sqlmock.NewRows(columns).
			AddRow("1", "a-template", "dev0", "512mb", "ubuntu-14-04-x64", "1,2", "userdata"))

	repo, err := NewRepository(db)
	assert.NoError(t, err)

	ogTmpl := Template{
		ID:       "1",
		Name:     "a-template",
		Region:   "dev0",
		Size:     "512mb",
		Image:    "ubuntu-14-04-x64",
		SSHKeys:  []string{"1", "2"},
		UserData: "userdata",
	}

	tmpl, err := repo.GetTemplate("a-template")
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
			AddRow("1", "template-1", "dev0", "512mb", "ubuntu-14-04-x64", "1,2", "userdata").
			AddRow("2", "template-2", "dev0", "512mb", "ubuntu-14-04-x64", "3,4", "userdata"))

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
		WithArgs("group", "as", 3, "load", "a-template").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("abcdefg"))

	repo, err := NewRepository(db)
	assert.NoError(t, err)

	cgr := CreateGroupRequest{
		Name:         "group",
		BaseName:     "as",
		BaseSize:     3,
		MetricType:   "load",
		TemplateName: "a-template",
	}

	id, err := repo.CreateGroup(cgr)
	assert.NoError(t, err)

	expected := Group{
		ID:           "abcdefg",
		Name:         "group",
		BaseName:     "as",
		BaseSize:     3,
		MetricType:   "load",
		TemplateName: "a-template",
	}

	assert.Equal(t, expected, id)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateGroup_InvalidName(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	defer db.Close()

	repo, err := NewRepository(db)
	assert.NoError(t, err)

	cgr := CreateGroupRequest{
		Name:         "!!!",
		BaseName:     "as",
		BaseSize:     3,
		MetricType:   "load",
		TemplateName: "a-template",
	}

	_, err = repo.CreateGroup(cgr)

	assert.True(t, errors.Is(err, ValidationErr))

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetGroup(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	defer db.Close()

	columns := []string{"id", "name", "base_name", "base_size", "metric_type", "template_name"}

	mock.ExpectQuery("SELECT (.+) from groups (.+)").
		WithArgs("as").
		WillReturnRows(sqlmock.NewRows(columns).
			AddRow("abc", "group", "as", 3, "load", "a-template"))

	repo, err := NewRepository(db)
	assert.NoError(t, err)

	ogGroup := Group{
		ID:           "abc",
		Name:         "group",
		BaseName:     "as",
		BaseSize:     3,
		MetricType:   "load",
		TemplateName: "a-template",
	}

	group, err := repo.GetGroup("as")
	assert.NoError(t, err)
	assert.Equal(t, ogGroup, group)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestListGroups(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	defer db.Close()

	columns := []string{"id", "name", "base_name", "base_size", "metric_type", "template_name"}

	mock.ExpectQuery("SELECT (.+) from groups").
		WillReturnRows(sqlmock.NewRows(columns).
			AddRow("abc", "group1", "as", 3, "load", "template-1").
			AddRow("def", "group2", "as2", 3, "load", "template-1"))

	repo, err := NewRepository(db)
	assert.NoError(t, err)

	groups, err := repo.ListGroups()
	assert.NoError(t, err)
	assert.Len(t, groups, 2)

	assert.NoError(t, mock.ExpectationsWereMet())
}
