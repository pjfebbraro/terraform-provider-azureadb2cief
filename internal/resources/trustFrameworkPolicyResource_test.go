package resources_test

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/pjfebbraro/terraform-provider-azureadb2cief/internal/acceptance"
	"github.com/pjfebbraro/terraform-provider-azureadb2cief/internal/client"
	"net/http"
	"os"
	"regexp"
	"testing"
)

func TestAccCustomPolicy(t *testing.T) {
	resourceName := "azureadb2cief_trust_framework_policy.TrustFrameworkBase"

	resource.Test(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"azureadb2cief": func() (*schema.Provider, error) {
				return acceptance.AzureADB2CProvider, nil
			},
		},
		PreCheck: func() {
			preCheckEnv(t)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccResourceScaffolding,
				Check: resource.ComposeTestCheckFunc(
					testPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "B2C_1A_TrustFrameworkBase"),
				),
			},
		},
	})
}

func TestAccCustomPolicyE2E(t *testing.T) {
	baseResourceName := "azureadb2cief_trust_framework_policy.TrustFrameworkBase"
	extResourceName := "azureadb2cief_trust_framework_policy.TrustFrameworkExtensions"

	resource.Test(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"azureadb2cief": func() (*schema.Provider, error) {
				return acceptance.AzureADB2CProvider, nil
			},
		},
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: "0.7.2",
			},
			"azuread": {
				Source: "hashicorp/azuread",
			},
			"random": {
				Source: "hashicorp/random",
			},
		},
		PreCheck: func() {
			preCheckEnv(t)
		},
		Steps: []resource.TestStep{
			{
				Config: testEnd2End,
				Check: resource.ComposeTestCheckFunc(
					testPolicyExists(baseResourceName),
					testPolicyExists(extResourceName),
					resource.TestCheckResourceAttr(baseResourceName, "name", "B2C_1A_TrustFrameworkBase"),
					resource.TestCheckResourceAttr(extResourceName, "name", "B2C_1A_TrustFrameworkExtensions"),
				),
			},
		},
	})
}

func TestAccCustomPolicyInvalidXml(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"azureadb2cief": func() (*schema.Provider, error) {
				return acceptance.AzureADB2CProvider, nil
			},
		},
		PreCheck: func() {
			preCheckEnv(t)
		},
		Steps: []resource.TestStep{
			{
				Config:      testInvalidPolicyXml,
				ExpectError: regexp.MustCompile("Invalid Policy XML"),
			},
		},
	})
}
func TestAccCustomPolicyInvalidXml2(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"azureadb2cief": func() (*schema.Provider, error) {
				return acceptance.AzureADB2CProvider, nil
			},
		},
		PreCheck: func() {
			preCheckEnv(t)
		},
		Steps: []resource.TestStep{
			{
				Config:      testInvalidPolicyXml2,
				ExpectError: regexp.MustCompile("Invalid Policy XML"),
			},
		},
	})
}
func preCheckEnv(t *testing.T) {
	variables := []string{
		"TF_VAR_tenant_name",
		"TF_VAR_tenant_object_id",
	}

	for _, variable := range variables {
		value := os.Getenv(variable)
		if value == "" {
			t.Fatalf("`%s` must be set for acceptance tests %s", variable, t.Name())
		}
	}
}

func testPolicyExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		policyClient := acceptance.AzureADB2CProvider.Meta().(*client.Client).TrustFrameworkPolicyClient
		id := state.RootModule().Resources[resourceName].Primary.ID
		r, status, err := policyClient.Get(context.Background(), id)
		if err != nil {
			return err
		}

		if status == http.StatusOK && r != nil && r.Name == id {
			return nil
		}

		return fmt.Errorf("trustframeworkpolicy policyClient returned invalid status (%v) when checking azure", status)
	}
}

const testAccResourceScaffolding = `
provider "azureadb2cief" {}
resource "azureadb2cief_trust_framework_key_set" "TokenSigningKeyContainer" {
  name = "B2C_1A_TokenSigningKeyContainer"
  use = "sig"
  kty = "RSA"
}
resource "azureadb2cief_trust_framework_key_set" "TokenEncryptionKeyContainer" {
  name = "B2C_1A_TokenEncryptionKeyContainer"
  use = "enc"
  kty = "RSA"
}
resource "azureadb2cief_trust_framework_policy" "TrustFrameworkBase" {
  name = "B2C_1A_TrustFrameworkBase"
  policy = templatefile("${path.module}/testdata/B2C_1A_TrustFrameworkBase.xml", local.template_vars)
}

locals {
  template_vars = {
    tenant_name = var.tenant_name
    tenant_object_id = var.tenant_object_id
    is_development = var.is_development
	token_signing_key_container = azureadb2cief_trust_framework_key_set.TokenSigningKeyContainer.name
	token_encryption_key_container = azureadb2cief_trust_framework_key_set.TokenEncryptionKeyContainer.name
  }
}

variable "tenant_name" {
  type = string
}

variable "tenant_object_id" {
  type = string
}

variable "is_development" {
  type = bool
  default = true
}
`

const testEnd2End = `
provider "azureadb2cief" {}
provider "azuread" {}
provider "random" {}

variable "tenant_name" {
  type = string
}

variable "tenant_object_id" {
  type = string
}

variable "is_development" {
  type = bool
  default = true
}

data "azuread_client_config" "current" {}
data "azuread_service_principal" "MicrosoftGraph" {
  display_name = "Microsoft Graph"
}
resource "random_uuid" "AccessIEFScopeId" {}

resource "azuread_application" "IdentityExperienceFramework" {
  display_name                   = "TestIdentityExperienceFramework"
  oauth2_post_response_required  = false
  owners                         = [
    data.azuread_client_config.current.object_id,
  ]
  sign_in_audience               = "AzureADMyOrg"
  api {
    known_client_applications      = []
    mapped_claims_enabled          = false
    requested_access_token_version = 1

    oauth2_permission_scope {
      admin_consent_description  = "Allow the application to access the IdentityExperienceFramework on behalf of the signed-in user."
      admin_consent_display_name = "Access IdentityExperienceFramework"
      enabled                    = true
      id                         = random_uuid.AccessIEFScopeId.result
      type                       = "Admin"
      value                      = "user_impersonation"
    }
  }

  required_resource_access {
    resource_app_id = data.azuread_service_principal.MicrosoftGraph.application_id

    resource_access {
      id   = data.azuread_service_principal.MicrosoftGraph.oauth2_permission_scope_ids["openid"]
      type = "Scope"
    }
    resource_access {
      id   = data.azuread_service_principal.MicrosoftGraph.oauth2_permission_scope_ids["offline_access"]
      type = "Scope"
    }
  }

  web {
    redirect_uris = [
      "https://${var.tenant_name}.b2clogin.com/${var.tenant_name}.onmicrosoft.com",
    ]

    implicit_grant {
      access_token_issuance_enabled = false
      id_token_issuance_enabled     = false
    }
  }

}

resource "azuread_application" "ProxyIdentityExperienceFramework" {
  display_name                   = "TestProxyIdentityExperienceFramework"
  fallback_public_client_enabled = true
  owners                         = [
    data.azuread_client_config.current.object_id,
  ]
  sign_in_audience               = "AzureADMyOrg"
  tags                           = []

  api {
    known_client_applications      = []
    mapped_claims_enabled          = false
    requested_access_token_version = 1
  }
  public_client {
    redirect_uris = [
      "myapp://auth",
    ]
  }

  required_resource_access {
    resource_app_id = data.azuread_service_principal.MicrosoftGraph.application_id

    resource_access {
      id   = data.azuread_service_principal.MicrosoftGraph.oauth2_permission_scope_ids["openid"]
      type = "Scope"
    }
    resource_access {
      id   = data.azuread_service_principal.MicrosoftGraph.oauth2_permission_scope_ids["offline_access"]
      type = "Scope"
    }
  }
  required_resource_access {
    resource_app_id = azuread_application.IdentityExperienceFramework.application_id

    resource_access {
      id = azuread_application.IdentityExperienceFramework.oauth2_permission_scope_ids["user_impersonation"]
      type = "Scope"
    }
  }

  single_page_application {
    redirect_uris = []
  }

  timeouts {}

  web {
    redirect_uris = []

    implicit_grant {
      access_token_issuance_enabled = false
      id_token_issuance_enabled     = false
    }
  }
}

resource "azureadb2cief_trust_framework_key_set" "TokenSigningKeyContainer" {
  name = "B2C_1A_TokenSigningKeyContainer"
  use = "sig"
  kty = "RSA"
}
resource "azureadb2cief_trust_framework_key_set" "TokenEncryptionKeyContainer" {
  name = "B2C_1A_TokenEncryptionKeyContainer"
  use = "enc"
  kty = "RSA"
}

resource "azureadb2cief_trust_framework_policy" "TrustFrameworkBase" {
  name = "B2C_1A_TrustFrameworkBase"
  policy = templatefile("${path.module}/testdata/B2C_1A_TrustFrameworkBase.xml", local.trust_framework_base_template_vars)
}

resource "azureadb2cief_trust_framework_policy" "TrustFrameworkExtensions" {
  name = "B2C_1A_TrustFrameworkExtensions"
  policy = templatefile("${path.module}/testdata/TrustFrameworkExtensions.xml", local.trust_framework_ext_template_vars)
}

locals {
  trust_framework_base_template_vars = {
    tenant_name = var.tenant_name
    tenant_object_id = var.tenant_object_id
    is_development = var.is_development
	token_signing_key_container = azureadb2cief_trust_framework_key_set.TokenSigningKeyContainer.name
	token_encryption_key_container = azureadb2cief_trust_framework_key_set.TokenEncryptionKeyContainer.name
  }
  trust_framework_ext_template_vars = {
    tenant_name = var.tenant_name
    tenant_object_id = var.tenant_object_id
    is_development = var.is_development
	token_signing_key_container = azureadb2cief_trust_framework_key_set.TokenSigningKeyContainer.name
	token_encryption_key_container = azureadb2cief_trust_framework_key_set.TokenEncryptionKeyContainer.name
	base_policy = azureadb2cief_trust_framework_policy.TrustFrameworkBase.name
	ProxyIdentityExperienceFrameworkAppId = azuread_application.ProxyIdentityExperienceFramework.application_id
	IdentityExperienceFrameworkAppId = azuread_application.IdentityExperienceFramework.application_id
  }
}
`

const testInvalidPolicyXml = `
provider "azureadb2cief" {}
resource "azureadb2cief_trust_framework_key_set" "TokenSigningKeyContainer" {
  name = "B2C_1A_TokenSigningKeyContainer"
  use = "sig"
  kty = "RSA"
}
resource "azureadb2cief_trust_framework_key_set" "TokenEncryptionKeyContainer" {
  name = "B2C_1A_TokenEncryptionKeyContainer"
  use = "enc"
  kty = "RSA"
}
resource "azureadb2cief_trust_framework_policy" "TrustFrameworkBase" {
  name = "B2C_1A_TrustFrameworkBase"
  policy = templatefile("${path.module}/testdata/B2C_1A_TrustFrameworkBase_invalid_xml.xml", local.template_vars)
}

locals {
  template_vars = {
    tenant_name = var.tenant_name
    tenant_object_id = var.tenant_object_id
    is_development = var.is_development
	token_signing_key_container = azureadb2cief_trust_framework_key_set.TokenSigningKeyContainer.name
	token_encryption_key_container = azureadb2cief_trust_framework_key_set.TokenEncryptionKeyContainer.name
  }
}

variable "tenant_name" {
  type = string
}

variable "tenant_object_id" {
  type = string
}

variable "is_development" {
  type = bool
  default = true
}
`

const testInvalidPolicyXml2 = `
provider "azureadb2cief" {}
resource "azureadb2cief_trust_framework_key_set" "TokenSigningKeyContainer" {
  name = "B2C_1A_TokenSigningKeyContainer"
  use = "sig"
  kty = "RSA"
}
resource "azureadb2cief_trust_framework_key_set" "TokenEncryptionKeyContainer" {
  name = "B2C_1A_TokenEncryptionKeyContainer"
  use = "enc"
  kty = "RSA"
}
resource "azureadb2cief_trust_framework_policy" "TrustFrameworkBase" {
  name = "B2C_1A_TrustFrameworkBase"
  policy = templatefile("${path.module}/testdata/B2C_1A_TrustFrameworkBase_invalid_xml.xml", local.template_vars)
}

locals {
  template_vars = {
    tenant_name = var.tenant_name
    tenant_object_id = var.tenant_object_id
    is_development = var.is_development
	token_signing_key_container = azureadb2cief_trust_framework_key_set.TokenSigningKeyContainer.name
	token_encryption_key_container = azureadb2cief_trust_framework_key_set.TokenEncryptionKeyContainer.name
  }
}

variable "tenant_name" {
  type = string
}

variable "tenant_object_id" {
  type = string
}

variable "is_development" {
  type = bool
  default = true
}
`
