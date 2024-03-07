package server

import (
	"context"
	"errors"
	"fmt"

	"git.helsinki.tools/mittagessen-gmbh/settleup/pkg/api"
)

type (
	MembersAllInput struct {
		CommonInput
	}

	MembersAllOutput struct {
		// CommonOutput

		Body struct {
			Members []api.Member `json:"members"`
		}
	}
)

func makeMembersAllHandler(
	apiClient *api.Client,
) func(context.Context, *MembersAllInput) (*MembersAllOutput, error) {
	return func(ctx context.Context, _ *MembersAllInput) (*MembersAllOutput, error) {
		l := loggerFromRequest(ctx)

		ms, err := apiClient.AllMembersInGroup(ctx)
		if err != nil {
			l.Warn(fmt.Errorf("failed to query all members: %w", err).Error())
			return nil, errors.New("failed to query all members")
		}

		var out MembersAllOutput
		out.Body.Members = ms

		return &out, nil
	}
}
