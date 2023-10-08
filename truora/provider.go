package truora

import (
	"context"
	"time"

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
			"truora_flow":          dataSourceFlow(),
			"truora_flow_document": dataSourceFlowDocument(),
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

func mergeSchemaMaps(maps ...map[string]*schema.Schema) map[string]*schema.Schema {
	result := map[string]*schema.Schema{}
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

func timeFromRFC3339(s string) *time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return nil
	}
	return &t
}

func getBoolFromMap(data map[string]interface{}, key string) bool {
	if v, ok := data[key]; ok {
		return v.(bool)
	}
	return false
}

func getStringFromMap(data map[string]interface{}, key string) string {
	if v, ok := data[key]; ok {
		return v.(string)
	}
	return ""
}

func getTimeFromMap(data map[string]interface{}, key string) *time.Time {
	if v, ok := data[key]; ok {
		return timeFromRFC3339(v.(string))
	}
	return nil
}

func getInt64FromMap(data map[string]interface{}, key string) int64 {
	if v, ok := data[key]; ok {
		return int64(v.(int))
	}
	return 0
}

func parseStringArray(data map[string]interface{}, key string) []string {
	if v, ok := data[key]; ok {
		stringArray := make([]string, 0)

		for _, item := range v.([]interface{}) {
			stringArray = append(stringArray, item.(string))
		}

		return stringArray
	}
	return nil
}

func getMapFromSingleElementList(data map[string]interface{}, key string) map[string]interface{} {
	if v, ok := data[key]; ok {
		list := v.([]interface{})
		if len(list) > 0 {
			return list[0].(map[string]interface{})
		}
	}
	return nil
}
