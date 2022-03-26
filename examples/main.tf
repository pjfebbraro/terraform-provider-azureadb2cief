terraform {
  required_providers {
    # azuread = {
    #   source = "hashicorp/azuread"
    #   version = "~> 2.19"
    # }
    azureadb2cief = {
      source  = "local/providers/azureadb2c"
      version = "~> 1.0.3"
    }
    # random = {
    #   source = "hashicorp/random"
    #   version = ">= 3.1.0"
    # }
  }
}
provider "azuread" {}
provider "azureadb2cief" {
  tenant_id = var.tenant_id
  use_cli   = true
}
provider "random" {}

variable "tenant_name" {
  type = string
}
variable "tenant_id" {
  type = string
}
data "azureadb2cief_client_config" "current" {}
data "azuread_client_config" "current" {}
data "azuread_service_principal" "MicrosoftGraph" {
  display_name = "Microsoft Graph"
}
