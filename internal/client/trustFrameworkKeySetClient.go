package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pjfebbraro/terraform-provider-azureadb2cief/internal/models"
	"io"
	"net/http"
)

type TrustFrameworkKeySetClient struct {
	*baseClient
}

func newKeySetClient(baseClient *baseClient) *TrustFrameworkKeySetClient {
	return &TrustFrameworkKeySetClient{
		baseClient,
	}
}

func (c *TrustFrameworkKeySetClient) GetKeySet(ctx context.Context, id string) (*models.TrustFrameworkKeySet, int, error) {
	path := fmt.Sprintf("/trustFramework/keySets/%s", id)
	response, err := c.doRequest(ctx, path, http.MethodGet, http.NoBody, nil)
	var status int

	if err != nil {
		return nil, status, err
	}
	status = response.StatusCode
	if status >= 400 {
		return nil, status, formatHttpErrorResponse(response)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, status, err
	}

	keyset := models.TrustFrameworkKeySet{}

	if err := json.Unmarshal(body, &keyset); err != nil {
		return nil, status, err
	}

	return &keyset, status, nil
}

func (c *TrustFrameworkKeySetClient) GetActiveKey(ctx context.Context, id string) (*models.TrustFrameworkKey, int, error) {
	path := fmt.Sprintf("/trustFramework/keySets/%s/getActiveKey", id)
	response, err := c.doRequest(ctx, path, http.MethodGet, http.NoBody, nil)

	var status int
	if err != nil {
		return nil, status, err
	}
	status = response.StatusCode
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, status, err
	}

	key := models.TrustFrameworkKey{}

	if err := json.Unmarshal(body, &key); err != nil {
		return nil, status, err
	}

	return &key, status, nil
}

func (c *TrustFrameworkKeySetClient) CreateKey(ctx context.Context, id string) (*models.TrustFrameworkKeySet, int, error) {
	path := "/trustFramework/keySets"

	keyId := map[string]string{
		"id": id,
	}
	var status int
	body, err := json.Marshal(keyId)

	if err != nil {
		return nil, status, err
	}
	contentType := "application/json"
	response, err := c.doRequest(ctx, path, http.MethodPost, bytes.NewBuffer(body), &contentType)

	if err != nil {
		return nil, status, err
	}
	status = response.StatusCode
	if status != 201 {
		return nil, status, formatHttpErrorResponse(response)
	}
	defer response.Body.Close()
	body, err = io.ReadAll(response.Body)
	if err != nil {
		return nil, status, err
	}

	keyset := models.TrustFrameworkKeySet{}
	err = json.Unmarshal(body, &keyset)
	if err != nil {
		return nil, status, err
	}

	return &keyset, status, nil
}

func (c *TrustFrameworkKeySetClient) GenerateKey(ctx context.Context, id string, key models.TrustFrameworkKey) (int, error) {
	path := fmt.Sprintf("/trustFramework/keySets/%s/generateKey", id)

	body, err := json.Marshal(key)
	var status int
	if err != nil {
		return status, err
	}
	contentType := "application/json"
	response, err := c.doRequest(ctx, path, http.MethodPost, bytes.NewBuffer(body), &contentType)
	if err != nil {
		return status, err
	}
	status = response.StatusCode
	if status != 200 {
		return status, formatHttpErrorResponse(response)
	}
	return status, nil
}
func (c *TrustFrameworkKeySetClient) UploadSecret(ctx context.Context, id string, key models.TrustFrameworkKey) (int, error) {
	path := fmt.Sprintf("/trustFramework/keySets/%s/uploadSecret", id)

	body, err := json.Marshal(key)
	var status int
	if err != nil {
		return status, err
	}
	contentType := "application/json"
	response, err := c.doRequest(ctx, path, http.MethodPost, bytes.NewBuffer(body), &contentType)
	if err != nil {
		return status, err
	}
	status = response.StatusCode
	if status != 200 {
		return status, formatHttpErrorResponse(response)
	}
	return status, nil
}
func (c *TrustFrameworkKeySetClient) UploadPFX(ctx context.Context, id string, key models.TrustFrameworkPfxKey) (int, error) {
	path := fmt.Sprintf("/trustFramework/keySets/%s/uploadPkcs12", id)

	body, err := json.Marshal(key)
	var status int
	if err != nil {
		return status, err
	}
	contentType := "application/json"
	response, err := c.doRequest(ctx, path, http.MethodPost, bytes.NewBuffer(body), &contentType)
	if err != nil {
		return status, err
	}
	status = response.StatusCode
	if status != 200 {
		return status, formatHttpErrorResponse(response)
	}
	return status, nil
}
func (c *TrustFrameworkKeySetClient) UploadCER(ctx context.Context, id string, key models.TrustFrameworkCerKey) (int, error) {
	path := fmt.Sprintf("/trustFramework/keySets/%s/uploadCertificate", id)

	body, err := json.Marshal(key)
	var status int
	if err != nil {
		return status, err
	}
	contentType := "application/json"
	response, err := c.doRequest(ctx, path, http.MethodPost, bytes.NewBuffer(body), &contentType)
	if err != nil {
		return status, err
	}
	status = response.StatusCode
	if status != 200 {
		return status, formatHttpErrorResponse(response)
	}
	return status, nil
}

func (c *TrustFrameworkKeySetClient) DeleteKey(ctx context.Context, id string) (int, error) {
	path := fmt.Sprintf("/trustFramework/keySets/%s", id)
	response, err := c.doRequest(ctx, path, http.MethodDelete, http.NoBody, nil)
	var status int
	if err != nil {
		return status, err
	}
	status = response.StatusCode
	if status != 204 {
		return status, formatHttpErrorResponse(response)
	}

	return status, nil
}
