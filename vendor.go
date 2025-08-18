package main

import "time"

type Vendor struct {
	ID        int
	Name      string
	OwnerID   *int
	IsNPC     bool
	CreatedAt time.Time
}

func SellProduct(user *User, product *Product, quantity int) {
	AddToInventory(user, product.ID, quantity)
	user.Balance = user.Balance - product.Price
	product.Stock = product.Stock - quantity
}
