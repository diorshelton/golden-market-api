package main

import "errors"

import "time"

type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       Coins     `json:"coins"`
	Stock       int       `json:"stock"`
	RestockRate int       `json:"restock_rate"`
	MaxStock    int       `json:"max_stock"`
	VendorID    int       `json:"vendor_id"`
	LastRestock time.Time `json:"last_restock"`
	CreatedAt   time.Time `json:"created_at"`
}

type Item struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type Coins int32

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

func RestockProduct(product *Product) {
	now := time.Now()
	if now.Sub(product.LastRestock) >= time.Hour {
		newStock := product.Stock + product.RestockRate
		if newStock > product.MaxStock {
			newStock = product.MaxStock
		}
		product.Stock = newStock
		product.LastRestock = now
	}
}
