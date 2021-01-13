package flow

import (
	"github.com/fatih/structs"
	"github.com/json-iterator/go"
	"github.com/pkg/errors"
)

const (
	// RefundStatusCreated is a refund created and awaiting processing.
	RefundStatusCreated = "created"

	// RefundStatusAccepted is a refund accepted and ready to be transferred.
	RefundStatusAccepted = "accepted"

	// RefundStatusRejected is a refund that was not accepted.
	RefundStatusRejected = "rejected"

	// RefundStatusRefunded is a refund that was accepted and transferred.
	RefundStatusRefunded = "refunded"

	// RefundStatusCanceled is a refund that was canceled by one of the parties.
	RefundStatusCanceled = "canceled"
)

// Refund represents a request for a refund.
type Refund struct {
	// OrderID is the ID of the order being refunded.
	OrderID       string `structs:"refundCommerceOrder"`

	// ReceiverEmail is the email of the refunded payer.
	ReceiverEmail string `structs:"receiverEmail"`

	// Amount is the amount of money the refund is for.
	Amount        uint64 `structs:"amount"`

	// CallbackURL is where Flow will notify the server about the refund status.
	CallbackURL   string `structs:"urlCallBack"`
}

// RefundStatus is the data related to a refund.
type RefundStatus struct {
	// Token is an identifier for the refund.
	Token       string `json:"token"`

	// RefundOrder is the OrderID of the order being refunded.
	RefundOrder string `json:"flowRefundOrder"`

	// Date is the date on which the refund request was created. It format the format yyyy-mm-dd hh:mm:ss
	Date        string `json:"date"`

	// Status is the status of the refund. It might be one of:
	// created, accepted, rejected, refunded, canceled
	Status      string `json:"status"`

	// Amount is the amount of money being refunded.
	Amount      uint64 `json:"amount"`

	// Fee is the fee being charged for the refund.
	Fee         uint64 `json:"fee"`
}

// CreateRefund starts a new refund request.
func (c Client) CreateRefund(r Refund) (*RefundStatus, error) {
	url, body := c.buildPOST("/refund/create", structs.Map(r))

	data, err := c.post(url, body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to transact with the server")
	}

	var status RefundStatus
	err = jsoniter.Unmarshal(data, &status)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse response")
	}

	return &status, err
}

// CancelRefund cancels a refund request.
func (c Client) CancelRefund(token string) (*RefundStatus, error) {
	url, body := c.buildPOST("/refund/cancel", map[string]interface{}{
		"token": token,
	})

	data, err := c.post(url, body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to transact with the server")
	}

	var status RefundStatus
	err = jsoniter.Unmarshal(data, &status)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse response")
	}

	return &status, err
}

// GetRefundStatus fetches the status of a refund request.
func (c Client) GetRefundStatus(token string) (*RefundStatus, error) {
	url := c.buildGET("/refund/getStatus", map[string]interface{}{
		"token": token,
	})

	data, err := c.get(url)
	if err != nil {
		return nil, errors.Wrap(err, "unable to transact with the server")
	}

	var status RefundStatus
	err = jsoniter.Unmarshal(data, &status)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse response")
	}

	return &status, err
}