# Terraform AzureADB2CIEF Example

This example demonstrates using the provider to manage the AzureAD B2C Custom Policy Starter Pack

1. `B2C_1A_TokenSigningKeyContainer` required token signing key container
2. `B2C_1A_TokenEncryptionKeyContainer` required token encryption key container
3. `IdentityExperienceFramework` required `azuread_application` referenced by custom policies
4. `ProxyIdentityExpreienceFramework` required `azuread_application` referenced by custom policies
5. `B2C_1A_TrustFrameworkBase`
6. `B2C_1A_TrustFrameworkExtensions`
7. `B2C_1A_signup_signin`
