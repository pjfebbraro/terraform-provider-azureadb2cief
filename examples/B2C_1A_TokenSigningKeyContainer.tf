resource "azureadb2cief_trust_framework_key_set" "B2C_1A_TokenSigningKeyContainer" {
  name = "B2C_1A_TokenSigningKeyContainer"
  use  = "sig"
  kty  = "RSA"
}