package tools

type VectorResult struct {
	Entity string
	ID     int64
}

func GroupByEntityAndID(results []VectorResult) map[string][]int64 {
	group := map[string][]int64{}
	for _, r := range results {
		group[r.Entity] = append(group[r.Entity], r.ID)
	}
	return group
}
