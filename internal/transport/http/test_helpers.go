package http

import (
	"encoding/json"
	"errors"
	"io"
)

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("error")
}

func makeJSON(data map[string]any) string {
	bytes, _ := json.Marshal(data)
	return string(bytes)
}

func parseBody(body io.Reader) map[string]any {
	b, _ := io.ReadAll(body)
	bodyMap := make(map[string]any)
	_ = json.Unmarshal(b, &bodyMap)
	return bodyMap
}
