
resource "azuread_application" "ProxyIdentityExperienceFramework" {
  display_name                   = "ProxyIdentityExperienceFramework"
  fallback_public_client_enabled = true
  owners = [
    data.azuread_client_config.current.object_id,
  ]
  sign_in_audience = "AzureADMyOrg"
  tags             = []

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
      id   = azuread_application.IdentityExperienceFramework.oauth2_permission_scope_ids["user_impersonation"]
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