package main

import "errors"

import "time"

type Product struct {
	ID int
	Name string
	Stock int
	RestockRate int
	MaxStock int
	LastRestock time.Time
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

func RemoveFromInventory(user *User, productID, quantity int) error {
	for i, item := range user.Inventory {
		if item.ProductID == productID {
			if item.Quantity < quantity {
				return errors.New("not enough items to remove")
			}
			user.Inventory[i].Quantity -= quantity
			if user.Inventory[i].Quantity == 0 {
				user.Inventory = append(user.Inventory[:i], user.Inventory[i+1:]...)
			}
			return nil
		}
	}
	return errors.New("item not found inventory")
}

func RestockProduct(p *Product)  {
	now := time.Now()
	if now.Sub(p.LastRestock) >= time.Hour {
		newStock := p.Stock + p.RestockRate
		if newStock > p.MaxStock {
			newStock = p.MaxStock
		}
		p.Stock= newStock
		p.LastRestock = now
	}
}