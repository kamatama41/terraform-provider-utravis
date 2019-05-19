package utravis

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/shuheiktgw/go-travis"
)

func Provider() terraform.ResourceProvider {
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("TRAVIS_API_TOKEN", nil),
				Description: "The API token used to connect to Travis CI.",
			},
			"base_url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("TRAVIS_BASE_URL", nil),
				Description: fmt.Sprintf("The Travis CI Base API URL (must be either %s or %s)", travis.ApiOrgUrl, travis.ApiComUrl),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"utravis_env_var":  resourceTravisEnvVar(),
			"utravis_key_pair": resourceTravisKeyPair(),
		},
	}
	p.ConfigureFunc = providerConfigure(p)

	return p
}

func providerConfigure(p *schema.Provider) schema.ConfigureFunc {
	return func(d *schema.ResourceData) (interface{}, error) {
		baseUrl := d.Get("base_url").(string)
		token := d.Get("token").(string)
		cfg, err := NewConfig(baseUrl, token)
		if err != nil {
			return nil, err
		}
		return cfg, nil
	}
}
