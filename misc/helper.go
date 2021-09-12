package misc

import (
	"encoding/json"
)

func GetJsonFromJsonObjs(obj interface{}) ([]byte, error) {

	jsonBytes, err := json.Marshal(&obj)
	return jsonBytes, err
}
