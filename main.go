package main

type Item struct {
	ProductID int
	Quantity  int
}

type User struct {
	ID        int
	Inventory []Item
}

func AddToInventory(u *User, productID, quantity int) {
	for i, item := range u.Inventory {
		if item.ProductID == productID {
			u.Inventory[i].Quantity += quantity
			return
		}
	}
	u.Inventory = append(u.Inventory, Item{
		ProductID: productID,
		Quantity:  quantity,
	})
}

func main() {
}
