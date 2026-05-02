package expr

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

// assertBSON marshals both got and want to Extended JSON and compares them.
// Use this in tests to get readable diffs.
func assertBSON(t *testing.T, got any, want any) {
	t.Helper()
	g, err := bson.MarshalExtJSON(got, true, true)
	require.NoError(t, err, "marshal got")
	w, err := bson.MarshalExtJSON(want, true, true)
	require.NoError(t, err, "marshal want")
	assert.JSONEq(t, string(w), string(g))
}

// mustJSON is a debug helper used in test output.
func mustJSON(v any) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}

func TestField(t *testing.T) {
	assert.Equal(t, "$price", Field("price").Build())
	assert.Equal(t, "$items.status", Field("items.status").Build())
}

func TestVar(t *testing.T) {
	assert.Equal(t, "$$item", Var("item").Build())
	assert.Equal(t, "$$ROOT", Root.Build())
	assert.Equal(t, "$$REMOVE", Remove.Build())
}

func TestValue(t *testing.T) {
	assert.Equal(t, 42, Value(42).Build())
	assert.Equal(t, "hello", Value("hello").Build())
	assert.Nil(t, Value(nil).Build())
}

func TestEq(t *testing.T) {
	got := Eq(Field("status"), Value("active")).Build()
	want := bson.D{{Key: "$eq", Value: bson.A{"$status", "active"}}}
	assertBSON(t, got, want)
}

func TestNe(t *testing.T) {
	got := Ne(Field("count"), Value(0)).Build()
	want := bson.D{{Key: "$ne", Value: bson.A{"$count", 0}}}
	assertBSON(t, got, want)
}

func TestGtGteLtLte(t *testing.T) {
	assertBSON(t,
		Gt(Field("age"), Value(18)).Build(),
		bson.D{{Key: "$gt", Value: bson.A{"$age", 18}}},
	)
	assertBSON(t,
		Gte(Field("age"), Value(21)).Build(),
		bson.D{{Key: "$gte", Value: bson.A{"$age", 21}}},
	)
	assertBSON(t,
		Lt(Field("stock"), Value(10)).Build(),
		bson.D{{Key: "$lt", Value: bson.A{"$stock", 10}}},
	)
	assertBSON(t,
		Lte(Field("score"), Value(100)).Build(),
		bson.D{{Key: "$lte", Value: bson.A{"$score", 100}}},
	)
}

func TestAnd(t *testing.T) {
	got := And(
		Eq(Field("a"), Value(1)),
		Eq(Field("b"), Value(2)),
	).Build()
	want := bson.D{{Key: "$and", Value: bson.A{
		bson.D{{Key: "$eq", Value: bson.A{"$a", 1}}},
		bson.D{{Key: "$eq", Value: bson.A{"$b", 2}}},
	}}}
	assertBSON(t, got, want)
}

func TestOr(t *testing.T) {
	got := Or(
		Eq(Field("role"), Value("admin")),
		Eq(Field("role"), Value("mod")),
	).Build()
	want := bson.D{{Key: "$or", Value: bson.A{
		bson.D{{Key: "$eq", Value: bson.A{"$role", "admin"}}},
		bson.D{{Key: "$eq", Value: bson.A{"$role", "mod"}}},
	}}}
	assertBSON(t, got, want)
}

func TestCond(t *testing.T) {
	got := Cond(
		Eq(Field("status"), Value("active")),
		Field("price"),
		Value(0),
	).Build()
	want := bson.D{{Key: "$cond", Value: bson.D{
		{Key: "if", Value: bson.D{{Key: "$eq", Value: bson.A{"$status", "active"}}}},
		{Key: "then", Value: "$price"},
		{Key: "else", Value: 0},
	}}}
	assertBSON(t, got, want)
}

func TestIfNull(t *testing.T) {
	got := IfNull(Field("discount"), Value(0)).Build()
	want := bson.D{{Key: "$ifNull", Value: bson.A{"$discount", 0}}}
	assertBSON(t, got, want)
}

func TestFilter(t *testing.T) {
	got := Filter(
		Field("items"), "item",
		Gt(Var("item.qty"), Value(5)),
	).Build()
	want := bson.D{{Key: "$filter", Value: bson.D{
		{Key: "input", Value: "$items"},
		{Key: "as", Value: "item"},
		{Key: "cond", Value: bson.D{{Key: "$gt", Value: bson.A{"$$item.qty", 5}}}},
	}}}
	assertBSON(t, got, want)
}

func TestMergeObjects(t *testing.T) {
	got := MergeObjects(Root, Raw(bson.D{{Key: "extra", Value: 1}})).Build()
	want := bson.D{{Key: "$mergeObjects", Value: bson.A{"$$ROOT", bson.D{{Key: "extra", Value: 1}}}}}
	assertBSON(t, got, want)
}

func TestAdd(t *testing.T) {
	got := Add(Field("a"), Field("b"), Value(10)).Build()
	want := bson.D{{Key: "$add", Value: bson.A{"$a", "$b", 10}}}
	assertBSON(t, got, want)
}

func TestConcat(t *testing.T) {
	got := Concat(Field("first"), Value(" "), Field("last")).Build()
	want := bson.D{{Key: "$concat", Value: bson.A{"$first", " ", "$last"}}}
	assertBSON(t, got, want)
}

func TestLet(t *testing.T) {
	got := Let(
		[]LetBinding{{"total", Field("price")}},
		Multiply(Var("total"), Value(1.1)),
	).Build()
	_ = mustJSON(got) // ensure it is serialisable
	assertBSON(t, got, bson.D{{Key: "$let", Value: bson.D{
		{Key: "vars", Value: bson.D{{Key: "total", Value: "$price"}}},
		{Key: "in", Value: bson.D{{Key: "$multiply", Value: bson.A{"$$total", 1.1}}}},
	}}})
}

func TestSwitch(t *testing.T) {
	got := Switch(
		[]SwitchBranch{
			{Eq(Field("tier"), Value("gold")), Value(0.2)},
			{Eq(Field("tier"), Value("silver")), Value(0.1)},
		},
		Value(0),
	).Build()
	_ = mustJSON(got)
	assertBSON(t, got, bson.D{{Key: "$switch", Value: bson.D{
		{Key: "branches", Value: bson.A{
			bson.D{{Key: "case", Value: bson.D{{Key: "$eq", Value: bson.A{"$tier", "gold"}}}}, {Key: "then", Value: 0.2}},
			bson.D{{Key: "case", Value: bson.D{{Key: "$eq", Value: bson.A{"$tier", "silver"}}}}, {Key: "then", Value: 0.1}},
		}},
		{Key: "default", Value: 0},
	}}})
}

func TestRegexMatch(t *testing.T) {
	got := RegexMatch(Field("email"), `^.+@.+\..+$`, "i").Build()
	want := bson.D{{Key: "$regexMatch", Value: bson.D{
		{Key: "input", Value: "$email"},
		{Key: "regex", Value: `^.+@.+\..+$`},
		{Key: "options", Value: "i"},
	}}}
	assertBSON(t, got, want)
}
