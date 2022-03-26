package client

import (
	"bytes"
	"context"
	"fmt"
	"github.com/pjfebbraro/terraform-provider-azureadb2cief/internal/models"
	"io"
	"net/http"
)

type TrustFrameworkPolicyClient struct {
	*baseClient
}

func newPolicyClient(baseClient *baseClient) *TrustFrameworkPolicyClient {
	return &TrustFrameworkPolicyClient{
		baseClient,
	}
}

func (c *TrustFrameworkPolicyClient) Get(ctx context.Context, name string) (*models.Policy, int, error) {
	var status int
	path := fmt.Sprintf("/trustframework/policies/%s/$value", name)
	response, err := c.doRequest(ctx, path, http.MethodGet, http.NoBody, nil)
	if err != nil {
		return nil, status, err
	}
	status = response.StatusCode
	if status != http.StatusOK {
		return nil, status, formatHttpErrorResponse(response)
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, status, err
	}

	xml := string(body)

	return &models.Policy{
		Name:   name,
		Policy: xml,
	}, status, nil
}
func (c *TrustFrameworkPolicyClient) Create(ctx context.Context, policyXml *string) (int, error) {
	url := "/trustFramework/policies"
	var status int
	body := bytes.NewBuffer([]byte(*policyXml))
	contentType := "application/xml"
	response, err := c.doRequest(ctx, url, http.MethodPost, body, &contentType)
	if err != nil {
		return status, err
	}
	status = response.StatusCode
	if status != http.StatusCreated {
		return status, formatHttpErrorResponse(response)
	}

	return status, nil
}

func (c *TrustFrameworkPolicyClient) Update(ctx context.Context, policy *models.Policy) (int, error) {
	var status int
	path := fmt.Sprintf("/trustframework/policies/%s/$value", policy.Name)
	contentType := "application/xml"
	response, err := c.doRequest(ctx, path, http.MethodPut, bytes.NewBuffer([]byte(policy.Policy)), &contentType)
	if err != nil {
		return status, err
	}
	status = response.StatusCode
	if status != http.StatusCreated && status != http.StatusOK {
		return status, formatHttpErrorResponse(response)
	}

	return status, nil
}

func (c *TrustFrameworkPolicyClient) Delete(ctx context.Context, name string) (int, error) {
	path := fmt.Sprintf("/trustframework/policies/%s", name)
	var status int
	response, err := c.doRequest(ctx, path, http.MethodDelete, http.NoBody, nil)

	if err != nil {
		return status, err
	}
	status = response.StatusCode
	if status != http.StatusNoContent {
		return status, formatHttpErrorResponse(response)
	}

	return status, nil
}
