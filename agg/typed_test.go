package agg

import (
	"testing"
)

type Address struct {
	City    string `bson:"city"`
	ZipCode string `bson:"zip"`
}

type User struct {
	ID      string   `bson:"_id"`
	Email   string   `bson:"email"`
	Name    string   // no tag, falls back to field name
	Ignored string   `bson:"-"`
	Address Address  `bson:"address"`
	Pointer *Address `bson:"ptrAddr"`
}

func TestKey(t *testing.T) {
	tests := []struct {
		field string
		want  string
		panic bool
	}{
		{field: "Email", want: "email"},
		{field: "ID", want: "_id"},
		{field: "Name", want: "Name"},
		{field: "Address", want: "address"},
		{field: "City", want: "address.city", panic: true}, // we map by root fields, not flattened by default if we ask for child field directly
		{field: "NonExistent", panic: true},
	}

	for _, tt := range tests {
		t.Run(tt.field, func(t *testing.T) {
			if tt.panic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("expected panic for field %s", tt.field)
					}
				}()
			}
			
			got := Key[User](tt.field)
			if got != tt.want {
				t.Errorf("Key[User](%q) = %q, want %q", tt.field, got, tt.want)
			}
		})
	}
	
	// Test cache hit by running again
	if got := Key[User]("Email"); got != "email" {
		t.Errorf("Cache hit failed, got %q", got)
	}
}
