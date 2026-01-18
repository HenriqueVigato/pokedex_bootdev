package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func getData(url string) ([]byte, error) {
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

	return body, nil
}

func convertToJSON(data []byte) map[string]any {
	var dataJSON map[string]any
	if err := json.Unmarshal(data, &dataJSON); err != nil {
		fmt.Errorf("error during Unmarshal the response body %v", err)
		return nil
	}

	return dataJSON
}
