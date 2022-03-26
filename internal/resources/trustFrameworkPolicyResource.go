package resources

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pjfebbraro/terraform-provider-azureadb2cief/internal/client"
	"github.com/pjfebbraro/terraform-provider-azureadb2cief/internal/models"
	"github.com/pjfebbraro/terraform-provider-azureadb2cief/internal/util"
	"log"
	"net/http"
	"strings"
)

func TrustFrameworkPolicyResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the policy.  The name must begin with B2C_1A_",
				Computed:    false,
				ForceNew:    true,
				ValidateFunc: func(i interface{}, s string) (warnings []string, errors []error) {
					name := i.(string)
					if strings.Index(name, "B2C_1A_") == 0 {
						return
					}

					errors = append(errors, fmt.Errorf("custom_policy_key_set name (%s) must begin with B2C_1A_", name))
					return
				},
			},
			"policy": {
				Description:      "The policy XML",
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: util.XmlDiff,
			},
		},
		SchemaVersion: 1,
		CreateContext: policyResourceCreate,
		ReadContext:   policyResourceRead,
		UpdateContext: policyResourceUpdate,
		DeleteContext: policyResourceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts:    nil,
		Description: "",
	}
}

func policyResourceDelete(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	policyClient := i.(*client.Client).TrustFrameworkPolicyClient

	_, err := policyClient.Delete(ctx, data.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func policyResourceUpdate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	policyClient := i.(*client.Client).TrustFrameworkPolicyClient
	id := data.Id()

	if data.HasChange("policy") {
		xml := data.Get("policy").(string)

		policy := models.Policy{
			Name:   id,
			Policy: xml,
		}
		_, err := policyClient.Update(ctx, &policy)

		if err != nil {
			return diag.FromErr(err)
		}
	}

	return policyResourceRead(ctx, data, i)
}

func policyResourceCreate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	policyClient := i.(*client.Client).TrustFrameworkPolicyClient

	id := data.Get("name").(string)
	xml := data.Get("policy").(string)

	_, err := policyClient.Create(ctx, &xml)

	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(id)
	return policyResourceRead(ctx, data, i)
}

func policyResourceRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	policyClient := i.(*client.Client).TrustFrameworkPolicyClient

	policy, status, err := policyClient.Get(ctx, data.Id())
	if err != nil {
		if status == http.StatusNotFound {
			log.Printf("[DEBUG] Trust Framework Policy with Name %q was not found - removing from state", data.Id())
			data.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	data.Set("policy", policy.Policy)
	data.Set("name", policy.Name)
	return nil
}
