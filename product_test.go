package main

import (
	"testing"
)

func TestAddToInventory(t *testing.T) {

	user := User{
		ID:0,
		Inventory: []Item{},
	}


	AddToInventory(&user, 0, 1)

	if len(user.Inventory) != 1 {
		t.Errorf("Expected inventory length 1, got %d", len(user.Inventory))
	}

}