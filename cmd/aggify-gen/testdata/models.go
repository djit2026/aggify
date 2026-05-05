package models

type Address struct {
	City    string `bson:"city"`
	ZipCode string `bson:"zip"`
}

type User struct {
	ID        string  `bson:"_id"`
	Email     string  `bson:"email"`
	Name      string  `bson:"name,omitempty"`
	Ignored   string  `bson:"-"`
	Addresses Address `bson:"address"`
	Pointer   *Address `bson:"ptrAddr"`
}

type OrderItem struct {
	ProductID string `bson:"productId"`
	Quantity  int    `bson:"quantity"`
}

type Order struct {
	ID    string      `bson:"_id"`
	Items []OrderItem `bson:"items"`
}
