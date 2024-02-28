package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Based on https://github.com/firebase/firebase-admin-go/blob/b04387eff11f911cf10e9a76f4fcbf517b6c9a62/integration/auth/auth_test.go

const (
	loginURLFormat = "https://www.googleapis.com/identitytoolkit/v3/relyingparty/verifyPassword?key=%s"
)

func (c *Client) Login(ctx context.Context) (string, error) {
	reqBody, err := json.Marshal(map[string]any{
		"email":             c.conf.APIConf.SettleUpConf.Username,
		"password":          c.conf.APIConf.SettleUpConf.Password,
		"returnSecureToken": true,
	})
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %w", err)
	}

	// nosemgrep: gosec.G107-1
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf(loginURLFormat, c.conf.APIConf.Key),
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	body, err := c.Request(req)
	if err != nil {
		return "", fmt.Errorf("failed to do request: %w", err)
	}

	var authRes struct {
		IDToken string `json:"idToken"` //nolint:tagliatelle // The Firebase API sucks
	}
	if err := json.Unmarshal(body, &authRes); err != nil {
		return "", fmt.Errorf("failed to unmarshal response body %s: %w", string(body), err)
	}

	c.setAuthToken(authRes.IDToken)

	return authRes.IDToken, nil
}

func (c *Client) setAuthToken(token string) {
	c.authToken = token
}
