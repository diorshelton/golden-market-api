package main

import(
	"testing"
)

func TestRecordTransaction(t *testing.T) {
	SellProduct()
	got := RecordTransaction()

	 {
		t.Errorf("Expected record %v but got %v", record, got)
	}
}