package client

type Client struct {
	TrustFrameworkPolicyClient *TrustFrameworkPolicyClient
	TrustFrameworkKeySetClient *TrustFrameworkKeySetClient
	Config                     *MsGraphClientConfig
}

func New(config MsGraphClientConfig) (*Client, error) {
	bc, err := newBaseClient(config)
	if err != nil {
		return nil, err
	}

	return &Client{
		TrustFrameworkPolicyClient: newPolicyClient(bc),
		TrustFrameworkKeySetClient: newKeySetClient(bc),
		Config:                     &config,
	}, nil
}
