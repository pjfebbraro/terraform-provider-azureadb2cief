package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pjfebbraro/terraform-provider-azureadb2cief/internal/client"
	"github.com/pjfebbraro/terraform-provider-azureadb2cief/internal/datasources"
	"github.com/pjfebbraro/terraform-provider-azureadb2cief/internal/resources"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	// schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
	// 	desc := s.Description
	// 	if s.Default != nil {
	// 		desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
	// 	}
	// 	return strings.TrimSpace(desc)
	// }
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"client_id": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("ARM_CLIENT_ID", ""),
					Description: "The Client ID which should be used for service principal authentication",
				},
				"tenant_id": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("ARM_TENANT_ID", ""),
					Description: "The Tenant ID which should be used. Works with all authentication methods except Managed Identity",
				},
				"client_secret": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("ARM_CLIENT_SECRET", ""),
					Description: "The application password to use when authenticating as a Service Principal using a Client Secret",
				},
				"use_cli": {
					Type:        schema.TypeBool,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("ARM_USE_CLI", true),
					Description: "Allow Azure CLI to be used for Authentication",
				},
			},
			ResourcesMap: map[string]*schema.Resource{
				"azureadb2cief_trust_framework_policy":  resources.TrustFrameworkPolicyResource(),
				"azureadb2cief_trust_framework_key_set": resources.TrustFrameworkKeySetResource(),
			},
			DataSourcesMap: map[string]*schema.Resource{
				"azureadb2cief_client_config": datasources.ClientConfigDataSource(),
			},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {

		authConfig := client.MsGraphClientConfig{
			TenantID:            d.Get("tenant_id").(string),
			ClientID:            d.Get("client_id").(string),
			ClientSecret:        d.Get("client_secret").(string),
			EnableAzureCliToken: d.Get("use_cli").(bool),
		}

		return buildClient(authConfig)
	}
}

func buildClient(config client.MsGraphClientConfig) (*client.Client, diag.Diagnostics) {
	apiClient, err := client.New(config)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return apiClient, nil
}
