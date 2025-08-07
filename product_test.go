package main

import (
	"fmt"
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

		AddToInventory(&user, 4, 1)

		RemoveFromInventory(&user, 4)


		got := user.Inventory[0].Quantity

		want := 0
		fmt.Printf("Got %v and an Inventory of:%v", got,user.Inventory)
		if got != want {
			t.Errorf("Wanted %v but got %v", want, got)
		}
	})
}
