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

func TrustFrameworkKeySetResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the key set.  The name must begin with B2C_1A_",
				ValidateFunc: func(i interface{}, s string) (warnings []string, errors []error) {
					name := i.(string)
					if strings.Index(name, "B2C_1A_") == 0 {
						return
					}

					errors = append(errors, fmt.Errorf("trustframeworkkeyset name (%s) must begin with B2C_1A_", name))
					return
				},
			},
			"use": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Specifies if the key use is signature or encryption.  Valid values are 'sig' or 'enc'.",
				ValidateFunc: func(i interface{}, s string) (warnings []string, errors []error) {
					use := i.(string)
					use = strings.ToLower(use)
					if use == "sig" || use == "enc" {
						return
					}

					err := fmt.Errorf("invalid 'use' value: %s, valid values are 'sig' and 'enc'", use)
					errors = append(errors, err)
					return
				},
				DiffSuppressFunc: util.NotCaseSensitive,
			},
			"kty": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Specifies the key encryption type.  Valid values are 'rsa' and 'oct'",
				ValidateFunc: func(i interface{}, s string) (warnings []string, errors []error) {
					kty := i.(string)
					kty = strings.ToLower(kty)
					if kty == "rsa" || kty == "oct" {
						return
					}

					err := fmt.Errorf("invalid 'kty' value: %s, valid values are 'rsa' and 'oct'", kty)
					errors = append(errors, err)
					return
				},
				DiffSuppressFunc: util.NotCaseSensitive,
			},
			"k": {
				Description: "The optional secret value to upload.",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Sensitive:   true,
				Computed:    true,
			},
			"nbf": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "Not valid before date, this value is a NumericDate as defined in RFC 7519.",
			},
			"exp": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "Expiration date, this value is a NumericDate as defined in RFC 7519.",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true, //will require migration
				ForceNew:    true,
				Description: "Specifies the key type.  Valid values are 'secret', 'pfx' and 'pem'",
				ValidateFunc: func(i interface{}, s string) (warnings []string, errors []error) {
					if i == nil {
						return
					}

					typeVal := i.(string)
					typeVal = strings.ToLower(typeVal)
					if typeVal == "pfx" || typeVal == "cer" || typeVal == "secret" {
						return
					}

					err := fmt.Errorf("invalid 'type' value: %s, valid values are 'secret' and 'cer', and 'pfx'", typeVal)
					errors = append(errors, err)
					return
				},
			},
			"value": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Sensitive:   true,
				Description: "When type is pfx or cer, the Base64 encoded pfx or cer file.",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Sensitive:   true,
				Description: "When type is pfx, the password for the pfx file.",
			},
		},
		CreateContext: createKey,
		ReadContext:   readKey,
		DeleteContext: deleteKey,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		SchemaVersion: 1,
	}
}

func deleteKey(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	keySetClient := i.(*client.Client).TrustFrameworkKeySetClient

	_, err := keySetClient.DeleteKey(ctx, data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId("")
	return nil
}

func readKey(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	keySetClient := i.(*client.Client).TrustFrameworkKeySetClient

	id := data.Id()

	keyset, stat, err := keySetClient.GetKeySet(ctx, id)

	if err != nil {
		if stat == http.StatusNotFound {
			log.Printf("[DEBUG] Trust Framework Key Set with Name %q was not found - removing from state", data.Id())
			data.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	data.Set("name", id)

	if len(keyset.Keys) == 0 {
		return nil
	}

	key, stat, err := keySetClient.GetActiveKey(ctx, id)
	if err != nil {
		if stat == http.StatusNotFound {
			log.Printf("[DEBUG] Trust Framework Key Set with Name %q was not found - removing from state", data.Id())
			data.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	data.Set("use", key.Use)
	data.Set("kty", key.Kty)

	if key.Nbf != nil {
		data.Set("nbf", *key.Nbf)
	}
	if key.Exp != nil {
		data.Set("exp", *key.Exp)
	}

	return nil
}

func createKey(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	keySetClient := i.(*client.Client).TrustFrameworkKeySetClient

	name := data.Get("name").(string)
	kty := data.Get("kty").(string)
	use := data.Get("use").(string)
	keyType := data.Get("type").(string)

	if keyType == "pfx" {
		if _, ok := data.GetOk("value"); !ok {
			return diag.Errorf("The 'value' field is required for key type of 'pfx'")
		}
		if _, ok := data.GetOk("password"); !ok {
			return diag.Errorf("The 'password' field is required for key type of 'pfx'")
		}
		if _, ok := data.GetOk("use"); ok {
			return diag.Errorf("The 'use' field should not be specified for key type of 'pfx'")
		}
	} else if keyType == "cer" {
		if _, ok := data.GetOk("value"); !ok {
			return diag.Errorf("The 'value' field is required for key type of 'cer'")
		}
		if _, ok := data.GetOk("use"); ok {
			return diag.Errorf("The 'use' field should not be specified for key type of 'cer'")
		}
	} else {
		if _, ok := data.GetOk("use"); !ok {
			return diag.Errorf("The 'use' field is required for key type of 'secret'")
		}
	}

	key := models.TrustFrameworkKey{
		Kty: &kty,
		Use: &use,
	}

	if expRaw, ok := data.GetOk("exp"); ok {
		exp := expRaw.(int)
		key.Exp = &exp
	}
	if nbfRaw, ok := data.GetOk("nbf"); ok {
		nbf := nbfRaw.(int)
		key.Nbf = &nbf
	}

	keyset, _, err := keySetClient.CreateKey(ctx, name)
	if err != nil {
		return diag.FromErr(err)
	}
	data.SetId(*keyset.Id)

	if keyType == "pfx" {
		password := data.Get("password").(string)
		pfx := data.Get("value").(string)
		pfxKey := models.TrustFrameworkPfxKey{
			Key:      &pfx,
			Password: &password,
		}
		_, err = keySetClient.UploadPFX(ctx, *keyset.Id, pfxKey)
	} else if keyType == "cer" {
		cer := data.Get("value").(string)
		cerKey := models.TrustFrameworkCerKey{
			Key: &cer,
		}
		_, err = keySetClient.UploadCER(ctx, *keyset.Id, cerKey)
	} else {
		specifiedSecret := false
		if rawSecret, ok := data.GetOk("k"); ok {
			secret := rawSecret.(string)
			key.K = &secret
			specifiedSecret = true
		}
		if specifiedSecret {
			_, err = keySetClient.UploadSecret(ctx, *keyset.Id, key)
		} else {
			_, err = keySetClient.GenerateKey(ctx, *keyset.Id, key)
		}
	}

	if err != nil {
		return diag.FromErr(err)
	}

	return readKey(ctx, data, i)
}
