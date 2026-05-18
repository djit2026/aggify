package q

import (
	"testing"

	"github.com/djit2026/aggify/geo"
	"go.mongodb.org/mongo-driver/bson"
)

func ptrF(v float64) *float64 { return &v }

func TestGeoWithin(t *testing.T) {
	ring := []geo.Coord{
		geo.LngLat(-74.0, 40.7),
		geo.LngLat(-73.9, 40.7),
		geo.LngLat(-73.9, 40.8),
		geo.LngLat(-74.0, 40.8),
		geo.LngLat(-74.0, 40.7),
	}
	got := GeoWithin("location", geo.Polygon(ring))
	if got[0].Key != "location" {
		t.Errorf("expected field key 'location', got %s", got[0].Key)
	}
	inner := got[0].Value.(bson.D)
	if inner[0].Key != "$geoWithin" {
		t.Errorf("expected $geoWithin, got %s", inner[0].Key)
	}
}

func TestGeoWithinShape(t *testing.T) {
	got := GeoWithinShape("loc", geo.Box(geo.LngLat(-74, 40), geo.LngLat(-73, 41)))
	inner := got[0].Value.(bson.D)
	if inner[0].Key != "$geoWithin" {
		t.Errorf("expected $geoWithin, got %s", inner[0].Key)
	}
}

func TestGeoIntersects(t *testing.T) {
	got := GeoIntersects("area", geo.Point(-73.98, 40.75))
	inner := got[0].Value.(bson.D)
	if inner[0].Key != "$geoIntersects" {
		t.Errorf("expected $geoIntersects, got %s", inner[0].Key)
	}
}

func TestNear(t *testing.T) {
	max := 5000.0
	min := 100.0
	got := Near("location", geo.Point(-73.98, 40.75), NearOptions{
		MaxDistance: &max,
		MinDistance: &min,
	})
	inner := got[0].Value.(bson.D)
	if inner[0].Key != "$near" {
		t.Errorf("expected $near, got %s", inner[0].Key)
	}
	nearDoc := inner[0].Value.(bson.D)
	if nearDoc[0].Key != "$geometry" {
		t.Errorf("expected $geometry, got %s", nearDoc[0].Key)
	}
	if nearDoc[1].Key != "$maxDistance" {
		t.Errorf("expected $maxDistance, got %s", nearDoc[1].Key)
	}
}

func TestNearSphere(t *testing.T) {
	got := NearSphere("location", geo.Point(-73.98, 40.75))
	inner := got[0].Value.(bson.D)
	if inner[0].Key != "$nearSphere" {
		t.Errorf("expected $nearSphere, got %s", inner[0].Key)
	}
}
