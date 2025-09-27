package main

import (
	// "reflect"
	"testing"
	// "time"
)

func TestCreateTransaction(t *testing.T) {
	t.Skip()
	// t.Skip("Create transaction with fields: buyer, vendor, product quantity and total price", func(t *testing.T) {
	// 	// Test transaction creation with correct buyer, vendor, product, quantity and total price

	// 	sonny := Vendor{ID: 1}
	// 	tomoyo := User{ID: 27, Username: "tomoyo", Balance: 1000}

	// 	bacon := Product{ID: 93, Name: "bacon", Price: 9, Stock: 475, VendorID: sonny.ID}
	// 	sweetPotato := Product{ID: 1, Name: "sweet potato", Price: 4, Stock: 54, VendorID: sonny.ID}

	// 	purchases := []PurchaseItem{
	// 		{Product: &sweetPotato, Quantity: 6},
	// 		{Product: &bacon, Quantity: 6},
	// 	}

	// 	got := CreateTransaction(&tomoyo, purchases)

	// 	want := &Transaction{
	// 		BuyerID:  27,
	// 		VendorID: 1,
	// 		Items: []TransactionItem{
	// 			{ProductID: sweetPotato.ID, Quantity: 6, Subtotal: 24},
	// 			{ProductID: bacon.ID, Quantity: 6, Subtotal: 54},
	// 		},
	// 		TotalPrice: 78,
	// 	}
	// 	// Compare fields but ignore ID/TimeStamp
	// 	if got.BuyerID != want.BuyerID ||
	// 		got.VendorID != want.VendorID ||
	// 		got.TotalPrice != want.TotalPrice ||
	// 		!reflect.DeepEqual(got.Items, want.Items) {
	// 		t.Errorf("\ngot:%+v\nwant:%+v", got, want)
	// 	}
	// })
	// t.Run("Check if TimeStamp exists", func(t *testing.T) {
	// 	t.Skip("skipping test")

	// 	sonny := Vendor{ID: 1}
	// 	dandara := User{ID: 1654, Username: "dandara", Balance: 4000}

	// 	tambourine := Product{ID: 93, Name: "tambourine", Price: 54, Stock: 475, VendorID: sonny.ID}

	// 	drum := Product{ID: 1, Name: "drum", Price: 94, Stock: 353, VendorID: sonny.ID}

	// 	purchases := []PurchaseItem{
	// 		{Product: &tambourine, Quantity: 3},
	// 		{Product: &drum, Quantity: 1},
	// 	}

	// 	trx := CreateTransaction(&dandara, purchases)

	// 	if trx.TimeStamp.IsZero() {
	// 		t.Errorf("expected timestamp to be set but got zero %v", trx.TimeStamp)
	// 	}

	// 	now := time.Now().UTC()
	// 	if trx.TimeStamp.Before(now.Add(-2*time.Second)) || trx.TimeStamp.After(now.Add(2*time.Second)) {
	// 		t.Errorf("timestamp not within 2 second range, got %v, expected around %v", trx.TimeStamp, now)
	// 	}

	// 	if trx.TimeStamp.Location() != time.UTC {
	// 		t.Errorf("timestamp not stored in UTC, got %v", trx.TimeStamp.Location())
	// 	}
	// })
}
