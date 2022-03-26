resource "azureadb2cief_trust_framework_key_set" "B2C_1A_TokenEncryptionKeyContainer" {
  name = "B2C_1A_TokenEncryptionKeyContainer"
  use  = "enc"
  kty  = "RSA"
}