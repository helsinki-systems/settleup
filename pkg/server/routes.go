package server

import (
	"git.helsinki.tools/mittagessen-gmbh/settleup/pkg/api"

	huma "github.com/danielgtaylor/huma/v2"
)

func registerRoutes(hapi huma.API, apiClient *api.Client) {
	huma.Get(hapi, "/members", makeMembersAllHandler(
		apiClient,
	))

	huma.Get(hapi, "/transactions", makeTransactionsAllHandler(
		apiClient,
	))
	huma.Post(hapi, "/transactions", makeTransactionsCreateHandler(
		apiClient,
	))
}
