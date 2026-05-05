package stage

import (
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

func TestOut(t *testing.T) {
	got := Out("authors").Build()
	want := bson.D{{Key: "$out", Value: "authors"}}
	if got[0].Key != want[0].Key || got[0].Value != want[0].Value {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestMerge(t *testing.T) {
	got := Merge(MergeOptions{
		IntoCollection: "newDailySales",
		On:             []string{"date"},
		WhenMatched:    "replace",
		WhenNotMatched: "insert",
	}).Build()

	if got[0].Key != "$merge" {
		t.Errorf("got %v, want $merge", got[0].Key)
	}
}
