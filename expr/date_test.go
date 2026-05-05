package expr

import (
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

func TestDateOperators(t *testing.T) {
	tests := []struct {
		name     string
		expr     Expr
		expected any
	}{
		{
			name: "Year",
			expr: Year(Field("createdAt")),
			expected: bson.D{{Key: "$year", Value: "$createdAt"}},
		},
		{
			name: "DateAdd",
			expr: DateAdd(Field("createdAt"), Value(5), "day"),
			expected: bson.D{{Key: "$dateAdd", Value: bson.D{
				{Key: "startDate", Value: "$createdAt"},
				{Key: "amount", Value: 5},
				{Key: "unit", Value: "day"},
			}}},
		},
		{
			name: "DateTrunc with binSize",
			expr: DateTrunc(Field("date"), "month", 3),
			expected: bson.D{{Key: "$dateTrunc", Value: bson.D{
				{Key: "date", Value: "$date"},
				{Key: "unit", Value: "month"},
				{Key: "binSize", Value: 3},
			}}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.expr.Build()
			assertBSON(t, got, tt.expected)
		})
	}
}
