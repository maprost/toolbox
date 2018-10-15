package tbarray

func RemoveInt64Duplicates(list *[]int64) {
	found := make(map[int64]struct{})
	newListSize := 0
	for i, x := range *list {
		if _, ok := found[x]; !ok {
			found[x] = struct{}{}
			(*list)[newListSize] = (*list)[i]
			newListSize++
		}
	}
	*list = (*list)[:newListSize]
}

func RemoveStringDuplicates(list *[]string) {
	found := make(map[string]struct{})
	newListSize := 0
	for i, x := range *list {
		if _, ok := found[x]; !ok {
			found[x] = struct{}{}
			(*list)[newListSize] = (*list)[i]
			newListSize++
		}
	}
	*list = (*list)[:newListSize]
}
