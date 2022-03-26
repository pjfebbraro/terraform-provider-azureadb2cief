package models

type TrustFrameworkKey struct {
	Exp *int    `json:"exp"`
	Kty *string `json:"kty"`
	Nbf *int    `json:"nbf"`
	Use *string `json:"use"`
	K   *string `json:"k"`
}

type TrustFrameworkKeySet struct {
	Id   *string             `json:"id"`
	Keys []TrustFrameworkKey `json:"keys"`
}
