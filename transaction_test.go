package main

import (
	"testing"
)

func TestCreateTransaction(t *testing.T) {
	//Test transaction creation with correct buyer, vendor, product, quantity, total price
	sonny := Vendor{ID: 1, Name: "Sonny", IsNPC: true}
	tomoyo := User{ID: 27, Username: "tomoyo", Balance: 1000}

	bacon := Product{ID: 93, Name: "bacon", Price: 9, Stock: 475, VendorID: 1}
	sweetPotato := Product{ID: 1, Name: "sweet potato", Price: 4, Stock: 54, VendorID: 1}

	item1 := PurchaseItem{&sweetPotato, 6}
	item2 := PurchaseItem{&bacon, 6}

	purchases := []PurchaseItem{item1, item2}

	got := CreateTransaction(&tomoyo, purchases)

	if got.BuyerID != tomoyo.ID {
		t.Errorf("Expected buyerID of %v but got %v", tomoyo.ID, got.BuyerID)
	}

	if got.VendorID != sonny.ID {
		t.Errorf("Expected vendorID of %v but got %v", sonny.ID, got.VendorID)
	}

	if int(got.TotalPrice) != 78 {
		t.Errorf("Expected 27 but got %v", got.TotalPrice)
	}

	if len(got.Items) != 2 {
		t.Errorf("Expected 2 but got %v", len(got.Items))
	}

	if got.Items[0].Subtotal != 24 {
		t.Errorf("Expected 24 but got %v", got.Items[0].Subtotal)
	}

	if got.Items[1].Subtotal != 54 {
		t.Errorf("Expected 54 but got %v", got.Items[1].Subtotal)
	}
}
