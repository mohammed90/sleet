package authorization

import (
	"fmt"

	"github.com/BoltApp/sleet"
)

type SourceOfFunds struct {
	Type     string    `json:"type,omitempty"`
	Provided *Provided `json:"provided,omitempty"`
}

type Provided struct {
	Card *Card `json:"card,omitempty"`
}

type Card struct {
	Expiry        *Expiry         `json:"expiry,omitempty"`
	Number        string          `json:"number,omitempty"`
	NameOnCard    string          `json:"nameOnCard,omitempty"`
	SecurityCode  string          `json:"securityCode,omitempty"`
	StoredOnFile  CardOnFileState `json:"storedOnFile,omitempty"`
	DevicePayment *DevicePayment  `json:"devicePayment,omitempty"`
}

type Expiry struct {
	Month string `json:"month,omitempty"`
	Year  string `json:"year,omitempty"`
}

type DevicePayment struct {
	PaymentToken string `json:"paymentToken,omitempty"`
}

type CardOnFileState uint8

const (
	CardOnFileStateNotStored CardOnFileState = iota
	CardOnFileStateToBeStored
	CardOnFileStateStored
)

const (
	storedOnFileNotStored  string = "NOT_STORED"
	storedOnFileToBeStored string = "TO_BE_STORED"
	storedOnFileStored     string = "STORED"
)

var cofStateToString = map[CardOnFileState]string{
	CardOnFileStateNotStored:  storedOnFileNotStored,
	CardOnFileStateToBeStored: storedOnFileToBeStored,
	CardOnFileStateStored:     storedOnFileStored,
}
var stringToCofState = map[string]CardOnFileState{
	storedOnFileNotStored:  CardOnFileStateNotStored,
	storedOnFileToBeStored: CardOnFileStateToBeStored,
	storedOnFileStored:     CardOnFileStateStored,
}
var initiatorTypeToStoredOnFileState = map[sleet.ProcessingInitiatorType]CardOnFileState{
	sleet.ProcessingInitiatorTypeInitialCardOnFile:         CardOnFileStateToBeStored,
	sleet.ProcessingInitiatorTypeInitialRecurring:          CardOnFileStateToBeStored,
	sleet.ProcessingInitiatorTypeStoredCardholderInitiated: CardOnFileStateStored,
	sleet.ProcessingInitiatorTypeStoredMerchantInitiated:   CardOnFileStateStored,
	sleet.ProcessingInitiatorTypeFollowingRecurring:        CardOnFileStateStored,
}
var initiatorTypeToStoredOnFileString = map[sleet.ProcessingInitiatorType]string{
	sleet.ProcessingInitiatorTypeInitialCardOnFile:         storedOnFileToBeStored,
	sleet.ProcessingInitiatorTypeInitialRecurring:          storedOnFileToBeStored,
	sleet.ProcessingInitiatorTypeStoredCardholderInitiated: storedOnFileStored,
	sleet.ProcessingInitiatorTypeStoredMerchantInitiated:   storedOnFileStored,
	sleet.ProcessingInitiatorTypeFollowingRecurring:        storedOnFileStored,
}

func (cfs *CardOnFileState) UnmarshalText(text []byte) error {
	s := string(text)
	state, ok := stringToCofState[s]
	if !ok {
		return fmt.Errorf("unrecognized CoF state from MPGS: %s", s)
	}
	*cfs = state
	return nil
}

func (cfs CardOnFileState) MarshalText() ([]byte, error) {
	v := cofStateToString[cfs]
	return []byte(v), nil
}
