package autoscale

import (
	"encoding/json"
	"fmt"
	"testing"

	"golang.org/x/net/context"

	"github.com/go-errors/errors"
	"github.com/stretchr/testify/require"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

type dbTestFn func(context.Context, Repository, sqlmock.Sqlmock)

func withDBMock(t *testing.T, fn dbTestFn) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer db.Close()

	repo, err := NewRepository(db)
	require.NoError(t, err)

	ctx := context.Background()
	fn(ctx, repo, mock)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateTemplate(t *testing.T) {
	withDBMock(t, func(ctx context.Context, repo Repository, mock sqlmock.Sqlmock) {
		mock.ExpectBegin()
		mock.ExpectQuery("INSERT into templates (.+) RETURNING id").
			WithArgs("id").
			WithArgs("a-template", "dev0", "512mb", "ubuntu-14-04-x64", []uint8(`[{"ID":1},{"ID":2}]`), "userdata").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("id"))
		mock.ExpectCommit()

		in := Template{
			Name:   "a-template",
			Region: "dev0",
			Size:   "512mb",
			Image:  "ubuntu-14-04-x64",
			SSHKeys: []SSHKey{
				{ID: 1},
				{ID: 2},
			},
			UserData: "userdata",
		}

		expected := &Template{
			ID:     "id",
			Name:   "a-template",
			Region: "dev0",
			Size:   "512mb",
			Image:  "ubuntu-14-04-x64",
			SSHKeys: []SSHKey{
				{ID: 1},
				{ID: 2},
			},
			UserData: "userdata",
		}

		tmpl, err := repo.CreateTemplate(ctx, in)
		require.NoError(t, err)

		require.Equal(t, expected, tmpl)

	})
}

func TestCreateTemplate_InvalidName(t *testing.T) {
	withDBMock(t, func(ctx context.Context, repo Repository, mock sqlmock.Sqlmock) {
		in := Template{
			Name:   "!!!",
			Region: "dev0",
			Size:   "512mb",
			Image:  "ubuntu-14-04-x64",
			SSHKeys: []SSHKey{
				{ID: 1},
				{ID: 2},
			},
			UserData: "userdata",
		}

		_, err := repo.CreateTemplate(ctx, in)

		require.True(t, errors.Is(err, ValidationErr))
	})
}

func TestGetTemplate(t *testing.T) {
	withDBMock(t, func(ctx context.Context, repo Repository, mock sqlmock.Sqlmock) {
		columns := []string{"id", "name", "region", "size", "image", "ssh_keys", "user_data"}

		mock.ExpectQuery("SELECT (.+) from templates (.+)").
			WithArgs("a-template").
			WillReturnRows(sqlmock.NewRows(columns).
				AddRow("1", "a-template", "dev0", "512mb", "ubuntu-14-04-x64", []uint8(`[{"ID":1},{"ID":2}]`), "userdata"))

		ogTmpl := &Template{
			ID:     "1",
			Name:   "a-template",
			Region: "dev0",
			Size:   "512mb",
			Image:  "ubuntu-14-04-x64",
			SSHKeys: []SSHKey{
				{ID: 1},
				{ID: 2},
			},
			UserData: "userdata",
		}

		tmpl, err := repo.GetTemplate(ctx, "a-template")
		require.NoError(t, err)
		require.Equal(t, ogTmpl, tmpl)
	})
}

func TestListTemplates(t *testing.T) {
	withDBMock(t, func(ctx context.Context, repo Repository, mock sqlmock.Sqlmock) {
		columns := []string{"id", "name", "region", "size", "image", "ssh_keys", "user_data"}

		mock.ExpectQuery("SELECT (.+) from templates").
			WillReturnRows(sqlmock.NewRows(columns).
				AddRow("1", "template-1", "dev0", "512mb", "ubuntu-14-04-x64", []uint8(`[{"ID":1},{"ID":2}]`), "userdata").
				AddRow("2", "template-2", "dev0", "512mb", "ubuntu-14-04-x64", []uint8(`[{"ID":3},{"ID":4}]`), "userdata"))

		tmpls, err := repo.ListTemplates(ctx)
		require.NoError(t, err)
		require.Len(t, tmpls, 2)

	})
}

func TestDeleteTemplate(t *testing.T) {
	withDBMock(t, func(ctx context.Context, repo Repository, mock sqlmock.Sqlmock) {
		mock.ExpectBegin()
		mock.ExpectExec("DELETE from templates").WithArgs("a-template").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.DeleteTemplate(ctx, "a-template")
		require.NoError(t, err)
	})
}

func TestCreateGroup(t *testing.T) {
	withDBMock(t, func(ctx context.Context, repo Repository, mock sqlmock.Sqlmock) {
		m, err := NewFileLoad()
		require.NoError(t, err)

		metricJSON, err := json.Marshal(&m)
		require.NoError(t, err)

		vp := defaultValuePolicy
		vpJSON, err := json.Marshal(&vp)
		require.NoError(t, err)

		mock.ExpectBegin()
		mock.ExpectQuery("INSERT into groups (.+) RETURNING id").
			WithArgs("id").
			WithArgs("group", "as", "a-template", "load", []uint8(metricJSON), "value", []uint8(vpJSON)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("abcdefg"))
		mock.ExpectCommit()

		cgr := CreateGroupRequest{
			Name:         "group",
			BaseName:     "as",
			TemplateName: "a-template",
			MetricType:   "load",
			Metric:       metricJSON,
			PolicyType:   "value",
			Policy:       vpJSON,
		}

		g, err := repo.CreateGroup(ctx, cgr)
		require.NoError(t, err)

		require.Equal(t, "abcdefg", g.ID)
	})
}

func TestCreateGroup_InvalidName(t *testing.T) {
	withDBMock(t, func(ctx context.Context, repo Repository, mock sqlmock.Sqlmock) {
		policyJSON := []byte(`{
    "scale_up_value": 0.8,
    "scale_up_by": 2,
    "scale_down_value": 0.2,
    "scale_down_by": 1,
    "warm_up_duration": "10s"
  }`)

		metricJSON := []byte(`{
    "stats_dir": "/tmp"
  }`)

		cgr := CreateGroupRequest{
			Name:         "!!!",
			BaseName:     "as",
			MetricType:   "load",
			Metric:       metricJSON,
			PolicyType:   "value",
			Policy:       policyJSON,
			TemplateName: "a-template",
		}

		_, err := repo.CreateGroup(ctx, cgr)

		require.True(t, errors.Is(err, ValidationErr), fmt.Sprintf("received %#v", err))

	})
}

func TestGetGroup(t *testing.T) {
	withDBMock(t, func(ctx context.Context, repo Repository, mock sqlmock.Sqlmock) {
		groupColumns := []string{"id", "base_name", "template_name", "metric_type", "metric", "policy_type", "policy"}

		m, err := NewFileLoad()
		require.NoError(t, err)

		mJSON, err := json.Marshal(m)
		require.NoError(t, err)

		pJSON, err := json.Marshal(defaultValuePolicy)
		require.NoError(t, err)

		mock.ExpectQuery("SELECT (.+) from groups (.+)").
			WithArgs("as").
			WillReturnRows(sqlmock.NewRows(groupColumns).
				AddRow("abc", "as", "template-1", "load", []uint8(mJSON), "value", []uint8(pJSON)))

		group, err := repo.GetGroup(ctx, "as")
		require.NoError(t, err)
		require.Equal(t, "abc", group.ID)

	})
}

func TestListGroups(t *testing.T) {
	withDBMock(t, func(ctx context.Context, repo Repository, mock sqlmock.Sqlmock) {
		groupColumns := []string{"id", "name", "base_name", "template_name", "metric_type", "metric", "policy_type", "policy"}

		m, err := NewFileLoad()
		require.NoError(t, err)

		mJSON, err := json.Marshal(m)
		require.NoError(t, err)

		pJSON, err := json.Marshal(defaultValuePolicy)
		require.NoError(t, err)

		mock.ExpectQuery("SELECT (.+) from groups").
			WillReturnRows(sqlmock.NewRows(groupColumns).
				AddRow("abc", "group1", "as", "template-1", "load", []uint8(mJSON), "value", []uint8(pJSON)).
				AddRow("def", "group2", "as", "template-1", "load", []uint8(mJSON), "value", []uint8(pJSON)))

		groups, err := repo.ListGroups(ctx)
		require.NoError(t, err)
		require.Len(t, groups, 2)
	})
}

func TestDeleteGroup(t *testing.T) {
	withDBMock(t, func(ctx context.Context, repo Repository, mock sqlmock.Sqlmock) {
		mock.ExpectBegin()
		mock.ExpectExec("DELETE from groups").WithArgs("a-group").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.DeleteGroup(ctx, "a-group")
		require.NoError(t, err)
	})
}

func TestUpdateGroup(t *testing.T) {
	withDBMock(t, func(ctx context.Context, repo Repository, mock sqlmock.Sqlmock) {
		m, err := NewFileLoad()
		require.NoError(t, err)

		p, err := NewValuePolicy()
		require.NoError(t, err)
		p.vpd = defaultValuePolicy

		mock.ExpectBegin()
		mock.ExpectExec("UPDATE groups").WithArgs(&m, &p, "group").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		g := Group{
			ID:           "abc",
			Name:         "group",
			BaseName:     "as",
			TemplateName: "a-template",
			MetricType:   "load",
			Metric:       m,
			PolicyType:   "value",
			Policy:       p,
		}

		err = repo.SaveGroup(ctx, g)
		require.NoError(t, err)
	})
}
