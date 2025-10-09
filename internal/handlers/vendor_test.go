package handlers

import (
	"testing"
	// "time"
)

func TestSellProduct(t *testing.T) {
	t.Skip("skipping test")

	// var sonny = Vendor{
	// 	ID: 01,
	// }

	// var product = Product{
	// 	ID:       1,
	// 	Name:     "Sweet Potato",
	// 	Price:    4,
	// 	Stock:    5,
	// 	VendorID: sonny.ID,
	// }

	// var user = User{
	// 	ID:      11,
	// 	Balance: 34,
	// }

	// itemCount := 1
	// expectedBal := user.Balance - product.Price
	// expectedStock := product.Stock - itemCount

	// SellProduct(&user, &product, itemCount)

	// if user.Balance != expectedBal {
	// 	t.Errorf("User's balance should be %v got %v", expectedBal, user.Balance)
	// }

	// if product.Stock != expectedStock {
	// 	t.Errorf("Product stock should be %v got %v", expectedStock, product.Stock)
	// }

	// got := user.Inventory[0].Quantity
	// want := itemCount

	// if got != want {
	// 	t.Errorf("got user inventory of %v but want %v, Inventory:%+v", got, want, user.Inventory)
	// }

	// t.Run("user can purchase multiple items", func(t *testing.T) {
	// 	var sonny = Vendor{
	// 		ID: 01,
	// 	}

	// 	var product = Product{
	// 		ID:       1,
	// 		Name:     "Peach",
	// 		Price:    4,
	// 		Stock:    5,
	// 		VendorID: sonny.ID,
	// 	}

	// 	var user = User{
	// 		ID:      11,
	// 		Balance: 96,
	// 	}

	// 	itemCount := 4
	// 	saleTotal := product.Price * Coins(itemCount)
	// 	expectedBal = user.Balance - saleTotal

	// 	SellProduct(&user, &product, itemCount)

	// 	if user.Balance != expectedBal {
	// 		t.Errorf("Incorrect balance got %v but expected %v", user.Balance, expectedBal)
	// 	}

	// 	if user.Inventory[0].Quantity != 4 {
	// 		t.Errorf("Expected an inventory of 4 but got %+v", user.Inventory)
	// 	}
	// })

	// t.Run("Should have enough stock to make sale", func(t *testing.T) {
	// 	var sonny = Vendor{
	// 		ID: 01,
	// 	}

	// 	var product = Product{
	// 		ID:          1,
	// 		Name:        "Tomato",
	// 		Price:       3,
	// 		Stock:       5,
	// 		MaxStock:    10,
	// 		LastRestock: time.Now(),
	// 		VendorID:    sonny.ID,
	// 	}

	// 	var user = User{
	// 		ID:      11,
	// 		Balance: 96,
	// 	}

	// 	itemCount := 7
	// 	got := SellProduct(&user, &product, itemCount)

	// 	if product.Stock < 0 {
	// 		t.Errorf("Product stock should not be negative. product stock:%+v", product.Stock)
	// 	}

	// 	if got == nil {
	// 		t.Error("Expected error but did not get one")
	// 	}
	// })

	// t.Run("User must have enough coins for sale", func(t *testing.T) {
	// 	var sonny = Vendor{
	// 		ID: 01,
	// 	}

	// 	var product = Product{
	// 		ID:       1,
	// 		Name:     "Lemon Cake Slice",
	// 		Price:    2,
	// 		Stock:    5,
	// 		VendorID: sonny.ID,
	// 	}

	// 	var user = User{
	// 		ID:      11,
	// 		Balance: 3,
	// 	}

	// 	itemCount := 3
	// 	got := SellProduct(&user, &product, itemCount)

	// 	saleTotal := product.Price * Coins(itemCount)
	// 	expectedBal = user.Balance - saleTotal

	// 	if got == nil {
	// 		t.Errorf("Wanted an error but did not get one User")
	// 	}
	// 	if user.Balance < 0 {
	// 		t.Errorf("user should not have a negative balance User Bal:%v", user.Balance)
	// 	}
	// })
}
