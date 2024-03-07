package server

import (
	"github.com/google/uuid"
)

type (
	// CommonOutput struct {
	// 	RequestID uuid.UUID `header:"X-Request-Id" readOnly:"true" hidden:"true"`
	// }

	CommonInput struct {
		RequestID uuid.UUID `header:"X-Request-Id" hidden:"true" readOnly:"true"`
	}
)
