package autoscale

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestSaveTemplate(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	defer db.Close()

	mock.ExpectQuery("INSERT into templates (.+) RETURNING id").
		WithArgs("id").
		WithArgs("dev0", "512mb", "ubuntu-14-04-x64", "1,2", "userdata").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("1"))

	repo, err := NewRepository(db)
	assert.NoError(t, err)

	tmpl := &Template{
		Region:     "dev0",
		Size:       "512mb",
		Image:      "ubuntu-14-04-x64",
		RawSSHKeys: "1,2",
		UserData:   "userdata",
	}

	id, err := repo.SaveTemplate(tmpl)
	assert.NoError(t, err)

	assert.Equal(t, 1, id)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTemplate(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	defer db.Close()

	columns := []string{"region", "size", "image", "ssh_keys", "user_data"}

	mock.ExpectQuery("SELECT (.+) from templates (.+)").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows(columns).AddRow("dev0", "512mb", "ubuntu-14-04-x64", "1,2", "userdata"))

	repo, err := NewRepository(db)
	assert.NoError(t, err)

	ogTmpl := &Template{
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

	columns := []string{"id", "region", "size", "image", "ssh_keys", "user_data"}

	mock.ExpectQuery("SELECT (.+) from templates").
		WillReturnRows(sqlmock.NewRows(columns).
			AddRow(1, "dev0", "512mb", "ubuntu-14-04-x64", "1,2", "userdata").
			AddRow(2, "dev0", "512mb", "ubuntu-14-04-x64", "3,4", "userdata"))

	repo, err := NewRepository(db)
	assert.NoError(t, err)

	tmpls, err := repo.ListTemplates()
	assert.NoError(t, err)
	assert.Len(t, tmpls, 2)

	assert.NoError(t, mock.ExpectationsWereMet())
}
