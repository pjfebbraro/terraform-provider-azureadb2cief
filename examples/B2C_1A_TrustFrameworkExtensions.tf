resource "azureadb2cief_trust_framework_policy" "B2C_1A_TrustFrameworkExtensions" {
  name   = "B2C_1A_TrustFrameworkExtensions"
  policy = <<-EOT
<?xml version="1.0" encoding="utf-8" ?>
<TrustFrameworkPolicy 
  xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" 
  xmlns:xsd="http://www.w3.org/2001/XMLSchema" 
  xmlns="http://schemas.microsoft.com/online/cpim/schemas/2013/06" 
  PolicySchemaVersion="0.3.0.0" 
  TenantId="${var.tenant_name}" 
  PolicyId="B2C_1A_TrustFrameworkExtensions" 
  PublicPolicyUri="http://${var.tenant_name}/B2C_1A_TrustFrameworkExtensions">
  
  <BasePolicy>
    <TenantId>${var.tenant_name}</TenantId>
    <PolicyId>${azureadb2cief_trust_framework_policy.B2C_1A_TrustFrameworkBase.name}</PolicyId>
  </BasePolicy>
  <BuildingBlocks>

  </BuildingBlocks>

  <ClaimsProviders>


    <ClaimsProvider>
      <DisplayName>Local Account SignIn</DisplayName>
      <TechnicalProfiles>
         <TechnicalProfile Id="login-NonInteractive">
          <Metadata>
            <Item Key="client_id">${azuread_application.ProxyIdentityExperienceFramework.application_id}</Item>
            <Item Key="IdTokenAudience">${azuread_application.IdentityExperienceFramework.application_id}</Item>
          </Metadata>
          <InputClaims>
            <InputClaim ClaimTypeReferenceId="client_id" DefaultValue="${azuread_application.ProxyIdentityExperienceFramework.application_id}" />
            <InputClaim ClaimTypeReferenceId="resource_id" PartnerClaimType="resource" DefaultValue="${azuread_application.IdentityExperienceFramework.application_id}" />
          </InputClaims>
        </TechnicalProfile>
      </TechnicalProfiles>
    </ClaimsProvider>

  </ClaimsProviders>

    <!--UserJourneys>
	
	</UserJourneys-->

</TrustFrameworkPolicy>

    EOT
}