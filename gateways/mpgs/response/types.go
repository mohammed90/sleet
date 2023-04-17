package response

type Result string

const (
	ResultFailure Result = "FAILURE"
	ResultPending        = "PENDING"
	ResultSuccess        = "SUCCESS"
	ResultUnknown        = "UNKNOWN"
)

type GatewayCode string

const (
	GatewayCodeAborted                     GatewayCode = "ABORTED"
	GatewayCodeAcquirerSystemError                     = "ACQUIRER_SYSTEM_ERROR"
	GatewayCodeApproved                                = "APPROVED"
	GatewayCodeApprovedAuto                            = "APPROVED_AUTO"
	GatewayCodeApprovedPendingSettlement               = "APPROVED_PENDING_SETTLEMENT"
	GatewayCodeAuthenticationFailed                    = "AUTHENTICATION_FAILED"
	GatewayCodeAuthenticationInProgress                = "AUTHENTICATION_IN_PROGRESS"
	GatewayCodeBalanceAvailable                        = "BALANCE_AVAILABLE"
	GatewayCodeBalanceUnknown                          = "BALANCE_UNKNOWN"
	GatewayCodeBlocked                                 = "BLOCKED"
	GatewayCodeCancelled                               = "CANCELLED"
	GatewayCodeDeclined                                = "DECLINED"
	GatewayCodeDeclinedAVS                             = "DECLINED_AVS"
	GatewayCodeDeclinedAVSCSC                          = "DECLINED_AVS_CSC"
	GatewayCodeDeclinedCSC                             = "DECLINED_CSC"
	GatewayCodeDeclinedDoNotContact                    = "DECLINED_DO_NOT_CONTACT"
	GatewayCodeDeclinedPaymentPlan                     = "DECLINED_PAYMENT_PLAN"
	GatewayCodeDeferredTransactionReceived             = "DEFERRED_TRANSACTION_RECEIVED"
	GatewayCodeExceedRetryLimit                        = "EXCEEDED_RETRY_LIMIT"
	GatewayCodeExpiredCard                             = "EXPIRED_CARD"
	GatewayCodeInsufficientFunds                       = "INSUFFICIENT_FUNDS"
	GatewayCodeInvalidCSC                              = "INVALID_CSC"
	GatewayCodeNotSupported                            = "NOT_SUPPORTED"
	GatewayCodePartiallyApproved                       = "PARTIALLY_APPROVED"
	GatewayCodePending                                 = "PENDING"
	GatewayCodeReferred                                = "REFERRED"
	GatewayCodeSystemError                             = "SYSTEM_ERROR"
	GatewayCodeTimedOut                                = "TIMED_OUT"
	GatewayCodeUnknown                                 = "UNKNOWN"
	GatewayCodeUnspecifiedFailure                      = "UNSPECIFIED_FAILURE"
)
