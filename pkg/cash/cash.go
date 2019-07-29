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


//ParseRates retrieves rates object from the decoded json response
func ParseRates(data []byte) map[string]interface{} {
	var r map[string]interface{}
	res, err := unmarshalJSON(data)
	if err != nil {
		log.Fatal(err)
	}
	if rates, ok := res["rates"]; ok {
		r = rates.(map[string]interface{})
	}
	return r
}


//ParseBase retrieves base currency value from the decoded json response
func ParseBase(data []byte) string {
	var b string
	res, err := unmarshalJSON(data)
	if err != nil {
		log.Fatal(err)
	}
	if base, ok := res["base"]; ok {
		b = base.(string)
	}
	return b
}