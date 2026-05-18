package expr

import (
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

func TestTrigOperators(t *testing.T) {
	tests := []struct {
		name     string
		expr     Expr
		expected any
	}{
		{
			name:     "Sin",
			expr:     Sin(Field("angle")),
			expected: bson.D{{Key: "$sin", Value: "$angle"}},
		},
		{
			name:     "Cos",
			expr:     Cos(Field("angle")),
			expected: bson.D{{Key: "$cos", Value: "$angle"}},
		},
		{
			name:     "Atan2",
			expr:     Atan2(Field("y"), Field("x")),
			expected: bson.D{{Key: "$atan2", Value: bson.A{"$y", "$x"}}},
		},
		{
			name:     "DegreesToRadians",
			expr:     DegreesToRadians(Field("angleDeg")),
			expected: bson.D{{Key: "$degreesToRadians", Value: "$angleDeg"}},
		},
		{
			name:     "RadiansToDegrees",
			expr:     RadiansToDegrees(Field("angleRad")),
			expected: bson.D{{Key: "$radiansToDegrees", Value: "$angleRad"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertBSON(t, tt.expr.Build(), tt.expected)
		})
	}
}

func TestDatePartsOperators(t *testing.T) {
	t.Run("Week", func(t *testing.T) {
		assertBSON(t,
			Week(Field("date")).Build(),
			bson.D{{Key: "$week", Value: "$date"}},
		)
	})
	t.Run("ISOWeek", func(t *testing.T) {
		assertBSON(t,
			ISOWeek(Field("date")).Build(),
			bson.D{{Key: "$isoWeek", Value: "$date"}},
		)
	})
	t.Run("ISODayOfWeek", func(t *testing.T) {
		assertBSON(t,
			ISODayOfWeek(Field("date")).Build(),
			bson.D{{Key: "$isoDayOfWeek", Value: "$date"}},
		)
	})
	t.Run("DateFromParts", func(t *testing.T) {
		got := DateFromParts(DateFromPartsOptions{
			Year:  Value(2024),
			Month: Value(1),
			Day:   Value(15),
		}).Build()
		assertBSON(t, got, bson.D{{Key: "$dateFromParts", Value: bson.D{
			{Key: "year", Value: 2024},
			{Key: "month", Value: 1},
			{Key: "day", Value: 15},
		}}})
	})
	t.Run("DateToParts", func(t *testing.T) {
		got := DateToParts(Field("createdAt")).Build()
		assertBSON(t, got, bson.D{{Key: "$dateToParts", Value: bson.D{
			{Key: "date", Value: "$createdAt"},
		}}})
	})
	t.Run("DateToParts ISO", func(t *testing.T) {
		got := DateToParts(Field("createdAt"), DateToPartsOptions{ISO: true}).Build()
		assertBSON(t, got, bson.D{{Key: "$dateToParts", Value: bson.D{
			{Key: "date", Value: "$createdAt"},
			{Key: "iso8601", Value: true},
		}}})
	})
}

func TestFieldOpsOperators(t *testing.T) {
	t.Run("GetFieldSimple", func(t *testing.T) {
		assertBSON(t,
			GetFieldSimple("price.usd").Build(),
			bson.D{{Key: "$getField", Value: "price.usd"}},
		)
	})
	t.Run("GetField", func(t *testing.T) {
		assertBSON(t,
			GetField("price.usd", Root).Build(),
			bson.D{{Key: "$getField", Value: bson.D{
				{Key: "field", Value: "price.usd"},
				{Key: "input", Value: "$$ROOT"},
			}}},
		)
	})
	t.Run("SetField", func(t *testing.T) {
		assertBSON(t,
			SetField("status", Root, Value("active")).Build(),
			bson.D{{Key: "$setField", Value: bson.D{
				{Key: "field", Value: "status"},
				{Key: "input", Value: "$$ROOT"},
				{Key: "value", Value: "active"},
			}}},
		)
	})
	t.Run("UnsetField", func(t *testing.T) {
		assertBSON(t,
			UnsetField("password", Root).Build(),
			bson.D{{Key: "$unsetField", Value: bson.D{
				{Key: "field", Value: "password"},
				{Key: "input", Value: "$$ROOT"},
			}}},
		)
	})
}

func TestMiscOperators(t *testing.T) {
	t.Run("Literal", func(t *testing.T) {
		assertBSON(t,
			Literal("$notAField").Build(),
			bson.D{{Key: "$literal", Value: "$notAField"}},
		)
	})
	t.Run("Rand", func(t *testing.T) {
		got := Rand().Build()
		if _, ok := got.(bson.D); !ok {
			t.Errorf("Rand() should build to bson.D")
		}
	})
	t.Run("TsSecond", func(t *testing.T) {
		assertBSON(t,
			TsSecond(Field("ts")).Build(),
			bson.D{{Key: "$tsSecond", Value: "$ts"}},
		)
	})
	t.Run("TsIncrement", func(t *testing.T) {
		assertBSON(t,
			TsIncrement(Field("ts")).Build(),
			bson.D{{Key: "$tsIncrement", Value: "$ts"}},
		)
	})
}

func TestArrayV2Operators(t *testing.T) {
	t.Run("Range", func(t *testing.T) {
		assertBSON(t,
			Range(Value(0), Value(10), Value(2)).Build(),
			bson.D{{Key: "$range", Value: bson.A{0, 10, 2}}},
		)
	})
	t.Run("Range no step", func(t *testing.T) {
		assertBSON(t,
			Range(Value(0), Value(5)).Build(),
			bson.D{{Key: "$range", Value: bson.A{0, 5}}},
		)
	})
	t.Run("SortArray scalar", func(t *testing.T) {
		assertBSON(t,
			SortArray(Field("scores"), nil, Value(-1)).Build(),
			bson.D{{Key: "$sortArray", Value: bson.D{
				{Key: "input", Value: "$scores"},
				{Key: "sortBy", Value: -1},
			}}},
		)
	})
	t.Run("SortArray document", func(t *testing.T) {
		assertBSON(t,
			SortArray(Field("items"), []SortArrayField{{"price", SortArrayDesc}}, nil).Build(),
			bson.D{{Key: "$sortArray", Value: bson.D{
				{Key: "input", Value: "$items"},
				{Key: "sortBy", Value: bson.D{{Key: "price", Value: -1}}},
			}}},
		)
	})
	t.Run("FirstN", func(t *testing.T) {
		assertBSON(t,
			FirstN(Value(3), Field("scores")).Build(),
			bson.D{{Key: "$firstN", Value: bson.D{
				{Key: "n", Value: 3},
				{Key: "input", Value: "$scores"},
			}}},
		)
	})
	t.Run("MaxN", func(t *testing.T) {
		assertBSON(t,
			MaxN(Value(3), Field("scores")).Build(),
			bson.D{{Key: "$maxN", Value: bson.D{
				{Key: "n", Value: 3},
				{Key: "input", Value: "$scores"},
			}}},
		)
	})
}
