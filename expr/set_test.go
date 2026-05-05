package expr

import (
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

func TestSetOperators(t *testing.T) {
	tests := []struct {
		name     string
		expr     Expr
		expected any
	}{
		{
			name:     "SetEquals",
			expr:     SetEquals(Field("roles"), Value([]string{"admin", "editor"})),
			expected: bson.D{{Key: "$setEquals", Value: bson.A{"$roles", []string{"admin", "editor"}}}},
		},
		{
			name:     "SetIntersection",
			expr:     SetIntersection(Field("roles"), Field("permissions")),
			expected: bson.D{{Key: "$setIntersection", Value: bson.A{"$roles", "$permissions"}}},
		},
		{
			name:     "SetDifference",
			expr:     SetDifference(Field("setA"), Field("setB")),
			expected: bson.D{{Key: "$setDifference", Value: bson.A{"$setA", "$setB"}}},
		},
		{
			name:     "SetIsSubset",
			expr:     SetIsSubset(Value([]int{1, 2}), Field("numbers")),
			expected: bson.D{{Key: "$setIsSubset", Value: bson.A{[]int{1, 2}, "$numbers"}}},
		},
		{
			name:     "AnyElementTrue",
			expr:     AnyElementTrue(Field("booleans")),
			expected: bson.D{{Key: "$anyElementTrue", Value: "$booleans"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.expr.Build()
			assertBSON(t, got, tt.expected)
		})
	}
}
