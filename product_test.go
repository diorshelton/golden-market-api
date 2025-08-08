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
	t.Run("adds specified product quantity to user inventory", func(t *testing.T) {
		user := User{
			ID:        0,
			Inventory: []Item{},
		}

		AddToInventory(&user, 12, 7)
		AddToInventory(&user, 5, 1)
		AddToInventory(&user, 12, 1)

		expected := Item{ProductID: 12, Quantity: 8}

		if user.Inventory[0].ProductID != 12 || user.Inventory[0].Quantity != 8 {
			t.Errorf("Item not added correctly, got %+v expected %+v", user.Inventory[0], expected)
		}
	})

}

func TestRemoveFromInventory(t *testing.T) {
	t.Run("should remove an item from inventory", func(t *testing.T) {
		user := User{
			ID:        0,
			Inventory: []Item{},
		}

		AddToInventory(&user, 41, 1)
		AddToInventory(&user, 4, 3)
		AddToInventory(&user, 11, 7)
		AddToInventory(&user, 4, 3)
		RemoveFromInventory(&user, 4, 1)

		got := checkItemQuantity(&user, 4)
		want := 5

		assertItemQuantity(t, got, want, &user)
	})
	t.Run("Inventory should NOT be negative", func(t *testing.T) {
		user := User{
			ID:        0,
			Inventory: []Item{},
		}

		AddToInventory(&user, 6, 1)
		AddToInventory(&user, 9, 33)
		RemoveFromInventory(&user, 6, 4)

		got := checkItemQuantity(&user, 6)
		want := 0
		assertItemQuantity(t, got, want, &user)
	})
	t.Run("Remove an item that does not exist", func(t *testing.T) {
		t.Skip("NOT YET IMPLEMENTED")

	})
}

func assertItemQuantity(t testing.TB, got, want int, user *User) {
	t.Helper()

	if got != want {
		t.Errorf("Got %v but wanted %v", got, want)
		t.Errorf("UserInventory:%+v", user.Inventory)
	}
}

// check User inventory for item quantity when given productID
func checkItemQuantity(user *User, itemID int) int {
	for _, item := range user.Inventory {
		if item.ProductID == itemID {
			return item.Quantity
		}
	}
	return 0
}
