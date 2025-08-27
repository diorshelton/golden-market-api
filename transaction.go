package main

import (
	"time"
)

type Transaction struct {
	ID int `json:"id"`
	BuyerID int `json:"buyer_id"`
	VendorID int `json:"vendor_id"`
	Items []TransactionItem `json:"items"`
	TotalPrice Coins `json:"total_price"`
	TimeStamp time.Time `json:"time_stamp"`
}

type TransactionItem struct {
	ProductID int `json:"product_id"`
	Quantity int `json:"quantity"`
	Subtotal Coins `json:"subtotal"`
}

func RecordTransaction() Transaction {
	 var record Transaction
	return record
}