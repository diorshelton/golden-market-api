package main

type Item struct {
	ProductID int
	Quantity  int
}

type User struct {
	ID        int
	Inventory []Item
}

func AddToInventory(user *User, productID, quantity int) {
	for i, item := range user.Inventory {
		if item.ProductID == productID {
			user.Inventory[i].Quantity += quantity
			return
		}
	}
	user.Inventory = append(user.Inventory, Item{
		ProductID: productID,
		Quantity:  quantity,
	})
}

func RemoveFromInventory(user *User, productID, quantity int) {
	for i, item := range user.Inventory {
		itemCount := user.Inventory[i].Quantity
		if item.ProductID == productID && itemCount-quantity <= 0 {
			user.Inventory[i].Quantity = 0
			return
		}
		user.Inventory[i].Quantity -= quantity
	}
}

func main() {
}
