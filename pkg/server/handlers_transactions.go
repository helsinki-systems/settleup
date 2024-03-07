package server

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/helsinki-systems/settleup/pkg/api"
)

type (
	TransactionsAllInput struct {
		CommonInput
	}

	TransactionsAllOutput struct {
		// CommonOutput

		Body struct {
			Transactions []api.Transaction `json:"transactions"`
		}
	}
)

func makeTransactionsAllHandler(
	apiClient *api.Client,
) func(context.Context, *TransactionsAllInput) (*TransactionsAllOutput, error) {
	return func(ctx context.Context, _ *TransactionsAllInput) (*TransactionsAllOutput, error) {
		l := loggerFromRequest(ctx)

		ts, err := apiClient.AllTransactionsInGroup(ctx)
		if err != nil {
			l.Warn(fmt.Errorf("failed to query all transactions: %w", err).Error())
			return nil, errors.New("failed to query all transactions")
		}

		var out TransactionsAllOutput
		out.Body.Transactions = ts

		return &out, nil
	}
}

type (
	TransactionPayingUser struct {
		MemberID string  `json:"member_id"`
		Weight   float32 `json:"weight"`
	}

	TransactionsCreateInput struct {
		CommonInput

		Body struct {
			Purpose   string    `json:"purpose"`
			CreatedAt time.Time `json:"created_at"`

			What struct {
				Amount  float32                 `json:"amount"`
				ForWhom []TransactionPayingUser `json:"for_whom"`
			} `json:"what"`

			WhoPaid []TransactionPayingUser `json:"who_paid"`
		}
	}

	TransactionsCreateOutput struct {
		// CommonOutput

		Body struct {
			ID string `json:"id"`
		}
	}
)

func makeTransactionsCreateHandler(
	apiClient *api.Client,
) func(context.Context, *TransactionsCreateInput) (*TransactionsCreateOutput, error) {
	return func(ctx context.Context, input *TransactionsCreateInput) (*TransactionsCreateOutput, error) {
		l := loggerFromRequest(ctx)

		in := input.Body
		var forWhom []api.PayingUser
		for _, fm := range in.What.ForWhom {
			forWhom = append(forWhom, api.PayingUser(fm))
		}
		var whoPaid []api.PayingUser
		for _, wp := range in.WhoPaid {
			whoPaid = append(whoPaid, api.PayingUser(wp))
		}
		tr := api.Transaction{
			Purpose: in.Purpose,
			DateTime: api.SettleUpTime{
				Time: in.CreatedAt,
			},

			Items: []api.Item{
				{
					Amount:  in.What.Amount,
					ForWhom: forWhom,
				},
			},

			WhoPaid: whoPaid,

			// Hardcoded fields
			CurrencyCode: api.CurrencyCodeEuro,
			Type:         api.TransactionTypeExpense,
		}

		t, err := apiClient.CreateTransactionInGroup(ctx, tr)
		if err != nil {
			l.Warn(fmt.Errorf("failed to create transaction: %w", err).Error())
			return nil, errors.New("failed to create transaction")
		}

		var out TransactionsCreateOutput
		out.Body.ID = t.ID

		return &out, nil
	}
}
