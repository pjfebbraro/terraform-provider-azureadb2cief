# Azure Active Directory B2C Identity Experience Framework (b2cief) Provider

The Azure Provider can be used to configure custom policies in [Azure Active Directory B2C](https://docs.microsoft.com/en-us/azure/active-directory-b2c/custom-policy-overview) using the [Microsoft Graph Beta](https://docs.microsoft.com/en-us/graph/overview) API. Documentation regarding the [Data Sources](https://www.terraform.io/docs/language/data-sources/index.html) and [Resources](https://www.terraform.io/docs/language/resources/index.html) supported by the provider can be found in the navigation to the left.

## Example Usage

```hcl
# Configure Terraform
terraform {
  required_providers {
    azureadb2cief = {
      source  = "hashicorp/azureadb2cieif"
      version = "~> 0.1.0"
    }
  }
}

# Configure the Provider
provider "azuread" {
  tenant_id = "00000000-0000-0000-0000-000000000000"
}

```

## Authenticating to the Microsoft Graph API

Azure CLI and Service Principal with Client Id and Secret are Supported:

* Authenticating to Azure Active Directory using the Azure CLI
* [Authenticating to Azure Active Directory using a Service Principal and a Client Secret]


## Argument Reference

The following arguments are supported:
* `tenant_id` - (Required) The Tenant ID which should be used. This can also be sourced from the `ARM_TENANT_ID` environment variable.
---
When authenticating as a Service Principal using a Client Secret, the following fields can be set:
* `client_id` - (Optional) The Client ID which should be used when authenticating as a service principal. This can also be sourced from the `ARM_CLIENT_ID` environment variable.
* `client_secret` - (Optional) The application password to be used when authenticating using a client secret. This can also be sourced from the `ARM_CLIENT_SECRET` environment variable.

[More Information about setting up proper Microsoft Graph API Access can be found here](https://docs.microsoft.com/en-us/azure/active-directory-b2c/microsoft-graph-get-started?tabs=app-reg-ga)

---
For Azure CLI authentication, the following fields can be set:
* `use_cli` - (Optional) Should Azure CLI be used for authentication? This can also be sourced from the `ARM_USE_CLI` environment variable. Defaults to `true`.

Logging in with the CLI can be accomplished with the following:
```shell
az login --allow-no-subscriptions --output none --service-principal --tenant "${TENANT_ID}" --username "${CLIENT_ID}" --password "${CLIENT_SECRET}"
```
---

## Logging and Tracing

Logging output can be controlled with the `TF_LOG` or `TF_PROVIDER_LOG` environment variables. Exporting `TF_LOG=DEBUG` will increase the log verbosity and emit HTTP request and response traces to stdout when running Terraform. This output is very useful when reporting a bug in the provider.

Note that whilst we make every effort to remove authentication tokens from HTTP traces, they can still contain very identifiable and personal information which you should carefully censor before posting on our issue tracker.
