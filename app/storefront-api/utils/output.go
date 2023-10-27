package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func PrintRequestBody(r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	var prettyJSON bytes.Buffer
	if err = json.Indent(&prettyJSON, bodyBytes, "", "\t"); err != nil {
		fmt.Printf("JSON parse error: %v", err)
		return
	}
	fmt.Println(string(prettyJSON.Bytes()))
}
