package authorization

import (
	"fmt"

	"github.com/BoltApp/sleet"
)

type TransactionSource uint8

const (
	TransactionSourceInternet TransactionSource = iota
	TransactionSourceMerchant
)

const (
	transactionSourceInternetString string = "INTERNET"
	transactionSourceMerchantString string = "MERCHANT"
)

var txnSourceToString = map[TransactionSource]string{
	TransactionSourceInternet: transactionSourceInternetString,
	TransactionSourceMerchant: transactionSourceMerchantString,
}
var stringToTxnSource = map[string]TransactionSource{
	transactionSourceInternetString: TransactionSourceInternet,
	transactionSourceMerchantString: TransactionSourceMerchant,
}

func (cfs *TransactionSource) UnmarshalText(text []byte) error {
	s := string(text)
	state, ok := stringToTxnSource[s]
	if !ok {
		return fmt.Errorf("unrecognized CoF state from MPGS: %s", s)
	}
	*cfs = state
	return nil
}

func (cfs TransactionSource) MarshalText() ([]byte, error) {
	v := txnSourceToString[cfs]
	return []byte(v), nil
}

var initiatorTypeToTxnSource = map[sleet.ProcessingInitiatorType]TransactionSource{
	sleet.ProcessingInitiatorTypeInitialCardOnFile:         TransactionSourceInternet,
	sleet.ProcessingInitiatorTypeInitialRecurring:          TransactionSourceInternet,
	sleet.ProcessingInitiatorTypeStoredCardholderInitiated: TransactionSourceInternet,
	sleet.ProcessingInitiatorTypeStoredMerchantInitiated:   TransactionSourceMerchant,
	sleet.ProcessingInitiatorTypeFollowingRecurring:        TransactionSourceMerchant,
}
