package main

import (
	"testing"
)

func TestCreateTransaction(t *testing.T) {
	sonny := Vendor{ID: 1, Name: "Sonny", IsNPC: true}
	bacon := Product{ID: 93, Name: "bacon", Price: 9, Stock: 475, VendorID: 1}
	tomoyo := User{ID: 27, Username: "tomoyo", Balance: 1000}

	SellProduct(&tomoyo, &bacon, 3)

	got := CreateTransaction()

	if got.BuyerID != tomoyo.ID {
		t.Errorf("Expected buyerID of %v but got %v", tomoyo.ID, got.BuyerID)
	}

	if got.VendorID != sonny.ID {
		t.Errorf("Expected vendorID of %v but got %v", sonny.ID, got.VendorID)
	}

	if int(got.TotalPrice) != 27 {
		t.Errorf("Expected 27 but got %v", got.TotalPrice)
	}

	//	if got.Items[0].Quantity != 3 {
	//		t.Errorf("Expected 3 but got %v", got.Items[0].Quantity)
	//	}
}
