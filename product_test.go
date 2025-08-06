package main

import (
	"testing"
)

func TestAddToInventory(t *testing.T) {
	t.Run("Should add product to inventory", func(t *testing.T) {

		user := User{
			ID:        0,
			Inventory: []Item{},
		}

		//add one product with ID = 0 to inventory
		AddToInventory(&user, 0, 1)

		if len(user.Inventory) != 1 {
			t.Errorf("Expected inventory length 1, got %d", len(user.Inventory))
		}
	})
	t.Run("adds specific quantity of product to inventory", func(t *testing.T) {
		user := User{
			ID:        0,
			Inventory: []Item{},
		}

		AddToInventory(&user, 12, 5)

		if user.Inventory[0].ProductID != 12 || user.Inventory[0].Quantity != 5 {
			t.Errorf("Item not added correctly %+v", user.Inventory[0])
		}
	})

}
