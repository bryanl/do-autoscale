package autoscale

import (
	"pkg/do"
	"pkg/do/mocks"
	"pkg/doclient"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"golang.org/x/net/context"
)

func TestUserConfig(t *testing.T) {
	rs := &mocks.RegionsService{}
	regions := do.Regions{
		{}, {},
	}
	rs.On("List").Return(regions, nil)

	ss := &mocks.SizesService{}
	sizes := do.Sizes{
		{}, {},
	}
	ss.On("List").Return(sizes, nil)

	ks := &mocks.KeysService{}
	keys := do.SSHKeys{
		{}, {},
	}
	ks.On("List").Return(keys, nil)

	as := &mocks.AccountService{}
	a := &do.Account{
		Account: &godo.Account{UUID: "1"},
	}
	as.On("Get").Return(a, nil)

	dc := doclient.Client{
		RegionsService:  rs,
		SizesService:    ss,
		KeysService:     ks,
		AccountsService: as,
	}

	ctx := context.Background()
	uc, err := NewUserConfig(ctx, &dc)
	require.NoError(t, err)

	assert.Equal(t, "1", uc.ID)
	assert.Len(t, uc.Regions, 2)
	assert.Len(t, uc.Sizes, 2)
	assert.Len(t, uc.Keys, 2)

	assert.True(t, rs.AssertExpectations(t))
	assert.True(t, ss.AssertExpectations(t))
	assert.True(t, ks.AssertExpectations(t))
	assert.True(t, as.AssertExpectations(t))
}
