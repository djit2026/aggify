package stage

import (
	"testing"

	"github.com/djit2026/aggify/geo"
	"go.mongodb.org/mongo-driver/bson"
)

func ptrFloat(v float64) *float64 { return &v }

func TestGeoNear(t *testing.T) {
	got := GeoNear(GeoNearOptions{
		Near:          geo.Point(-73.9857, 40.7580),
		DistanceField: "dist.meters",
		Spherical:     true,
		MaxDistance:   ptrFloat(2000),
	}).Build()

	if got[0].Key != "$geoNear" {
		t.Errorf("expected $geoNear stage key, got %s", got[0].Key)
	}

	doc := got[0].Value.(bson.D)

	// Validate distanceField
	found := false
	for _, e := range doc {
		if e.Key == "distanceField" && e.Value == "dist.meters" {
			found = true
		}
	}
	if !found {
		t.Errorf("distanceField not set correctly in $geoNear")
	}
}

func TestGeoNearWithQuery(t *testing.T) {
	filter := bson.D{{Key: "category", Value: "restaurant"}}
	got := GeoNear(GeoNearOptions{
		Near:          geo.Point(-73.9857, 40.7580),
		DistanceField: "dist.meters",
		Spherical:     true,
		Query:         filter,
		Key:           "location",
	}).Build()

	doc := got[0].Value.(bson.D)
	keyFound := false
	for _, e := range doc {
		if e.Key == "key" && e.Value == "location" {
			keyFound = true
		}
	}
	if !keyFound {
		t.Errorf("key field not set correctly in $geoNear")
	}
}

func TestGeoNearPanicsOnMissingDistanceField(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic when distanceField is empty")
		}
	}()
	GeoNear(GeoNearOptions{Near: geo.Point(0, 0)})
}
