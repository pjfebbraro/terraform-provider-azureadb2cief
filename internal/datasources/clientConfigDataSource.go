package datasources

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pjfebbraro/terraform-provider-azureadb2cief/internal/client"
)

func ClientConfigDataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: clientConfigDataSourceRead,
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Description: "The tenant ID of the authenticated principal",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func clientConfigDataSourceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	graphClient := meta.(*client.Client)
	d.SetId(fmt.Sprintf("%s", graphClient.Config.TenantID))
	d.Set("tenant_id", graphClient.Config.TenantID)

	return nil
}
