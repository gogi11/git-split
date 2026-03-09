package helpers

import "sort"

type MapGroup[T any] struct {
	Key   string
	Value []T
}

func SortMap[T any](m map[string][]T) []MapGroup[T] {
	var groups []MapGroup[T]
	for key, value := range m {
		groups = append(groups, MapGroup[T]{
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
