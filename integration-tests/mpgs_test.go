package test

import (
	"strings"
	"testing"

	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/gateways/mpgs"
	sleet_testing "github.com/BoltApp/sleet/testing"
	"github.com/google/uuid"
)

var (
	mid      string = getEnv("MPGS_MID")
	host     string = getEnv("MPGS_HOST")
	password string = getEnv("MPGS_PASSWORD")
)

// TestMPGSAuthorizeFailed
//
// MPGS has test cards here: https://ap-gateway.mastercard.com/docs/testing#cards-responses
// Using a rejected card number
func TestMPGSAuthorizeFailed(t *testing.T) {
	client := mpgs.NewClient(mid, password, host)
	failedRequest := sleet_testing.BaseAuthorizationRequest()
	// set ClientTransactionReference to be empty
	failedRequest.CreditCard.Number = "5111111214111118"
	_, err := client.Authorize(failedRequest)
	if err == nil {
		t.Error("Authorize request should have failed")
	}

	if !strings.Contains(strings.ToLower(err.Error()), "invalid card number") {
		t.Errorf("Response should contain 'expired': %s", err.Error())
	}
}

// TestMPGSAuth
//
// This should successfully create an authorization on MPGS
func TestMPGSAuth(t *testing.T) {
	client := mpgs.NewClient(mid, password, host)
	request := sleet_testing.BaseAuthorizationRequest()
	request.MerchantOrderReference = uuid.NewString()
	request.CreditCard = &sleet.CreditCard{
		FirstName:       "John",
		LastName:        "Doe",
		Number:          "5111111111111118",
		ExpirationMonth: 1,
		ExpirationYear:  2039,
		CVV:             "100",
		Save:            false,
	}

	auth, err := client.Authorize(request)
	if err != nil {
		t.Errorf("Authorize request should not have failed: %s", err)
	}

	if !auth.Success {
		t.Error("Resulting auth should have been successful")
	}
}

// TestMPGSAuthFullCapture
//
// This should successfully create an authorization on MPGS then Capture for full amount
func TestMPGSAuthFullCapture(t *testing.T) {
	client := mpgs.NewClient(mid, password, host)
	authRequest := sleet_testing.BaseAuthorizationRequest()
	authRequest.MerchantOrderReference = uuid.NewString()
	authRequest.CreditCard = &sleet.CreditCard{
		FirstName:       "John",
		LastName:        "Doe",
		Number:          "5111111111111118",
		ExpirationMonth: 1,
		ExpirationYear:  2039,
		CVV:             "100",
		Save:            false,
	}

	auth, err := client.Authorize(authRequest)
	if err != nil {
		t.Error("Authorize request should not have failed")
	}

	if !auth.Success {
		t.Error("Resulting auth should have been successful")
	}

	captureRequest := &sleet.CaptureRequest{
		Amount:                 &authRequest.Amount,
		TransactionReference:   auth.TransactionReference,
		MerchantOrderReference: &authRequest.MerchantOrderReference,
	}
	capture, err := client.Capture(captureRequest)
	if err != nil {
		t.Error("Capture request should not have failed")
	}

	if !capture.Success {
		t.Error("Resulting capture should have been successful")
	}
}

// TestMPGSAuthPartialCapture
//
// This should successfully create an authorization on MPGS then Capture for a partial amount
// Since we auth for 1.00USD, we will capture for $0.50
func TestMPGSAuthPartialCapture(t *testing.T) {
	client := mpgs.NewClient(mid, password, host)
	authRequest := sleet_testing.BaseAuthorizationRequest()
	auth, err := client.Authorize(authRequest)
	if err != nil {
		t.Error("Authorize request should not have failed")
	}

	if !auth.Success {
		t.Error("Resulting auth should have been successful")
	}

	captureRequest := &sleet.CaptureRequest{
		Amount: &sleet.Amount{
			Amount:   50,
			Currency: "USD",
		}, TransactionReference: auth.TransactionReference,
	}
	capture, err := client.Capture(captureRequest)
	if err != nil {
		t.Error("Capture request should not have failed")
	}

	if !capture.Success {
		t.Error("Resulting capture should have been successful")
	}
}

// TestMPGSAuthVoid
//
// This should successfully create an authorization on MPGS then Void/Cancel the Auth
func TestMPGSAuthVoid(t *testing.T) {
	client := mpgs.NewClient(mid, password, host)
	authRequest := sleet_testing.BaseAuthorizationRequest()
	auth, err := client.Authorize(authRequest)
	if err != nil {
		t.Error("Authorize request should not have failed")
	}

	if !auth.Success {
		t.Error("Resulting auth should have been successful")
	}

	voidRequest := &sleet.VoidRequest{
		TransactionReference: auth.TransactionReference,
	}
	void, err := client.Void(voidRequest)
	if err != nil {
		t.Error("Void request should not have failed")
	}

	if !void.Success {
		t.Error("Resulting void should have been successful")
	}
}

// TestMPGSAuthCaptureRefund
//
// This should successfully create an authorization on MPGS then Capture for full amount, then refund for full amount
func TestMPGSAuthCaptureRefund(t *testing.T) {
	client := mpgs.NewClient(mid, password, host)
	authRequest := sleet_testing.BaseAuthorizationRequest()
	auth, err := client.Authorize(authRequest)
	if err != nil {
		t.Error("Authorize request should not have failed")
	}

	if !auth.Success {
		t.Error("Resulting auth should have been successful")
	}

	captureRequest := &sleet.CaptureRequest{
		Amount:               &authRequest.Amount,
		TransactionReference: auth.TransactionReference,
	}
	capture, err := client.Capture(captureRequest)
	if err != nil {
		t.Error("Capture request should not have failed")
	}

	if !capture.Success {
		t.Error("Resulting capture should have been successful")
	}

	refundRequest := &sleet.RefundRequest{
		Amount:               &authRequest.Amount,
		TransactionReference: capture.TransactionReference,
	}

	refund, err := client.Refund(refundRequest)
	if err != nil {
		t.Error("Refund request should not have failed")
	}

	if !refund.Success {
		t.Error("Resulting refund should have been successful")
	}
}
