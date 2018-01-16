package print

import "encoding/json"

// Struct prints a struct in a pretty way
func Struct(obj interface{}) {
	s, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}

	JsonByteSlice(s)
}
