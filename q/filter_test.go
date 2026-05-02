package q

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestEq(t *testing.T) {
	got := Eq("status", "active")
	want := bson.D{{Key: "status", Value: "active"}}
	assert.Equal(t, want, got)
}

func TestNe(t *testing.T) {
	got := Ne("count", 0)
	want := bson.D{{Key: "count", Value: bson.D{{Key: "$ne", Value: 0}}}}
	assert.Equal(t, want, got)
}

func TestGtGteLtLte(t *testing.T) {
	assert.Equal(t, bson.D{{Key: "age", Value: bson.D{{Key: "$gt", Value: 18}}}}, Gt("age", 18))
	assert.Equal(t, bson.D{{Key: "age", Value: bson.D{{Key: "$gte", Value: 21}}}}, Gte("age", 21))
	assert.Equal(t, bson.D{{Key: "stock", Value: bson.D{{Key: "$lt", Value: 10}}}}, Lt("stock", 10))
	assert.Equal(t, bson.D{{Key: "score", Value: bson.D{{Key: "$lte", Value: 100}}}}, Lte("score", 100))
}

func TestIn(t *testing.T) {
	got := In("role", "admin", "editor")
	want := bson.D{{Key: "role", Value: bson.D{{Key: "$in", Value: bson.A{"admin", "editor"}}}}}
	assert.Equal(t, want, got)
}

func TestNin(t *testing.T) {
	got := Nin("status", "banned", "deleted")
	want := bson.D{{Key: "status", Value: bson.D{{Key: "$nin", Value: bson.A{"banned", "deleted"}}}}}
	assert.Equal(t, want, got)
}

func TestExists(t *testing.T) {
	got := Exists("deletedAt", false)
	want := bson.D{{Key: "deletedAt", Value: bson.D{{Key: "$exists", Value: false}}}}
	assert.Equal(t, want, got)
}

func TestElemMatch(t *testing.T) {
	got := ElemMatch("items", And(Eq("status", "active"), Gt("qty", 0)))
	assert.Equal(t, "items", got[0].Key)
}

func TestAnd(t *testing.T) {
	got := And(Eq("a", 1), Eq("b", 2))
	assert.Equal(t, "$and", got[0].Key)
	arr, ok := got[0].Value.(bson.A)
	assert.True(t, ok)
	assert.Len(t, arr, 2)
}

func TestOr(t *testing.T) {
	got := Or(Eq("role", "admin"), Eq("role", "mod"))
	assert.Equal(t, "$or", got[0].Key)
}

func TestNor(t *testing.T) {
	got := Nor(Eq("status", "banned"))
	assert.Equal(t, "$nor", got[0].Key)
}
