package main

type Item struct {
	ID int
}

type User struct {
	ID int
	Inventory  []Item
}

func AddToInventory(u *User, productID, quantity int) []Item {
	newItem := Item{ID: 0}

	u.Inventory = append(u.Inventory, newItem)
	return u.Inventory
}

func main() {
}
