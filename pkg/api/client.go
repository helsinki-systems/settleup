package api

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	conf Config

	authToken string
}

func (c Client) makeURL(p string) string {
	return c.conf.APIConf.BaseURL + p
}

func (c *Client) Request(req *http.Request) ([]byte, error) {
	res, err := c.conf.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			c.conf.Logger.Warn(fmt.Sprintf("failed to close response body: %v", err))
		}
	}()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response status code %d: %s", res.StatusCode, body)
	}

	return body, nil
}

func (c *Client) AuthedRequest(req *http.Request) ([]byte, error) {
	if c.authToken == "" {
		return nil, errors.New("missing auth token")
	}

	// Include auth token, see https://firebase.google.com/docs/database/rest/auth#authenticate_with_an_id_token
	q := req.URL.Query()
	q.Add("auth", c.authToken)
	req.URL.RawQuery = q.Encode()

	return c.Request(req)
}

func New(c Config) *Client {
	return &Client{
		conf: c,
	}
}
