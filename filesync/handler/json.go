package handler

import "encoding/json"

func JSONToMap(j []byte) (out *map[string]any) {
	err := json.Unmarshal(j, &out)
	if err != nil {
		println(err)
	}

	return out
}
