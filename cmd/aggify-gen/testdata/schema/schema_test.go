package schema

import "testing"

func TestSchema(t *testing.T) {
	if User.ID != "_id" {
		t.Errorf("expected User.ID to be _id, got %q", User.ID)
	}
	if User.Addresses.City != "address.city" {
		t.Errorf("expected User.Addresses.City to be address.city, got %q", User.Addresses.City)
	}
	if User.Pointer.ZipCode != "ptrAddr.zip" {
		t.Errorf("expected User.Pointer.ZipCode to be ptrAddr.zip, got %q", User.Pointer.ZipCode)
	}
	if Order.Items.ProductID != "items.productId" {
		t.Errorf("expected Order.Items.ProductID to be items.productId, got %q", Order.Items.ProductID)
	}
}
