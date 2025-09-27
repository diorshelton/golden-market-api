package main

// import (
// 	"errors"
// 	"time"
// )

// type Vendor struct {
// 	ID        int
// 	Name      string
// 	OwnerID   *int
// 	IsNPC     bool
// 	CreatedAt time.Time
// }

// func SellProduct(user *User, product *Product, quantity int) error {
// 	if product.Stock < quantity {
// 		return errors.New("not enough stock for purchase")
// 	}

// 	total := quantity * int(product.Price)

// 	if total > int(user.Balance) {
// 		return errors.New("not enough coins for purchase")
// 	}

// 	AddToInventory(user, product.ID, quantity)
// 	user.Balance -= Coins(total)
// 	product.Stock = product.Stock - quantity

// 	return nil
// }
