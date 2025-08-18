package main

import (
	"testing"
	"time"
)

func TestSellProduct(t *testing.T) {
	var sonny = Vendor{
		ID: 01,
	}

	var product = Product{
		ID:          1,
		Name:        "Sweet Potato",
		Price: 4,
		Stock:       5,
		RestockRate: 3,
		MaxStock:    10,
		LastRestock: time.Now().Add(-2 * time.Hour),
		VendorID:    sonny.ID,
	}

	var user = User{
		ID:      11,
		Balance: 34,
	}

	itemQuantity := 1
	expectedBal := user.Balance - product.Price
	expectedStock := product.Stock - itemQuantity

	SellProduct(&user, &product, itemQuantity)

	if user.Balance != expectedBal {
		t.Errorf("User's balance should be %v got %v", expectedBal, user.Balance)
	}

	if product.Stock != expectedStock {
		t.Errorf("Product stock should be %v got %v", expectedStock, product.Stock)
	}

	expectedInventory := len(user.Inventory)
	want := itemQuantity

	if expectedInventory != 1 || user.Inventory[0].ProductID != product.ID {
		t.Errorf("got %v but want %v, user.Inventory:%v", expectedInventory, want, user.Inventory)
	}

}
