package main

type Item struct {
	ProductID int
	Quantity  int
}

type User struct {
	ID        int
	Inventory []Item
}

func AddToInventory(u *User, productID, quantity int) []Item {
	newItem := Item{ProductID: productID, Quantity: quantity}

	u.Inventory = append(u.Inventory, newItem)

	return u.Inventory
}

func main() {
}
