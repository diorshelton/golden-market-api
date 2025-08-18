package main

import "time"

type Vendor struct {
	ID int
	Name string
	OwnerID *int
	IsNPC bool
	CreatedAt time.Time
}