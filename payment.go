package flow

import (
	"fmt"
	"github.com/fatih/structs"
	"github.com/json-iterator/go"
	"github.com/pkg/errors"
)


// OrderRequest is the data needed to start a payment order.
type OrderRequest struct {
	// CommerceOrder is an optional ID created by the commerce.
	CommerceOrder   string                 `structs:"commerceOrder"`

	// Subject is the reason for the payment. It might be the items or service being bought.
	Subject         string                 `structs:"subject"`

	// Currency is optionally set to define the currency of the transaction.
	Currency        string                 `structs:"currency,omitempty"`

	// Amount represents the amount of money being charged.
	Amount          uint64                 `structs:"amount"`

	// PayerEmail is the email of the payer.
	PayerEmail      string                 `structs:"email"`

	// PaymentMethod can be set to define the allowed payment methods. It default to 9: all payment methods.
	PaymentMethod   int                    `structs:"paymentMethod,omitempty"`

	// ConfirmationURL is the URL to which Flow will send the order details after the payment is made. This URL must
	// implement the confirmation logic. See GinOrderConfirmationCallback and HTTPOrderConfirmationCallback.
	ConfirmationURL string                 `structs:"urlConfirmation,omitempty"`

	// ReturnURL is the URL to which the user will be sent after the payment is made and confirmed. Usually on this
	// URL a boucher or confirmation message will be showed.
	ReturnURL       string                 `structs:"urlReturn,omitempty"`

	// TimeoutSeconds is the number of seconds the order should stay active for.
	TimeoutSeconds  uint64                 `structs:"timeout,omitempty"`

	// MerchantID will be set if a merchant is being a middleman in the payment.
	MerchantID      string                 `structs:"merchantId,omitempty"`

	// PaymentCurrency is optionally set to force the payer to pay in an specific currency.
	PaymentCurrency string                 `structs:"payment_currency,omitempty"`
}

// OrderResponse is the values sent by Flow if an order is successfully created.
type OrderResponse struct {
	// FlowID is the Flow identifier for the order.
	FlowID int    `json:"flowOrder"`

	// Token is the token used to buildPOST the payment URL.
	Token  string `json:"token"`

	// URL is the base URL to redirect for payment.
	URL    string `json:"url"`
}

// GetPaymentURL parses the payment URL of an order
func (or OrderResponse) GetPaymentURL() string {
	return fmt.Sprintf("%s?token=%s", or.URL, or.Token)
}


// isValid checks that the mandatory fields are set.
func (or OrderRequest) isValid() bool {
	if or.CommerceOrder == "" || or.Subject == "" || or.Amount == 0 || or.PayerEmail == "" {
		return false
	}

	return true
}

// GetOrder fetches the an Order based on the provided order token.
func (c Client) GetOrder(token string) (*Order, error) {
	url := c.buildGET("/payment/getStatus", map[string]interface{}{
		"token": token,
	})

	data, err := c.get(url)
	if err != nil {
		return nil, errors.Wrap(err, "unable to transact with the server")
	}

	var order Order
	err = jsoniter.Unmarshal(data, &order)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse response")
	}

	return &order, err
}

// GetOrderByCommerceID fetches the an Order based on the provided commerce identifier.
func (c Client) GetOrderByCommerceID(commerceID string) (*Order, error) {
	url := c.buildGET("/payment/getStatusByCommerceId", map[string]interface{}{
		"commerceId": commerceID,
	})

	data, err := c.get(url)
	if err != nil {
		return nil, errors.Wrap(err, "unable to transact with the server")
	}

	var order Order
	err = jsoniter.Unmarshal(data, &order)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse response")
	}

	return &order, err
}

// GetOrderByFlowID fetches the an Order based on the provided Flow identifier.
func (c Client) GetOrderByFlowID(flowOrderID int) (*Order, error) {
	url := c.buildGET("/payment/getStatusByFlowOrder", map[string]interface{}{
		"flowOrder": flowOrderID,
	})

	data, err := c.get(url)
	if err != nil {
		return nil, errors.Wrap(err, "unable to transact with the server")
	}

	var order Order
	err = jsoniter.Unmarshal(data, &order)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse response")
	}

	return &order, err
}

// CreateOrder creates a new order and returns its ID and token.
func (c Client) CreateOrder(or OrderRequest) ( *OrderResponse, error) {
	if !or.isValid() {
		return nil, errors.New("invalid order request: unfilled required values")
	}

	url, body := c.buildPOST("/payment/create", structs.Map(or))

	data, err := c.post(url, body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to transact with the server")
	}

	var result OrderResponse
	err = jsoniter.Unmarshal(data, &result)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse response")
	}

	return &result, err
}

// CreateEmailOrder creates an Order to be sent by email and returns its ID and token.
func (c Client) CreateEmailOrder(or OrderRequest) (orderID int, token string, err error) {
	if !or.isValid() {
		return -1, "", errors.New("invalid order request: unfilled required values")
	}

	url, body := c.buildPOST("/payment/createEmail", structs.Map(or))

	data, err := c.post(url, body)
	if err != nil {
		return -1, "", errors.Wrap(err, "unable to transact with the server")
	}

	result := struct{
		FlowID int `json:"flowOrder"`
		Token string `json:"token"`
	}{}
	err = jsoniter.Unmarshal(data, &result)
	if err != nil {
		return -1, "", errors.Wrap(err, "unable to parse response")
	}

	return result.FlowID, result.Token, err
}