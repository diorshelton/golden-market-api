package main

import (
	"testing"
)

func TestSellProduct(t *testing.T) {
	sonny := Vendor{
		Name:"Sonny",
	}

	got := sonny.Name
	want := "Sonny"

	if got != want {
		t.Errorf("got %s but want %s", got,want)
	}
}