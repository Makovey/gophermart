package http

import (
	"encoding/json"
	"errors"
)

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("error")
}

func makeJSON(data map[string]any) string {
	bytes, _ := json.Marshal(data)
	return string(bytes)
}
