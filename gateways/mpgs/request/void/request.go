package void

import (
	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/common"
	"github.com/BoltApp/sleet/gateways/mpgs/request"
	"github.com/google/uuid"
)

func NewRequest(r *sleet.VoidRequest) Request {
	if r.ClientTransactionReference == nil {
		r.ClientTransactionReference = common.SPtr(uuid.NewString())
	}
	return Request{
		Base: request.Base{
			APIOperation: request.VoidOperation,
		},
		Transaction: Transaction{
			TargetTransactionID: r.TransactionReference,
		},
	}
}

type Request struct {
	request.Base
	Transaction Transaction `json:"transaction,omitempty"`
}

type Transaction struct {
	TargetTransactionID string `json:"targetTransactionId,omitempty"`
}
