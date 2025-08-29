package main

import (
	"reflect"
	"testing"
	// "time"
)

// Test transaction creation with correct buyer, vendor, product, quantity and total price
func TestCreateTransaction(t *testing.T) {
	t.Run("Create transaction with fields: buyer, vender, product quantity and total price", func(t *testing.T) {

		sonny := Vendor{ID: 1}
		tomoyo := User{ID: 27, Username: "tomoyo", Balance: 1000}

		bacon := Product{ID: 93, Name: "bacon", Price: 9, Stock: 475, VendorID: sonny.ID}
		sweetPotato := Product{ID: 1, Name: "sweet potato", Price: 4, Stock: 54, VendorID: sonny.ID}

		purchases := []PurchaseItem{
			{Product: &sweetPotato, Quantity: 6},
			{Product: &bacon, Quantity: 6},
		}

		got := CreateTransaction(&tomoyo, purchases)

		want := &Transaction{
			BuyerID:  27,
			VendorID: 1,
			Items: []TransactionItem{
				{ProductID: sweetPotato.ID, Quantity: 6, Subtotal: 24},
				{ProductID: bacon.ID, Quantity: 6, Subtotal: 54},
			},
			TotalPrice: 78,
		}

		// Compare only meaningful fields (ignore ID/TimeStamp)
		if got.BuyerID != want.BuyerID ||
			got.VendorID != want.VendorID ||
			got.TotalPrice != want.TotalPrice ||
			!reflect.DeepEqual(got.Items, want.Items) {
			t.Errorf("\ngot:%+v\nwant:%+v", got, want)
		}
	})
}
