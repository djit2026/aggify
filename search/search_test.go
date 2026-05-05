package search

import (
	"testing"
)

func TestSearch(t *testing.T) {
	got := Compound().Must(Text("title", "golang")).Build()

	if got[0].Key != "compound" {
		t.Errorf("expected compound, got %v", got[0].Key)
	}
}
