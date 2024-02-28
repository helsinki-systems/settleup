package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"git.helsinki.tools/mittagessen-gmbh/settleup/pkg/api"
)

func makeAllMembersHandler(
	l *slog.Logger,
	apiClient *api.Client,
) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		l := l
		rid := req.Header.Get(requestIDHeaderName)
		if rid != "" {
			if len(rid) != requestIDLength {
				l.Warn(fmt.Sprintf("invalid request ID length %d: %q", len(rid), rid))
			} else {
				l = l.With("rid", rid)
			}
		}

		ms, err := apiClient.AllMembersInGroup(req.Context())
		if err != nil {
			l.Warn(fmt.Sprintf("failed to query all members: %v", err))
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		b, err := json.MarshalIndent(ms, "", "\t")
		if err != nil {
			l.Warn(fmt.Sprintf("failed to marshal JSON: %v", err))
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		res.WriteHeader(http.StatusOK)
		res.Header().Add("Content-Type", "application/json")
		if _, err := res.Write(b); err != nil {
			l.Warn(fmt.Sprintf("failed to write response: %v", err))
		}
	}
}

func makeTransactionsHandler(
	l *slog.Logger,
	apiClient *api.Client,
) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet &&
			req.Method != http.MethodPost {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		l := l
		rid := req.Header.Get(requestIDHeaderName)
		if rid != "" {
			if len(rid) != requestIDLength {
				l.Warn(fmt.Sprintf("invalid request ID length %d: %q", len(rid), rid))
			} else {
				l = l.With("rid", rid)
			}
		}

		if req.Method == http.MethodGet {
			rs, err := apiClient.AllTransactionsInGroup(req.Context())
			if err != nil {
				l.Warn(fmt.Sprintf("failed to query all transactions: %v", err))
				res.WriteHeader(http.StatusInternalServerError)
				return
			}

			b, err := json.MarshalIndent(rs, "", "\t")
			if err != nil {
				l.Warn(fmt.Sprintf("failed to marshal JSON: %v", err))
				res.WriteHeader(http.StatusInternalServerError)
				return
			}

			res.WriteHeader(http.StatusOK)
			res.Header().Add("Content-Type", "application/json")
			if _, err := res.Write(b); err != nil {
				l.Warn(fmt.Sprintf("failed to write response: %v", err))
			}

			return
		}

		// Post

		body, err := io.ReadAll(req.Body)
		if err != nil {
			l.Warn(fmt.Sprintf("failed to read body: %v", err))
			res.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		var t api.Transaction
		if err := json.Unmarshal(body, &t); err != nil {
			l.Warn(fmt.Sprintf("failed to parse body: %v", err))
			res.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		if t.CurrencyCode == "" {
			t.CurrencyCode = api.CurrencyCodeEuro
		}
		if t.DateTime.IsZero() {
			t.DateTime = api.SettleUpTime{Time: time.Now()}
		}
		if t.Type == "" {
			t.Type = api.TransactionTypeExpense
		}

		tr, err := apiClient.CreateTransactionInGroup(req.Context(), t)
		if err != nil {
			l.Warn(fmt.Sprintf("failed to create transaction: %v", err))
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		b, err := json.MarshalIndent(tr, "", "\t")
		if err != nil {
			l.Warn(fmt.Sprintf("failed to marshal JSON: %v", err))
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		res.WriteHeader(http.StatusOK)
		res.Header().Add("Content-Type", "application/json")
		if _, err := res.Write(b); err != nil {
			l.Warn(fmt.Sprintf("failed to write response: %v", err))
		}
	}
}
