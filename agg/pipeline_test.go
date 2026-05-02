package agg

import (
	"testing"

	"github.com/djit2026/aggify/expr"
	"github.com/djit2026/aggify/q"
	"github.com/djit2026/aggify/stage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func TestNewPipeline(t *testing.T) {
	p := New()
	assert.Equal(t, 0, p.Len())
	assert.Empty(t, p.Build())
}

func TestMatchAndGroup(t *testing.T) {
	p := New().
		Match(q.Eq("status", "active")).
		Group(
			expr.Field("storeId"),
			stage.Acc("total", expr.Sum(expr.Field("price"))),
			stage.Acc("count", expr.Sum(expr.Value(1))),
		)

	pipeline := p.Build()
	require.Len(t, pipeline, 2)
	assert.Equal(t, "$match", pipeline[0][0].Key)
	assert.Equal(t, "$group", pipeline[1][0].Key)
}

func TestSortLimitSkip(t *testing.T) {
	p := New().SortDesc("createdAt").Skip(10).Limit(5)
	pipeline := p.Build()
	require.Len(t, pipeline, 3)
	assert.Equal(t, "$sort", pipeline[0][0].Key)
	assert.Equal(t, "$skip", pipeline[1][0].Key)
	assert.Equal(t, "$limit", pipeline[2][0].Key)
}

func TestUnwindAndProject(t *testing.T) {
	p := New().
		Unwind("$items").
		Stage(stage.Project().Include("name", "price").Exclude("_id"))
	pipeline := p.Build()
	require.Len(t, pipeline, 2)
	assert.Equal(t, "$unwind", pipeline[0][0].Key)
	assert.Equal(t, "$project", pipeline[1][0].Key)
}

func TestSetAndUnset(t *testing.T) {
	p := New().
		Set(stage.FE("upper", expr.ToUpper(expr.Field("name")))).
		Unset("_id", "internalField")
	pipeline := p.Build()
	require.Len(t, pipeline, 2)
	assert.Equal(t, "$set", pipeline[0][0].Key)
	assert.Equal(t, "$unset", pipeline[1][0].Key)
}

func TestRawEscapeHatch(t *testing.T) {
	p := New().Raw(bson.D{{Key: "$sample", Value: bson.D{{Key: "size", Value: 5}}}})
	pipeline := p.Build()
	require.Len(t, pipeline, 1)
	assert.Equal(t, "$sample", pipeline[0][0].Key)
}

func TestMustJSON(t *testing.T) {
	p := New().Match(q.Eq("active", true)).Limit(10)
	json := p.MustJSON()
	assert.Contains(t, json, "$match")
	assert.Contains(t, json, "$limit")
}

// TestComposableStage demonstrates the reusable stage pattern.
func TestComposableStage(t *testing.T) {
	// Define a reusable stage as a plain function — no library ceremony needed.
	activeItemsStage := func() stage.Stage {
		return stage.AddFields(
			stage.FE("activeItems", expr.Filter(
				expr.Field("items"),
				"item",
				expr.Eq(expr.Var("item.status"), expr.Value("active")),
			)),
		)
	}

	p := New().
		Match(q.Eq("userId", "user123")).
		Stage(activeItemsStage()).
		Count("total")

	pipeline := p.Build()
	require.Len(t, pipeline, 3)
	assert.Equal(t, "$match", pipeline[0][0].Key)
	assert.Equal(t, "$addFields", pipeline[1][0].Key)
	assert.Equal(t, "$count", pipeline[2][0].Key)
}
