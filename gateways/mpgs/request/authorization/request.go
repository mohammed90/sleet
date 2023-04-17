package authorization

import (
	"fmt"

	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/common"
	"github.com/BoltApp/sleet/gateways/mpgs/request"
	"github.com/google/uuid"
)

func NewRequest(r *sleet.AuthorizationRequest) Request {
	if r.MerchantOrderReference == "" {
		r.MerchantOrderReference = uuid.NewString()
	}
	if r.ClientTransactionReference == nil {
		r.ClientTransactionReference = common.SPtr(uuid.NewString())
	}
	mpgsReq := Request{
		Base: request.Base{APIOperation: request.AuthorizeOperation},
		Order: Order{
			Amount:    float32(r.Amount.Amount) / 100,
			Currency:  r.Amount.Currency,
			Reference: r.MerchantOrderReference,
		},
		SourceOfFunds: SourceOfFunds{
			Type: "CARD",
		},
	}
	if r.CreditCard != nil {
		mpgsReq.SourceOfFunds.Provided = &Provided{
			Card: &Card{
				Expiry: &Expiry{
					Month: fmt.Sprintf("%0d", r.CreditCard.ExpirationMonth),
					Year:  fmt.Sprintf("%d", r.CreditCard.ExpirationYear-2000),
				},
				Number:       r.CreditCard.Number,
				NameOnCard:   fmt.Sprintf("%s %s", r.CreditCard.FirstName, r.CreditCard.LastName),
				SecurityCode: r.CreditCard.CVV,
			},
		}
		if r.CreditCard.Save {
			mpgsReq.SourceOfFunds.Provided.Card.StoredOnFile = CardOnFileStateToBeStored
		}
	}
	if r.Cryptogram != "" {
		mpgsReq.Order.WalletProvider = "APPLE_PAY"
		mpgsReq.SourceOfFunds.Provided.Card.DevicePayment.PaymentToken = r.Cryptogram
	}

	if r.ThreeDS != nil {
		mpgsReq.Authentication = &request.Authentication{
			ThreeDS: &request.ThreeDS{
				ACSECI:              r.ECI,
				AuthenticationToken: r.ThreeDS.CAVV,
				TransactionID:       r.ThreeDS.DSTransactionID,
			},
			ThreeDS2: &request.ThreeDS2{
				ProtocolVersion:   r.ThreeDS.Version,
				StatusReasonCode:  "", // left empty
				TransactionStatus: r.ThreeDS.PAResStatus,
			},
		}
	}
	if r.ProcessingInitiator != nil {
		mpgsReq.Transaction = &Transaction{
			Source: initiatorTypeToTxnSource[*r.ProcessingInitiator],
		}
		mpgsReq.SourceOfFunds.Provided.Card.StoredOnFile = initiatorTypeToStoredOnFileState[*r.ProcessingInitiator]
	}
	return mpgsReq
}

type Request struct {
	request.Base
	Authentication *request.Authentication `json:"authentication,omitempty"`
	Order          Order                   `json:"order,omitempty"`
	SourceOfFunds  SourceOfFunds           `json:"sourceOfFunds,omitempty"`
	Transaction    *Transaction            `json:"transaction,omitempty"`
}
type Transaction struct {
	Source TransactionSource `json:"source,omitempty"`
}

type Order struct {
	Amount   float32 `json:"amount,omitempty"`
	Currency string  `json:"currency,omitempty"`

	// The identifier of the order
	Reference string `json:"reference,omitempty"`

	WalletProvider string `json:"walletProvider,omitempty"`
}
