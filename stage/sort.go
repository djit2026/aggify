package stage

import "go.mongodb.org/mongo-driver/bson"

// SortOrder is the sort direction.
type SortOrder int

const (
	// Asc sorts ascending (1).
	Asc SortOrder = 1
	// Desc sorts descending (-1).
	Desc SortOrder = -1
)

// SortField is a single sort key.
type SortField struct {
	Field string
	Order SortOrder
}

// sortStage implements Stage for $sort.
type sortStage struct{ fields []SortField }

// Sort creates a $sort stage.
// Fields are applied in the order provided (which matters for compound sorts).
//
//	stage.Sort(stage.SortField{"createdAt", stage.Desc}, stage.SortField{"name", stage.Asc})
func Sort(fields ...SortField) Stage {
	return sortStage{fields}
}

// SortAsc is a convenience function for a single ascending sort.
func SortAsc(field string) Stage {
	return sortStage{[]SortField{{field, Asc}}}
}

// SortDesc is a convenience function for a single descending sort.
func SortDesc(field string) Stage {
	return sortStage{[]SortField{{field, Desc}}}
}

func (s sortStage) Build() bson.D {
	doc := make(bson.D, len(s.fields))
	for i, f := range s.fields {
		doc[i] = bson.E{Key: f.Field, Value: int(f.Order)}
	}
	return bson.D{{Key: "$sort", Value: doc}}
}
