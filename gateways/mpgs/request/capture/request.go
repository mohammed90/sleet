package capture

import (
	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/common"
	"github.com/BoltApp/sleet/gateways/mpgs/request"
	"github.com/google/uuid"
)

func NewRequest(r *sleet.CaptureRequest) Request {
	if r.ClientTransactionReference == nil {
		r.ClientTransactionReference = common.SPtr(uuid.NewString())
	}
	return Request{
		Base: request.Base{APIOperation: request.CaptureOperation},
		Transaction: Transaction{
			Amount:   float32(r.Amount.Amount) / 100,
			Currency: r.Amount.Currency,
		},
	}
}

type Request struct {
	request.Base
	Transaction Transaction `json:"transaction,omitempty"`
}

type Transaction struct {
	Amount   float32 `json:"amount,omitempty"`
	Currency string  `json:"currency,omitempty"`
}
