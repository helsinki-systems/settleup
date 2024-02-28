//nolint:tagliatelle // The Settle Up API sucks
package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// See https://docs.google.com/document/d/18mxnyYSm39cbceA2FxFLiOfyyanaBY6ogG7oscgghxU/edit

type Member struct {
	ID *string `json:"id,omitempty"`

	Name          string  `json:"name"`
	Active        bool    `json:"active"`
	DefaultWeight float32 `json:"defaultWeight,string"`
}

func (c *Client) AllMembersInGroup(ctx context.Context) ([]Member, error) {
	// nosemgrep: gosec.G107-1
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		c.makeURL("/members/"+c.conf.APIConf.SettleUpConf.GroupID+".json"), //nolint:goconst // We want to be explicit here
		http.NoBody,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	body, err := c.AuthedRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %w", err)
	}

	var membersRes map[string]Member
	if err := json.Unmarshal(body, &membersRes); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body %s: %w", string(body), err)
	}

	ms := make([]Member, 0, len(membersRes))
	for k, v := range membersRes {
		k := k
		v.ID = &k
		ms = append(ms, v)
	}

	return ms, nil
}

type currencyCode string

const (
	CurrencyCodeEuro currencyCode = "EUR"
)

type transactionType string

const (
	TransactionTypeExpense transactionType = "expense"
)

type Item struct {
	Amount  float32      `json:"amount,string"`
	ForWhom []PayingUser `json:"forWhom"`
}

type PayingUser struct {
	MemberID string  `json:"memberId"`
	Weight   float32 `json:"weight,string"`
}

type SettleUpTime struct {
	time.Time
}

func (sut *SettleUpTime) UnmarshalJSON(b []byte) error {
	i, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse time %s: %w", b, err)
	}
	const msPerS = 1000
	sut.Time = time.Unix(i/msPerS, 0)
	return nil
}

//nolint:unparam // We need to return an error here to satisfy the json.Marshaler interface
func (sut SettleUpTime) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(sut.Time.UnixMilli(), 10)), nil
}

type Transaction struct {
	ID *string `json:"id,omitempty"`

	Category          *string         `json:"category,omitempty"`
	CurrencyCode      currencyCode    `json:"currencyCode"`
	DateTime          SettleUpTime    `json:"dateTime"`
	ExchangeRates     any             `json:"exchangeRates"`
	FixedExchangeRate bool            `json:"fixedExchangeRate"`
	Items             []Item          `json:"items"`
	Purpose           string          `json:"purpose"`
	ReceiptURL        *string         `json:"receiptUrl,omitempty"`
	Type              transactionType `json:"type"`
	WhoPaid           []PayingUser    `json:"whoPaid"`
}

func (c *Client) AllTransactionsInGroup(ctx context.Context) ([]Transaction, error) {
	// nosemgrep: gosec.G107-1
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		c.makeURL("/transactions/"+c.conf.APIConf.SettleUpConf.GroupID+".json"),
		http.NoBody,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	body, err := c.AuthedRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %w", err)
	}

	var transactionsRes map[string]Transaction
	if err := json.Unmarshal(body, &transactionsRes); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body %s: %w", string(body), err)
	}

	ts := make([]Transaction, 0, len(transactionsRes))
	for k, v := range transactionsRes {
		k := k
		v.ID = &k
		ts = append(ts, v)
	}

	return ts, nil
}

type CreateTransactionRes struct {
	ID string `json:"name"`
}

func (c *Client) CreateTransactionInGroup(ctx context.Context, t Transaction) (*CreateTransactionRes, error) {
	reqBody, err := json.Marshal(t)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	// nosemgrep: gosec.G107-1
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.makeURL("/transactions/"+c.conf.APIConf.SettleUpConf.GroupID+".json"),
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	body, err := c.AuthedRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %w", err)
	}

	var transactionRes CreateTransactionRes
	if err := json.Unmarshal(body, &transactionRes); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body %s: %w", string(body), err)
	}

	c.conf.Logger.Info("new transaction", "transaction", transactionRes)

	return &transactionRes, nil
}
