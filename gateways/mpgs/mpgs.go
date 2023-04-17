package mpgs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/BoltApp/sleet"
	"github.com/BoltApp/sleet/common"
	"github.com/BoltApp/sleet/gateways/mpgs/request/authorization"
	"github.com/BoltApp/sleet/gateways/mpgs/request/capture"
	"github.com/BoltApp/sleet/gateways/mpgs/request/refund"
	"github.com/BoltApp/sleet/gateways/mpgs/request/void"
	"github.com/BoltApp/sleet/gateways/mpgs/response"
)

var (
	_ sleet.Client            = &MPGSClient{}
	_ sleet.ClientWithContext = &MPGSClient{}
)

const version = "70"

type MPGSClient struct {
	mid        string
	password   string
	host       string
	version    string
	httpClient *http.Client
}

type mpgsRequest interface {
	authorization.Request | capture.Request | refund.Request | void.Request
}

func NewClient(mid, password, host string) *MPGSClient {
	return &MPGSClient{
		mid:        mid,
		password:   password,
		host:       host,
		version:    version,
		httpClient: common.DefaultHttpClient(),
	}
}

func NewWithHttpClient(mid, password, host string, client *http.Client) *MPGSClient {
	return &MPGSClient{
		mid:        mid,
		password:   password,
		host:       host,
		version:    version,
		httpClient: client,
	}
}

// Authorize implements sleet.ClientWithContext
func (mc *MPGSClient) Authorize(r *sleet.AuthorizationRequest) (*sleet.AuthorizationResponse, error) {
	return mc.AuthorizeWithContext(context.Background(), r)
}

// Capture implements sleet.ClientWithContext
func (mc *MPGSClient) Capture(r *sleet.CaptureRequest) (*sleet.CaptureResponse, error) {
	return mc.CaptureWithContext(context.Background(), r)
}

// Refund implements sleet.ClientWithContext
func (mc *MPGSClient) Refund(r *sleet.RefundRequest) (*sleet.RefundResponse, error) {
	return mc.RefundWithContext(context.Background(), r)
}

// Void implements sleet.ClientWithContext
func (mc *MPGSClient) Void(r *sleet.VoidRequest) (*sleet.VoidResponse, error) {
	return mc.VoidWithContext(context.Background(), r)
}

// AuthorizeWithContext implements sleet.ClientWithContext
func (mc *MPGSClient) AuthorizeWithContext(ctx context.Context, r *sleet.AuthorizationRequest) (*sleet.AuthorizationResponse, error) {
	mpgsReq := authorization.NewRequest(r)

	req, err := createHTTPRequest(ctx, mc, r.MerchantOrderReference, *r.ClientTransactionReference, mpgsReq)
	if err != nil {
		return nil, err
	}
	res, err := mc.doHTTPRequest(ctx, req)
	if err != nil {
		switch e := err.(type) {
		case *response.Error:
			return &sleet.AuthorizationResponse{Success: false, Response: e.Result, ErrorCode: e.Err.Cause}, err
		default:
			return nil, err
		}
	}

	authresponse := &sleet.AuthorizationResponse{
		Success:               true,
		TransactionReference:  res.txn.Transaction.ID,
		ExternalTransactionID: res.txn.Transaction.ID,
		Response:              string(res.txn.Result),
		Metadata: map[string]string{
			"receipt": res.txn.Transaction.Receipt,
		},
		StatusCode: res.statusCode,
		Header:     res.header,
	}

	return authresponse, nil
}

// CaptureWithContext implements sleet.ClientWithContext
func (mc *MPGSClient) CaptureWithContext(ctx context.Context, r *sleet.CaptureRequest) (*sleet.CaptureResponse, error) {
	mpgsReq := capture.NewRequest(r)

	req, err := createHTTPRequest(ctx, mc, *r.MerchantOrderReference, *r.ClientTransactionReference, mpgsReq)
	if err != nil {
		return nil, err
	}

	res, err := mc.doHTTPRequest(ctx, req)
	if err != nil {
		switch e := err.(type) {
		case *response.Error:
			return &sleet.CaptureResponse{Success: false, ErrorCode: common.SPtr(e.Err.Cause)}, err
		default:
			return nil, err
		}
	}

	return &sleet.CaptureResponse{
		Success:              true,
		TransactionReference: res.txn.Transaction.ID,
	}, nil
}

// RefundWithContext implements sleet.ClientWithContext
func (mc *MPGSClient) RefundWithContext(ctx context.Context, r *sleet.RefundRequest) (*sleet.RefundResponse, error) {
	mpgsReq := refund.NewRequest(r)

	req, err := createHTTPRequest(ctx, mc, *r.MerchantOrderReference, *r.ClientTransactionReference, mpgsReq)
	if err != nil {
		return nil, err
	}
	res, err := mc.doHTTPRequest(ctx, req)
	if err != nil {
		switch e := err.(type) {
		case *response.Error:
			return &sleet.RefundResponse{Success: false, ErrorCode: common.SPtr(e.Err.Cause)}, err
		default:
			return nil, err
		}
	}

	return &sleet.RefundResponse{
		Success:              true,
		TransactionReference: res.txn.Transaction.ID,
	}, nil
}

// VoidWithContext implements sleet.ClientWithContext
func (mc *MPGSClient) VoidWithContext(ctx context.Context, r *sleet.VoidRequest) (*sleet.VoidResponse, error) {
	mpgsReq := void.NewRequest(r)

	req, err := createHTTPRequest(ctx, mc, *r.MerchantOrderReference, *r.ClientTransactionReference, mpgsReq)
	if err != nil {
		return nil, err
	}
	res, err := mc.doHTTPRequest(ctx, req)
	if err != nil {
		switch e := err.(type) {
		case *response.Error:
			return &sleet.VoidResponse{Success: false, ErrorCode: common.SPtr(e.Err.Cause)}, err
		default:
			return nil, err
		}
	}

	return &sleet.VoidResponse{
		Success:              true,
		TransactionReference: res.txn.Transaction.ID,
	}, nil
}

func createHTTPRequest[MR mpgsRequest](ctx context.Context, client *MPGSClient, merchantOrderReference, clientTransactionReference string, r MR) (*http.Request, error) {
	body, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("https://%s/api/rest/version/%s/merchant/%s/order/%s/transaction/%s", client.host, client.version, client.mid, merchantOrderReference, clientTransactionReference)
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPut,
		url,
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(fmt.Sprintf("merchant.%s", client.mid), client.password)
	return req, nil
}

type apiResponse struct {
	txn        *response.RetrieveTransaction
	statusCode int
	header     http.Header
}

func (mc MPGSClient) doHTTPRequest(ctx context.Context, req *http.Request) (*apiResponse, error) {
	res, err := mc.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	if res.StatusCode != http.StatusCreated {
		errorRes := &response.Error{}
		decoder.Decode(errorRes)
		return nil, errorRes
	}
	txn := &response.RetrieveTransaction{}
	decoder.Decode(txn)
	return &apiResponse{
		txn:        txn,
		statusCode: 0,
		header:     res.Header,
	}, nil
}
