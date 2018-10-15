package tbarray

func CompareStringList(listA []string, listB []string) []string {
	equal := make(map[string]struct{})

	for _, v := range listA {
		equal[v] = struct{}{}
	}

	newListSize := 0
	for i, x := range listB {
		if _, ok := equal[x]; ok {
			listB[newListSize] = listB[i]
			newListSize++
		}
	}
	return listB[:newListSize]
}

// exclude all elements of list B out of list A and return the only the element that are in list A.
// e.g. ListA={1,2,3,4} ListB={3,4,5,6} -> Result={1,2}
func Exclude(list []string, exclude []string) []string {
	toExclude := make(map[string]struct{})

	for _, v := range exclude {
		toExclude[v] = struct{}{}
	}

	newListSize := 0
	for i, x := range list {
		if _, ok := toExclude[x]; !ok {
			list[newListSize] = list[i]
			newListSize++
		}
	}
	return list[:newListSize]
}
