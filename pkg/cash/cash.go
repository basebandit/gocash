package cash

import (
	"encoding/json"
	"log"
)

//UnmarshalJSON decode dynamic json data to a (key,value) pair for internal use
func unmarshalJSON(data []byte) (map[string]interface{}, error) {
	var f map[string]interface{}
	err := json.Unmarshal(data, &f)
	if err != nil {
		return nil, err
	}
	return f, nil
}


