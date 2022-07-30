package resources

import (
	"context"
	"encoding/xml"
	"fmt"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pjfebbraro/terraform-provider-azureadb2cief/internal/client"
	"github.com/pjfebbraro/terraform-provider-azureadb2cief/internal/models"
	"github.com/pjfebbraro/terraform-provider-azureadb2cief/internal/util"
	"io"
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
				DiffSuppressFunc: util.NotCaseSensitive,
			},
			"policy": {
				Description:      "The policy XML",
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: util.XmlDiff,
				ValidateDiagFunc: policyXmlValidate,
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

func policyXmlValidate(val interface{}, p cty.Path) diag.Diagnostics {
	policyXml := val.(string)
	var diags diag.Diagnostics
	reader := strings.NewReader(policyXml)
	decoder := xml.NewDecoder(reader)
	for {
		token, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				diagnostic := diag.Diagnostic{
					Severity:      diag.Error,
					Summary:       "Invalid Policy XML",
					AttributePath: p,
				}
				diags = append(diags, diagnostic)
				break
			}
		}

		if procInstToken, ok := token.(xml.ProcInst); ok {
			if procInstToken.Target != "xml" {
				diagnostic := diag.Diagnostic{
					Severity:      diag.Error,
					Summary:       "Error validating policy xml declaration",
					Detail:        fmt.Sprintf("Only xml is supported.  Xml declaration stated target %s", procInstToken.Target),
					AttributePath: p,
				}
				diags = append(diags, diagnostic)
				break
			}
			var procInst struct {
				Version string `xml:"version,attr"`
			}
			xmlWrappedProcInst := fmt.Sprintf("<procinst %s></procinst>", procInstToken.Inst)
			err = xml.Unmarshal([]byte(xmlWrappedProcInst), &procInst)
			if err != nil {
				diagnostic := diag.Diagnostic{
					Severity:      diag.Error,
					Summary:       "Error validating policy xml declaration",
					Detail:        fmt.Sprintf("Error parsing xml declaration.  Verify that it is valid: %s \n %s", procInstToken.Inst, err),
					AttributePath: p,
				}
				diags = append(diags, diagnostic)
				break
			}

			if procInst.Version != "1.0" {
				diagnostic := diag.Diagnostic{
					Severity:      diag.Error,
					Summary:       "Error validating policy xml declaration",
					Detail:        fmt.Sprintf("Only xml version 1.0 is supported.  Xml declaration stated version=%s", procInst.Version),
					AttributePath: p,
				}
				diags = append(diags, diagnostic)
				break
			}
		}
	}
	return diags
}
