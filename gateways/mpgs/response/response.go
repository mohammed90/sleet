package response

import "time"

type RetrieveTransaction struct {
	Merchant      string        `json:"merchant"`
	Result        Result        `json:"result"`
	Order         Order         `json:"order"`
	Response      Response      `json:"response"`
	SourceOfFunds SourceOfFunds `json:"sourceOfFunds,omitempty"`
	Transaction   Transaction   `json:"transaction"`
}

type Order struct {
	ID              string    `json:"id"`
	Amount          float32   `json:"amount"`
	CreationTime    time.Time `json:"creationTime"`
	Currency        string    `json:"currency"`
	LastUpdatedTime time.Time `json:"lastUpdatedTime"`

	MerchantAmount        float32 `json:"merchantAmount"`
	MerchantCurrency      string  `json:"merchantCurrency"`
	TotalAuthorizedAmount float32 `json:"totalAuthorizedAmount"`
	TotalCapturedAmount   float32 `json:"totalCapturedAmount"`
	TotalDisbursedAmount  float32 `json:"totalDisbursedAmount"`
	TotalRefundedAmount   float32 `json:"totalRefundedAmount"`
}

type Response struct {
	GatewayCode GatewayCode `json:"gatewayCode"`
}

type SourceOfFunds struct {
	Type     string    `json:"type,omitempty"`
	Provided *Provided `json:"provided,omitempty"`
}

type Provided struct {
	Card *Card `json:"card,omitempty"`
}

type Card struct {
	Brand         string         `json:"brand,omitempty"`
	Scheme        string         `json:"scheme,omitempty"`
	StoredOnFile  string         `json:"storedOnFile,omitempty"`
	FundingMethod string         `json:"fundingMethod,omitempty"`
	Expiry        *Expiry        `json:"expiry,omitempty"`
	Number        string         `json:"number,omitempty"`
	NameOnCard    string         `json:"nameOnCard,omitempty"`
	SecurityCode  string         `json:"securityCode,omitempty"`
	DevicePayment *DevicePayment `json:"devicePayment,omitempty"`
}

type Expiry struct {
	Month string `json:"month,omitempty"`
	Year  string `json:"year,omitempty"`
}

type DevicePayment struct {
	CryptogramFormat string `json:"cryptogramFormat,omitempty"`
}

type Transaction struct {
	ID       string  `json:"id"`
	Type     string  `json:"type"`
	Amount   float32 `json:"amount"`
	Currency string  `json:"currency"`
	Receipt  string  `json:"receipt,omitempty"`
}
