package handlers

// import (
// 	"time"
// )

// type Transaction struct {
// 	ID         int               `json:"id"`
// 	BuyerID    int               `json:"buyer_id"`
// 	VendorID   int               `json:"vendor_id"`
// 	Items      []TransactionItem `json:"items"`
// 	TotalPrice Coins             `json:"total_price"`
// 	TimeStamp  time.Time         `json:"time_stamp"`
// }

// type TransactionItem struct {
// 	ProductID int   `json:"product_id"`
// 	Quantity  int   `json:"quantity"`
// 	Subtotal  Coins `json:"subtotal"`
// }

// type PurchaseItem struct {
// 	Product  *Product
// 	Quantity int
// }

// func CreateTransaction(buyer *User, purchases []PurchaseItem) *Transaction {
// 	var items []TransactionItem
// 	var total Coins

// 	for _, p := range purchases {
// 		subTotal := Coins(p.Quantity * int(p.Product.Price))
// 		items = append(items, TransactionItem{
// 			ProductID: p.Product.ID,
// 			Quantity:  p.Quantity,
// 			Subtotal:  subTotal,
// 		})
// 		total += subTotal
// 	}

// 	return &Transaction{
// 		BuyerID:    buyer.ID,
// 		VendorID:   purchases[0].Product.VendorID,
// 		Items:      items,
// 		TotalPrice: total,
// 		TimeStamp:  time.Now().UTC(),
// 	}
// }
