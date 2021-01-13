package flow

const (
	// OrderStatusAwaitingPayment indicates that the order's payment is still pending.
	OrderStatusAwaitingPayment = iota + 1

	// OrderStatusPayed indicates that the order has been successfully payed and should be accepted.
	OrderStatusPayed

	// OrderStatusRejected indicates that the payment failed because it was rejected.
	OrderStatusRejected

	// OrderStatusCanceled indicates that some party canceled the order before it was fulfilled.
	OrderStatusCanceled
)

// Order represents a payment order.
type Order struct {
	// FlowOrder is the ID provided by Flow.
	FlowOrder int `json:"flowOrder,omitempty"`

	// CommerceOrder is an optional ID created by the commerce.
	CommerceOrder string `json:"commerceOrder,omitempty"`

	// RequestDate is the date the order was created. If none is set, the current date is used. It follows the format
	// yyyy-mm-dd hh:mm:ss
	RequestDate string `json:"requestDate,omitempty"`

	// Status is the current status of an order it might be one of the following:
	//  1 Awaiting payment - OrderStatusAwaitingPayment
	//  2 Payed 		   - OrderStatusPayed
	//  3 Rejected		   - OrderStatusRejected
	//  4 Canceled         - OrderStatusCanceled
	Status int `json:"status,omitempty"`

	// Subject is the reason for the payment. It might be the items or service being bought.
	Subject string `json:"subject,omitempty"`

	// Currency is optionally set to define the currency of the transaction.
	Currency string `json:"currency,omitempty"`

	// Amount represents the amount of money being charged.
	Amount string `json:"amount,omitempty"`

	// PayerEmail is the email of the payer.
	PayerEmail string `json:"payer,omitempty"`

	// Optional contains optional fields like RUT and ID.
	Optional Optional `json:"optional,omitempty"`

	// PendingInfo reports if a media is still waiting on the payment, and the date since it's waiting.
	PendingInfo PendingInfo `json:"pending_info,omitempty"`

	// PaymentData contains additional information about the payment.
	PaymentData PaymentData `json:"paymentData,omitempty"`

	// MerchantID will be set if a merchant is being a middleman in the payment.
	MerchantID string `json:"merchantId,omitempty"`
}

// Optional contains optional fields like RUT and ID.
type Optional struct {
	// RUT or Rol Unico Tributario is a unique national identifier
	RUT string `json:"RUT,omitempty"`

	// ID is the ID of the transaction.
	ID  string `json:"ID,omitempty"`
}

// PendingInfo reports if a media is still waiting on the payment, and the date since it's waiting.
type PendingInfo struct {
	// Media refers to a payment entity.
	Media string `json:"media,omitempty"`

	// Date is the date since the payment is pending. It format the format yyyy-mm-dd hh:mm:ss
	Date  string `json:"date,omitempty"`
}

// PaymentData contains additional information about the payment.
type PaymentData struct {
	// Date is the date of the payment.
	Date string `json:"date,omitempty"`

	// Media refers to a payment entity.
	Media string `json:"media,omitempty"`

	// ConversionDate is the date of the conversion between currencies if more than one was involved in the payment.
	// The format is yyyy-mm-dd hh:mm:ss
	ConversionDate string `json:"conversionDate,omitempty"`

	// ConversionRate is the rate used to convert between currencies if more than one was involved in the payment.
	ConversionRate float64 `json:"conversionRate,omitempty"`

	// Amount is the payed amount.
	Amount string `json:"amount,omitempty"`

	// Currency is the payment in which the payment was made.
	Currency string `json:"currency,omitempty"`

	// Fee is the fee applied to the transaction.
	Fee string `json:"fee,omitempty"`

	// Balance is the Amount minus the Fee
	Balance int64 `json:"balance,omitempty"`

	// TransferDate is the date in which the transfer was made. It format the format yyyy-mm-dd hh:mm:ss
	TransferDate string `json:"transferDate,omitempty"`
}

