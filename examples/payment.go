package main

import (
	"fmt"
	"github.com/CamiloHernandez/go-flow"
	"net/http"
)

func main() {
	c := flow.NewClient("your api key", "your secret key")
	c.SetSandbox()

	result, err := c.CreateOrder(flow.OrderRequest{
		CommerceOrder:   "123123",
		Subject:         "Test Order",
		Amount:          1000,
		PayerEmail:      "example@example.com",
		ConfirmationURL: "http://example.com/confirmation",
		ReturnURL:       "http://example.com/return",
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("ID:", result.FlowID)
	fmt.Println("Token:", result.Token)
	fmt.Println("URL:", result.GetPaymentURL())

	http.HandleFunc("/confirm", c.HTTPOrderConfirmationCallback(func(o *flow.Order) {
		fmt.Println("Order", o.CommerceOrder, "confirmed with status", o.Status)
	}))

	http.HandleFunc("/return", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Order completed!"))
	})

	panic(http.ListenAndServe(":80", nil))
}