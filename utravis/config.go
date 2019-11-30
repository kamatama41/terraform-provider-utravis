package utravis

import (
	"fmt"
	"sync"

	"github.com/shuheiktgw/go-travis"
)

type config struct {
	client *travis.Client
	mutex  sync.RWMutex
}

func NewConfig(baseUrl, token string) (*config, error) {
	if baseUrl != travis.ApiOrgUrl && baseUrl != travis.ApiComUrl {
		return nil, fmt.Errorf("base_url must be either %s or %s", travis.ApiOrgUrl, travis.ApiComUrl)
	}
	return &config{client: travis.NewClient(baseUrl, token)}, nil
}
