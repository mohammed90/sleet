package request

const (
	AuthorizeOperation string = "AUTHORIZE"
	CaptureOperation          = "CAPTURE"
	VoidOperation             = "VOID"
	RefundOperation           = "REFUND"
)

type Base struct {
	APIOperation string `json:"apiOperation,omitempty"`
}
