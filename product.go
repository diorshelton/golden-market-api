package main

type Product struct {
	ID int
}

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
		//remove item if quantity <= 0
		if item.ProductID == productID && itemCount-quantity <= 0 {
			user.Inventory = append(user.Inventory[:i], user.Inventory[i+1:]...)
			return
		}
		user.Inventory[i].Quantity -= quantity
	}
}
