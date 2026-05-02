package stage

import (
	"testing"

	"github.com/djit2026/aggify/expr"
	"github.com/djit2026/aggify/q"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func TestMatchStage(t *testing.T) {
	s := Match(q.Eq("status", "active"))
	got := s.Build()
	want := bson.D{{Key: "$match", Value: bson.D{{Key: "status", Value: "active"}}}}
	assert.Equal(t, want, got)
}

func TestGroupStage(t *testing.T) {
	s := Group(
		expr.Field("storeId"),
		Acc("total", expr.Sum(expr.Field("price"))),
	)
	got := s.Build()
	require.Equal(t, "$group", got[0].Key)

	inner := got[0].Value.(bson.D)
	assert.Equal(t, "$storeId", inner[0].Value)
	assert.Equal(t, "total", inner[1].Key)
}

func TestProjectStage(t *testing.T) {
	s := Project().Include("name", "email").Exclude("_id")
	got := s.Build()
	assert.Equal(t, "$project", got[0].Key)
	inner := got[0].Value.(bson.D)
	assert.Len(t, inner, 3)
}

func TestProjectComputed(t *testing.T) {
	s := Project().Computed("fullName",
		expr.Concat(expr.Field("first"), expr.Value(" "), expr.Field("last")),
	)
	got := s.Build()
	assert.Equal(t, "$project", got[0].Key)
}

func TestLookupStage(t *testing.T) {
	s := Lookup("users", "userId", "_id", "user")
	got := s.Build()
	assert.Equal(t, "$lookup", got[0].Key)
	inner := got[0].Value.(bson.D)
	assert.Equal(t, "users", inner[0].Value)
}

func TestUnwindStage(t *testing.T) {
	// Simple form.
	s := Unwind("$items")
	got := s.Build()
	assert.Equal(t, "$unwind", got[0].Key)
	assert.Equal(t, "$items", got[0].Value)

	// Full form.
	s2 := Unwind("$items").PreserveNullAndEmpty(true)
	got2 := s2.Build()
	assert.Equal(t, "$unwind", got2[0].Key)
	inner := got2[0].Value.(bson.D)
	assert.Equal(t, "$items", inner[0].Value)
}

func TestSortStage(t *testing.T) {
	s := Sort(SortField{"createdAt", Desc}, SortField{"name", Asc})
	got := s.Build()
	assert.Equal(t, "$sort", got[0].Key)
	inner := got[0].Value.(bson.D)
	assert.Equal(t, -1, inner[0].Value)
	assert.Equal(t, 1, inner[1].Value)
}

func TestLimitSkip(t *testing.T) {
	assert.Equal(t, bson.D{{Key: "$limit", Value: int64(10)}}, Limit(10).Build())
	assert.Equal(t, bson.D{{Key: "$skip", Value: int64(5)}}, Skip(5).Build())
}

func TestAddFieldsStage(t *testing.T) {
	s := AddFields(FE("doubled", expr.Multiply(expr.Field("price"), expr.Value(2))))
	got := s.Build()
	assert.Equal(t, "$addFields", got[0].Key)
}

func TestSetStage(t *testing.T) {
	s := Set(FE("upper", expr.ToUpper(expr.Field("name"))))
	got := s.Build()
	assert.Equal(t, "$set", got[0].Key)
}

func TestReplaceRootStage(t *testing.T) {
	s := ReplaceRoot(expr.MergeObjects(expr.Field("details"), expr.Root))
	got := s.Build()
	assert.Equal(t, "$replaceRoot", got[0].Key)
}

func TestUnsetStage(t *testing.T) {
	// Single field.
	s1 := Unset("password")
	got1 := s1.Build()
	assert.Equal(t, "$unset", got1[0].Key)
	assert.Equal(t, "password", got1[0].Value)

	// Multiple fields.
	s2 := Unset("password", "token")
	got2 := s2.Build()
	assert.Equal(t, "$unset", got2[0].Key)
	arr, ok := got2[0].Value.(bson.A)
	assert.True(t, ok)
	assert.Len(t, arr, 2)
}

func TestCountStage(t *testing.T) {
	s := Count("total")
	got := s.Build()
	assert.Equal(t, bson.D{{Key: "$count", Value: "total"}}, got)
}

func TestRawStage(t *testing.T) {
	d := bson.D{{Key: "$sample", Value: bson.D{{Key: "size", Value: 5}}}}
	s := Raw(d)
	assert.Equal(t, d, s.Build())
}

func TestPanicOnEmpty(t *testing.T) {
	assert.Panics(t, func() { Unset("") })
	assert.Panics(t, func() { Count("") })
	assert.Panics(t, func() { Lookup("", "a", "b", "c") })
}
