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
