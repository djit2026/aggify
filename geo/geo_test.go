package geo

import (
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

func assertBSON(t *testing.T, got, want any) {
	t.Helper()
	g, err := bson.MarshalExtJSON(got, true, true)
	if err != nil {
		t.Fatalf("marshal got: %v", err)
	}
	w, err := bson.MarshalExtJSON(want, true, true)
	if err != nil {
		t.Fatalf("marshal want: %v", err)
	}
	if string(g) != string(w) {
		t.Errorf("\ngot  %s\nwant %s", g, w)
	}
}

func TestPoint(t *testing.T) {
	got := Point(-73.9857, 40.7580).Build()
	assertBSON(t, got, bson.D{
		{Key: "type", Value: "Point"},
		{Key: "coordinates", Value: bson.A{-73.9857, 40.7580}},
	})
}

func TestLineString(t *testing.T) {
	got := LineString(LngLat(-73.98, 40.75), LngLat(-74.00, 40.71)).Build()
	assertBSON(t, got, bson.D{
		{Key: "type", Value: "LineString"},
		{Key: "coordinates", Value: bson.A{
			bson.A{-73.98, 40.75},
			bson.A{-74.00, 40.71},
		}},
	})
}

func TestPolygon(t *testing.T) {
	ring := []Coord{
		LngLat(-73.99, 40.76),
		LngLat(-73.98, 40.76),
		LngLat(-73.98, 40.75),
		LngLat(-73.99, 40.76),
	}
	got := Polygon(ring).Build()
	if got[0].Value.(string) != "Polygon" {
		t.Errorf("expected Polygon type")
	}
}

func TestBox(t *testing.T) {
	got := Box(LngLat(-74.00, 40.71), LngLat(-73.96, 40.75)).Build()
	assertBSON(t, got, bson.D{{Key: "$box", Value: bson.A{
		bson.A{-74.00, 40.71},
		bson.A{-73.96, 40.75},
	}}})
}

func TestCenterSphere(t *testing.T) {
	got := CenterSphere(LngLat(-73.98, 40.75), 5000.0/6378100.0).Build()
	if got[0].Key != "$centerSphere" {
		t.Errorf("expected $centerSphere, got %s", got[0].Key)
	}
}

func TestGeometryCollection(t *testing.T) {
	got := GeometryCollection(
		Point(-73.98, 40.75),
		Point(-74.00, 40.71),
	).Build()
	if got[0].Value.(string) != "GeometryCollection" {
		t.Errorf("expected GeometryCollection type")
	}
}
