package helpers

import "sort"

type MapGroup struct {
	Key   string
	Value []string
}

func SortMap(m map[string][]string) []MapGroup {
	var groups []MapGroup
	for key, value := range m {
		groups = append(groups, MapGroup{
			Key:   key,
			Value: value,
		})
	}

	sort.Slice(groups, func(i, j int) bool {
		// shortest key first
		if len(groups[i].Key) != len(groups[j].Key) {
			return len(groups[i].Key) < len(groups[j].Key)
		}
		// if same length, sort alphabetically
		return groups[i].Key < groups[j].Key
	})

	return groups
}
