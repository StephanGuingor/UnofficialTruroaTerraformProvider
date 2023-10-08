package truora

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	truora "terraform-provider-truora/truora/client"
)

type Config struct {
	APIKey    string
	APIServer string
}

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("TRUORA_API_KEY", nil),
				Description: "The API key for the Truora API",
			},
			"api_server": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("TRUORA_API_SERVER", "https://api.identity.truora.com"),
				Description: "The API server for the Truora API",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"truora_flow": resourceFlow(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"truora_flow": dataSourceFlow(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	config := Config{
		APIKey:    d.Get("api_key").(string),
		APIServer: d.Get("api_server").(string),
	}

	return config.Client()
}

func (c *Config) Client() (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics
	var opts []truora.TruoraClientOption

	if c.APIKey != "" {
		opts = append(opts, truora.WithAPIKey(c.APIKey))
	}

	if c.APIServer != "" {
		opts = append(opts, truora.WithAPIServer(c.APIServer))
	}

	rc, err := truora.NewClient(opts...)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return rc, diags
}
