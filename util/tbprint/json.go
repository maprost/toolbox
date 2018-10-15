package mpprint

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// Json prints a json in a pretty way
func Json(js string) {
	JsonByteSlice([]byte(js))
}

// JsonByteSlice prints a json in a pretty way
func JsonByteSlice(js []byte) {
	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, js, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(prettyJSON.String())
}
