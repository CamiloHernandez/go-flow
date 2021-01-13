package flow

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// GinOrderConfirmationCallback returns a gin.HandlerFunc that processes the Flow callback request.
// It validates the provided token, and if the token matches a accepted order, the onAccepted function gets called.
func (c *Client) GinOrderConfirmationCallback(onAccepted func(*Order)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokens, set := ctx.Request.PostForm["token"]
		if !set {
			ctx.Status(http.StatusBadRequest)
			return
		}

		if len(tokens) != 1 {
			ctx.Status(http.StatusBadRequest)
			return
		}

		token := tokens[0]

		order, err := c.GetOrder(token)
		if err != nil {
			ctx.Status(http.StatusInternalServerError)
			return
		}

		if order.Status == OrderStatusPayed {
			ctx.Status(http.StatusAccepted)
			onAccepted(order)
			return
		}

		ctx.Status(http.StatusUnauthorized)
		return
	}
}

// HTTPOrderConfirmationCallback returns a http.HandlerFunc that processes the Flow callback request.
// It validates the provided token, and if the token matches a accepted order, the onAccepted function gets called.
func (c *Client) HTTPOrderConfirmationCallback(onAccepted func(*Order)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		token := r.Form.Get("token")
		if token == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		order, err := c.GetOrder(token)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if order.Status == OrderStatusPayed {
			w.WriteHeader(http.StatusOK)
			onAccepted(order)
			return
		}

		w.WriteHeader(http.StatusUnauthorized)
		return
	}
}

// GinRefundConfirmationCallback returns a gin.HandlerFunc that processes the Flow callback request.
// It validates the provided token, and if the token matches an accepted refund, the onAccepted function gets called.
func (c *Client) GinRefundConfirmationCallback(onAccepted func(*RefundStatus)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := ctx.Request.ParseForm()
		if err != nil {
			ctx.Status(http.StatusBadRequest)
			return
		}

		tokens, set := ctx.Request.Form["token"]
		if !set {
			ctx.Status(http.StatusBadRequest)
			return
		}

		if len(tokens) != 1 {
			ctx.Status(http.StatusBadRequest)
			return
		}

		token := tokens[0]

		refund, err := c.GetRefundStatus(token)
		if err != nil {
			ctx.Status(http.StatusInternalServerError)
			return
		}

		if refund.Status == RefundStatusAccepted || refund.Status == RefundStatusRefunded {
			ctx.Status(http.StatusOK)
			onAccepted(refund)
			return
		}

		ctx.Status(http.StatusUnauthorized)
		return
	}
}

// HTTPRefundConfirmationCallback returns a http.HandlerFunc that processes the Flow callback request.
// It validates the provided token, and if the token matches an accepted refund, the onAccepted function gets called.
func (c *Client) HTTPRefundConfirmationCallback(onAccepted func(*RefundStatus)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		token := r.Form.Get("token")
		if token == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		refund, err := c.GetRefundStatus(token)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if refund.Status == RefundStatusAccepted || refund.Status == RefundStatusRefunded {
			w.WriteHeader(http.StatusAccepted)
			onAccepted(refund)
			return
		}

		w.WriteHeader(http.StatusUnauthorized)
		return
	}
}
