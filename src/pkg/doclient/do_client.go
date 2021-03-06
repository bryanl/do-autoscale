package doclient

import (
	"pkg/do"

	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
)

// Client is our interface to digitalocean.
type Client struct {
	TagsService     do.TagsService
	DropletsService do.DropletsService
	SizesService    do.SizesService
	RegionsService  do.RegionsService
	KeysService     do.KeysService
	AccountsService do.AccountService
}

type tokenSource struct {
	AccessToken string
}

// Token creates a token
func (t *tokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

// New creates an instance of Client.
func New(pat string) *Client {
	ts := &tokenSource{
		AccessToken: pat,
	}

	oc := oauth2.NewClient(oauth2.NoContext, ts)

	godoClient := godo.NewClient(oc)

	dc := &Client{
		DropletsService: do.NewDropletsService(godoClient),
		TagsService:     do.NewTagsService(godoClient),
		SizesService:    do.NewSizesService(godoClient),
		RegionsService:  do.NewRegionsService(godoClient),
		KeysService:     do.NewKeysService(godoClient),
		AccountsService: do.NewAccountService(godoClient),
	}

	return dc
}
