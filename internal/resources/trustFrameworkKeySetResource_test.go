package resources_test

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/pjfebbraro/terraform-provider-azureadb2cief/internal/acceptance"
	"github.com/pjfebbraro/terraform-provider-azureadb2cief/internal/client"
	"testing"
)

func TestAccTrustFrameworkKeyset_basic(t *testing.T) {
	resourceName := "azureadb2cief_trust_framework_key_set.TestKeySet"

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
		},
		Steps: []resource.TestStep{
			{
				Config: testAccTrustFrameworkKeyset1,
				Check: resource.ComposeTestCheckFunc(
					checkTrustFrameworkKeySetExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "nbf"),
					resource.TestCheckResourceAttrSet(resourceName, "exp"),
				),
			},
			{
				ImportStateVerify: true,
				ImportState:       true,
				ResourceName:      resourceName,
			},
		},
	})
}
func TestAccTrustFrameworkKeyset_nonbfexp(t *testing.T) {
	resourceName := "azureadb2cief_trust_framework_key_set.TestKeySet"

	resource.Test(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"azureadb2cief": func() (*schema.Provider, error) {
				return acceptance.AzureADB2CProvider, nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccTrustFrameworkKeyset2,
				Check: resource.ComposeTestCheckFunc(
					checkTrustFrameworkKeySetExists(resourceName),
					resource.TestCheckNoResourceAttr(resourceName, "nbf"),
					resource.TestCheckNoResourceAttr(resourceName, "exp"),
				),
			},
			{
				ImportStateVerify: true,
				ImportState:       true,
				ResourceName:      resourceName,
			},
		},
	})
}
func TestAccTrustFrameworkKeyset_basic2(t *testing.T) {
	resourceName := "azureadb2cief_trust_framework_key_set.TestKeySet"

	resource.Test(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"azureadb2cief": func() (*schema.Provider, error) {
				return acceptance.AzureADB2CProvider, nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccTrustFrameworkKeyset3,
				Check: resource.ComposeTestCheckFunc(
					checkTrustFrameworkKeySetExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				ImportStateVerify: true,
				ImportState:       true,
				ResourceName:      resourceName,
			},
		},
	})
}
func TestAccTrustFrameworkKeyset_UploadSecret(t *testing.T) {
	resourceName := "azureadb2cief_trust_framework_key_set.TestKeySetWithSecret"

	resource.Test(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"azureadb2cief": func() (*schema.Provider, error) {
				return acceptance.AzureADB2CProvider, nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccTrustFrameworkKeysetWithSecret,
				Check: resource.ComposeTestCheckFunc(
					checkTrustFrameworkKeySetExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"k"},
				ImportState:             true,
				ResourceName:            resourceName,
			},
		},
	})
}
func checkTrustFrameworkKeySetExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		keySetClient := acceptance.AzureADB2CProvider.Meta().(*client.Client).TrustFrameworkKeySetClient
		id := state.RootModule().Resources[resourceName].Primary.ID
		_, _, err := keySetClient.GetActiveKey(context.TODO(), id)
		if err != nil {
			return err
		}

		return nil
	}
}

const testAccTrustFrameworkKeyset1 = `
provider "azureadb2cief" {}
resource "time_static" "testkeyset_nbf" {}
resource "time_offset" "testkeyset_exp" {
  offset_years = 2
}
resource "azureadb2cief_trust_framework_key_set" "TestKeySet" {
  name = "B2C_1A_TestKeySet"
  use = "enc"
  kty = "RSA"
  nbf = time_static.testkeyset_nbf.unix
  exp = time_offset.testkeyset_exp.unix
}
`
const testAccTrustFrameworkKeyset2 = `
provider "azureadb2cief" {}
resource "azureadb2cief_trust_framework_key_set" "TestKeySet" {
  name = "B2C_1A_TestKeySet"
  use = "enc"
  kty = "RSA"
}
`
const testAccTrustFrameworkKeyset3 = `
provider "azureadb2cief" {}
resource "azureadb2cief_trust_framework_key_set" "TestKeySet" {
  name = "B2C_1A_TestKeySet"
  use = "enc"
  kty = "RSA"
}
`
const testAccTrustFrameworkKeysetWithSecret = `
provider "azureadb2cief" {}
resource "azureadb2cief_trust_framework_key_set" "TestKeySetWithSecret" {
  name = "B2C_1A_TestKeySetWithSecret"
  use = "enc"
  kty = "oct"
  k = "mypassword"
}
`
