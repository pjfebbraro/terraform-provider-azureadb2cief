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
func TestAccTrustFrameworkKeyset_UploadCer(t *testing.T) {
	resourceName := "azureadb2cief_trust_framework_key_set.TestKeySetWithCer"

	resource.Test(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"azureadb2cief": func() (*schema.Provider, error) {
				return acceptance.AzureADB2CProvider, nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccTrustFrameworkKeysetWithCer,
				Check: resource.ComposeTestCheckFunc(
					checkTrustFrameworkKeySetExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"k", "value", "type"},
				ImportState:             true,
				ResourceName:            resourceName,
			},
		},
	})
}
func TestAccTrustFrameworkKeyset_UploadPfx(t *testing.T) {
	resourceName := "azureadb2cief_trust_framework_key_set.TestKeySetWithPfx"

	resource.Test(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"azureadb2cief": func() (*schema.Provider, error) {
				return acceptance.AzureADB2CProvider, nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccTrustFrameworkKeysetWithPfx,
				Check: resource.ComposeTestCheckFunc(
					checkTrustFrameworkKeySetExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"k", "value", "type", "password"},
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
func TestAccTrustFrameworkKeyset_mixedCasesDiffIgnore(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"azureadb2cief": func() (*schema.Provider, error) {
				return acceptance.AzureADB2CProvider, nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccTrustFrameworkKeyset3,
			},
			{
				Config:   testAccTrustFrameworkKeyset3mixcase,
				PlanOnly: true,
			},
		},
	})
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
const testAccTrustFrameworkKeyset3mixcase = `
provider "azureadb2cief" {}
resource "azureadb2cief_trust_framework_key_set" "TestKeySet" {
  name = "B2C_1A_TestKeySet"
  use = "ENC"
  kty = "rsa"
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
const testAccTrustFrameworkKeysetWithCer = `
provider "azureadb2cief" {}
resource "azureadb2cief_trust_framework_key_set" "TestKeySetWithCer" {
  name = "B2C_1A_TestKeySetWithCer"
  kty = "rsa"
  type = "cer"
  value = "MIIC/DCCAeSgAwIBAgIBATANBgkqhkiG9w0BAQsFADAcMQ0wCwYDVQQDDAR0ZXN0MQswCQYDVQQGEwJVUzAeFw0yNDAzMDMyMjMyNDVaFw0yNTAzMDMyMjMyNDVaMBwxDTALBgNVBAMMBHRlc3QxCzAJBgNVBAYTAlVTMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAsZG3dFdTtCThjinAALmEQ6EbaDP6BTklqyzWsCZyFQjHND7fMPjImYouja2+IVMoFeI/ff7sqLhXAdrR3+uNQmvwKz8jyCupLSlX7nTKzvWATb/S1jXtduk85dPg6CrpuQCh7bBx/QLdC//gnHXbPPJVBZfUBVW5+7lGSdWWkxr85ERFlB7aSpu75hgjZU3IYmw6VTeqV1R4f22foaDozuEBZyivlgJNOiAb7dh0JjnGskY5HA6XWD7sRDqb7cvGRI7nwnssB0ZUOXqCViQ52jYszEXCNNbmMoL33tLMQmM8OERN8++cJaAPXCmrqIvpsyJ7IcVkHCoRG1MIwjRBmwIDAQABo0kwRzAOBgNVHQ8BAf8EBAMCB4AwFgYDVR0lAQH/BAwwCgYIKwYBBQUHAwMwHQYDVR0OBBYEFFX441UEohELyfMtHyoGvPo46bs+MA0GCSqGSIb3DQEBCwUAA4IBAQBe5diLmdVdD4UH3wXiHkzqgCrLZ6Xp9wu0qUuCko04Jg2vb9hB6l9CQ3opHBrOIXR0u8819+fGIjyFLTYQQu6t9ecSkokEYq6oSyaxYpZygNwzfaK6xlu45djb7oiyfj8mXt4sUgcaS6lZK9pJdFePfUvdJS2dLSWRB/qL+DuEbyRtIXNCYNn0JiH8i3DOCgWT+HZrQJud5Dw/sKBaxgLjJOvBjFZLY2sdP6JywHi4Mdpa1dyanOGNE7u2mj4LxAy6IOdq0VRqQMORVqhI0jRHOBh5N514rAR3vvLFK0WSNdUJBj3xgRsasDdLm3lnIxrhxPDbdDCNLPknnhxKR9Ss"
}
`

const testAccTrustFrameworkKeysetWithPfx = `
provider "azureadb2cief" {}
resource "azureadb2cief_trust_framework_key_set" "TestKeySetWithPfx" {
  name = "B2C_1A_TestKeySetWithPfx"
  kty = "rsa"
  type = "pfx"
  value = "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"
  password = "test"
}
`
