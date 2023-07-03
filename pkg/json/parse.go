package json

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func Parse[T any](reader io.Reader) (*T, error) {
	jsonData := new(T)

	if err := json.NewDecoder(reader).Decode(jsonData); err != nil {
		return nil, err
	}

	return jsonData, nil
}

func Unmarshal[T any](data []byte) (*T, error) {
	unmarshaled := new(T)

	if err := json.Unmarshal(data, unmarshaled); err != nil {
		return nil, err
	}

	return unmarshaled, nil
}

func Read[T any](path string) (*T, error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	return Parse[T](jsonFile)
}
