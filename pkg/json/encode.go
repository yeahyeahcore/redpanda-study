package json

import (
	"bytes"
	"encoding/json"
)

func EncodeBuffer(body interface{}) (*bytes.Buffer, error) {
	buffer := new(bytes.Buffer)

	if err := json.NewEncoder(buffer).Encode(body); err != nil {
		return nil, err
	}

	return buffer, nil
}
