package utils

import "encoding/json"

type JSONMap map[string]any

func JSONToMap(j []byte) (out map[string]any) {
	err := json.Unmarshal(j, &out)
	if err != nil {
		println(err)
	}

	return out
}
