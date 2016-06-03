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

		sshJSON := `[{"id":1,"fingerprint":""},{"id":2,"fingerprint":""}]`
		mock.ExpectQuery("INSERT into templates (.+) RETURNING id").
			WithArgs("id").
			WithArgs("a-template", "dev0", "512mb", "ubuntu-14-04-x64", []uint8(sshJSON), "userdata").
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
			WithArgs("1").
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

		tmpl, err := repo.GetTemplate(ctx, "1")
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
		mock.ExpectExec("DELETE from templates").WithArgs("1").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.DeleteTemplate(ctx, "1")
		require.NoError(t, err)
	})
}

func TestCreateGroup(t *testing.T) {
	withDBMock(t, func(ctx context.Context, repo Repository, mock sqlmock.Sqlmock) {
		m, err := NewFileLoad()
		require.NoError(t, err)

		metricJSON, err := json.Marshal(&m)
		require.NoError(t, err)

		vps := ValuePolicyScale(1, 10, 0.8, 2, 0.2, 1)
		vp, err := NewValuePolicy(vps)
		require.NoError(t, err)

		vpJSON, err := json.Marshal(&vp)
		require.NoError(t, err)

		mock.ExpectBegin()
		mock.ExpectQuery("INSERT into groups (.+) RETURNING id").
			WithArgs("id").
			WithArgs("group", "as", "a-template", "load", []uint8(metricJSON), "value", []uint8(vpJSON)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("abcdefg"))
		mock.ExpectCommit()

		group := Group{
			Name:       "group",
			BaseName:   "as",
			TemplateID: "a-template",
			MetricType: "load",
			Metric:     m,
			PolicyType: "value",
			Policy:     vp,
		}

		g, err := repo.CreateGroup(ctx, group)
		require.NoError(t, err)

		require.Equal(t, "abcdefg", g.ID)
	})
}

func TestCreateGroup_InvalidName(t *testing.T) {
	withDBMock(t, func(ctx context.Context, repo Repository, mock sqlmock.Sqlmock) {
		m, err := NewFileLoad()
		require.NoError(t, err)

		vps := ValuePolicyScale(1, 10, 0.8, 2, 0.2, 1)
		vp, err := NewValuePolicy(vps)

		group := Group{
			Name:       "!!!",
			BaseName:   "as",
			MetricType: "load",
			Metric:     m,
			PolicyType: "value",
			Policy:     vp,
			TemplateID: "a-template",
		}

		_, err = repo.CreateGroup(ctx, group)

		require.True(t, errors.Is(err, ValidationErr), fmt.Sprintf("received %#v", err))

	})
}

func TestGetGroup(t *testing.T) {
	withDBMock(t, func(ctx context.Context, repo Repository, mock sqlmock.Sqlmock) {
		groupColumns := []string{"name", "base_name", "template_id", "metric_type", "metric", "policy_type", "policy"}

		m, err := NewFileLoad()
		require.NoError(t, err)

		mJSON, err := json.Marshal(m)
		require.NoError(t, err)

		pJSON, err := json.Marshal(defaultValuePolicy)
		require.NoError(t, err)

		mock.ExpectQuery("SELECT (.+) from groups (.+)").
			WithArgs("abc").
			WillReturnRows(sqlmock.NewRows(groupColumns).
				AddRow("group-1", "as", "template-1", "load", []uint8(mJSON), "value", []uint8(pJSON)))

		group, err := repo.GetGroup(ctx, "abc")
		require.NoError(t, err)
		require.Equal(t, "group-1", group.Name)

	})
}

func TestListGroups(t *testing.T) {
	withDBMock(t, func(ctx context.Context, repo Repository, mock sqlmock.Sqlmock) {
		groupColumns := []string{"id", "name", "base_name", "template_id", "metric_type", "metric", "policy_type", "policy"}

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
		mock.ExpectExec("DELETE from groups").WithArgs("id").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.DeleteGroup(ctx, "id")
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
		mock.ExpectExec("UPDATE groups").WithArgs(&m, &p, "abc").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		g := Group{
			ID:         "abc",
			Name:       "group",
			BaseName:   "as",
			TemplateID: "a-template",
			MetricType: "load",
			Metric:     m,
			PolicyType: "value",
			Policy:     p,
		}

		err = repo.SaveGroup(ctx, g)
		require.NoError(t, err)
	})
}
