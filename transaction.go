package main

import (
	"time"
)

type Transaction struct {
	ID         int               `json:"id"`
	BuyerID    int               `json:"buyer_id"`
	VendorID   int               `json:"vendor_id"`
	Items      []TransactionItem `json:"items"`
	TotalPrice Coins             `json:"total_price"`
	TimeStamp  time.Time         `json:"time_stamp"`
}

type TransactionItem struct {
	ProductID int   `json:"product_id"`
	Quantity  int   `json:"quantity"`
	Subtotal  Coins `json:"subtotal"`
}

type PurchaseItem struct {
	Product  *Product
	Quantity int
}

func CreateTransaction(buyer *User, purchases []PurchaseItem) *Transaction {
	tx := Transaction{}
	tx.BuyerID = buyer.ID
	tx.VendorID = purchases[0].Product.VendorID
	tx.TotalPrice = 0
	// var txTotal Coins = 0

	for _, item := range purchases {
		//Calculate subtotal and total for items
		subTotal := item.Quantity * int(item.Product.Price)
		tx.TotalPrice += Coins(subTotal)
		//Build transaction item to add to Transaction struct
		txItem := TransactionItem{
			ProductID: item.Product.ID,
			Quantity:  item.Quantity,
			Subtotal:  Coins(subTotal),
		}

		tx.Items = append(tx.Items, txItem)
	}

	// tx.TimeStamp = time.Now()

	return &tx
}
