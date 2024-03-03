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

type TrustFrameworkPfxKey struct {
	Key      *string `json:"key"`
	Password *string `json:"password"`
}
type TrustFrameworkCerKey struct {
	Key *string `json:"key"`
}
