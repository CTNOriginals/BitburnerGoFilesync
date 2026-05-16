package utils

import (
	"encoding/json"
	"log"
)

type JSONMap map[string]any

func JSONToMap(j []byte) (out map[string]any) {
	err := json.Unmarshal(j, &out)
	if err != nil {
		log.Println(err)
	}

	return out
}
