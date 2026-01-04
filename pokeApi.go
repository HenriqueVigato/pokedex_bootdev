package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func getDataJSON(url string) (map[string]any, error) {
	res, err := http.Get(url)
	if res.StatusCode > 299 {
		return nil, fmt.Errorf("response failed with a StatusCode: %v", res.StatusCode)
	}
	if err != nil {
		return nil, fmt.Errorf("error Creating a request to the api %v", err)
	}

	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("error readind the body response %v", err)
	}

	var dataJSON map[string]any
	if err := json.Unmarshal(body, &dataJSON); err != nil {
		return nil, fmt.Errorf("error during Unmarshal the response body %v", err)
	}

	return dataJSON, nil
}
