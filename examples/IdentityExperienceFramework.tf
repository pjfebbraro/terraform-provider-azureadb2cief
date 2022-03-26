resource "random_uuid" "AccessIEFScopeId" {}

resource "azuread_application" "IdentityExperienceFramework" {
  display_name                  = "IdentityExperienceFramework"
  oauth2_post_response_required = false
  owners = [
    data.azuread_client_config.current.object_id,
  ]
  sign_in_audience = "AzureADMyOrg"
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


