package nmi

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-playground/form"

	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/common"
)

const (
	transactionEndpoint = "https://secure.networkmerchants.com/api/transact.php"
)

var (
	// assert client interface
	_ sleet.ClientWithContext = &NMIClient{}
)

// NMIClient represents an HTTP client and the associated authentication information required for making a Direct Post API request.
type NMIClient struct {
	testMode    bool
	securityKey string
	httpClient  *http.Client
}

// NewClient returns a new client for making NMI Direct Post API requests for a given merchant using a specified security key.
func NewClient(env common.Environment, securityKey string) *NMIClient {
	return NewWithHttpClient(env, securityKey, common.DefaultHttpClient())
}

// NewWithHttpClient returns a client for making NMI Direct Post API requests for a given merchant using a specified security key.
// The provided HTTP client will be used to make the requests.
func NewWithHttpClient(env common.Environment, securityKey string, httpClient *http.Client) *NMIClient {
	return &NMIClient{
		testMode:    nmiTestMode(env),
		securityKey: securityKey,
		httpClient:  httpClient,
	}
}

// Authorize makes a payment authorization request to NMI for the given payment details. If successful, the
// authorization response will be returned.
func (client *NMIClient) Authorize(request *sleet.AuthorizationRequest) (*sleet.AuthorizationResponse, error) {
	return client.AuthorizeWithContext(context.TODO(), request)
}

// AuthorizeWithContext makes a payment authorization request to NMI for the given payment details. If successful, the
// authorization response will be returned.
func (client *NMIClient) AuthorizeWithContext(ctx context.Context, request *sleet.AuthorizationRequest) (*sleet.AuthorizationResponse, error) {
	nmiAuthRequest := buildAuthRequest(client.testMode, client.securityKey, request)

	nmiResponse, httpResponse, err := client.sendRequest(ctx, nmiAuthRequest)
	if err != nil {
		return nil, err
	}

	responseHeader := sleet.GetHTTPResponseHeader(request.Options, *httpResponse)
	// "2" means declined and "3" means bad request
	if nmiResponse.Response != "1" {
		return &sleet.AuthorizationResponse{
			Success:    false,
			Response:   nmiResponse.ResponseCode,
			ErrorCode:  nmiResponse.ResponseCode,
			StatusCode: httpResponse.StatusCode,
			Header:     responseHeader,
		}, nil
	}

	return &sleet.AuthorizationResponse{
		Success:              true,
		TransactionReference: nmiResponse.TransactionID,
		AvsResult:            sleet.AVSResponseUnknown,
		CvvResult:            sleet.CVVResponseUnknown,
		Response:             nmiResponse.ResponseCode,
		AvsResultRaw:         nmiResponse.AVSResponseCode,
		CvvResultRaw:         nmiResponse.CVVResponseCode,
		StatusCode:           httpResponse.StatusCode,
		Header:               responseHeader,
	}, nil
}

// Capture captures an authorized payment through NMI. If successful, the capture response will be returned.
// Multiple captures cannot be made on the same authorization.
func (client *NMIClient) Capture(request *sleet.CaptureRequest) (*sleet.CaptureResponse, error) {
	return client.CaptureWithContext(context.TODO(), request)
}

// CaptureWithContext captures an authorized payment through NMI. If successful, the capture response will be returned.
// Multiple captures cannot be made on the same authorization.
func (client *NMIClient) CaptureWithContext(ctx context.Context, request *sleet.CaptureRequest) (*sleet.CaptureResponse, error) {
	nmiCaptureRequest := buildCaptureRequest(client.testMode, client.securityKey, request)

	nmiResponse, _, err := client.sendRequest(ctx, nmiCaptureRequest)
	if err != nil {
		return nil, err
	}

	// "2" means declined and "3" means bad request
	if nmiResponse.Response != "1" {
		return &sleet.CaptureResponse{
			Success: false,
			// transactionid is not always returned for bad captures, and, when it is, it's the id of the original transaction
			TransactionReference: request.TransactionReference,
			ErrorCode:            &nmiResponse.ResponseCode,
		}, nil
	}

	return &sleet.CaptureResponse{
		Success:              true,
		TransactionReference: nmiResponse.TransactionID,
	}, nil
}

// Void cancels a NMI transaction. If successful, the void response will be returned. A previously voided
// transaction or one that has already been settled cannot be voided.
func (client *NMIClient) Void(request *sleet.VoidRequest) (*sleet.VoidResponse, error) {
	return client.VoidWithContext(context.TODO(), request)
}

// VoidWithContext cancels a NMI transaction. If successful, the void response will be returned. A previously voided
// transaction or one that has already been settled cannot be voided.
func (client *NMIClient) VoidWithContext(ctx context.Context, request *sleet.VoidRequest) (*sleet.VoidResponse, error) {
	nmiVoidRequest := buildVoidRequest(client.testMode, client.securityKey, request)

	nmiResponse, _, err := client.sendRequest(ctx, nmiVoidRequest)
	if err != nil {
		return nil, err
	}

	// "2" means declined and "3" means bad request
	if nmiResponse.Response != "1" {
		return &sleet.VoidResponse{
			Success: false,
			// transactionid is not always returned for bad voids, and, when it is, it's the id of the original transaction
			TransactionReference: request.TransactionReference,
			ErrorCode:            &nmiResponse.ResponseCode,
		}, nil
	}

	return &sleet.VoidResponse{
		Success:              true,
		TransactionReference: nmiResponse.TransactionID,
	}, nil
}

// Refund refunds a NMI transaction that has been captured or settled.
// If successful, the refund response will be returned.
// Multiple refunds can be made on the same payment, but the total amount refunded should not exceed the payment total.
func (client *NMIClient) Refund(request *sleet.RefundRequest) (*sleet.RefundResponse, error) {
	return client.RefundWithContext(context.TODO(), request)
}

// RefundWithContext refunds a NMI transaction that has been captured or settled.
// If successful, the refund response will be returned.
// Multiple refunds can be made on the same payment, but the total amount refunded should not exceed the payment total.
func (client *NMIClient) RefundWithContext(ctx context.Context, request *sleet.RefundRequest) (*sleet.RefundResponse, error) {
	nmiRefundRequest := buildRefundRequest(client.testMode, client.securityKey, request)

	nmiResponse, _, err := client.sendRequest(ctx, nmiRefundRequest)
	if err != nil {
		return nil, err
	}

	// "2" means declined and "3" means bad request
	if nmiResponse.Response != "1" {
		return &sleet.RefundResponse{
			Success: false,
			// No transactionid is returned for unsuccessful refunds because refunds create new transactions
			TransactionReference: request.TransactionReference,
			ErrorCode:            &nmiResponse.ResponseCode,
		}, nil
	}

	return &sleet.RefundResponse{
		Success:              true,
		TransactionReference: nmiResponse.TransactionID,
	}, nil
}

// sendRequest sends an API request with the given payload to the NMI transaction endpoint.
// If the request is successfully sent, its response message will be returned.
func (client *NMIClient) sendRequest(ctx context.Context, data *Request) (*Response, *http.Response, error) {
	encoder := form.NewEncoder()
	formData, err := encoder.Encode(data)
	if err != nil {
		return nil, nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, transactionEndpoint, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, nil, err
	}

	parsedUrl, err := url.Parse(transactionEndpoint)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Add("Host", parsedUrl.Hostname())
	req.Header.Add("User-Agent", common.UserAgent())
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			// TODO log
		}
	}()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	parsedFormData, err := url.ParseQuery(string(respBody))
	if err != nil {
		return nil, nil, err
	}
	decoder := form.NewDecoder()
	nmiResponse := Response{}
	err = decoder.Decode(&nmiResponse, parsedFormData)
	if err != nil {
		return nil, nil, err
	}

	return &nmiResponse, resp, nil
}
