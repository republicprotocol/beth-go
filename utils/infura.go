package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// SendRequest will send a request to infura and return the unmarshalled data
// back to the caller. It will retry until a valid response is returned, or
// until the context times out.
func SendRequest(
	ctx context.Context,
	url string,
	request string,
	data interface{},
) (interface{}, error) {

	// Retry until a valid response is returned or until context times out
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(5 * time.Millisecond):
		}

		// Create a new http  POST request
		req, err := http.NewRequest(
			"POST",
			url,
			bytes.NewBuffer([]byte(request)),
		)
		if err != nil {
			continue
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			continue
		}

		// Read the response status
		if resp.StatusCode != http.StatusOK {
			if resp.Body != nil {
				resp.Body.Close()
			}
			continue
		}

		// Get the result
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			if resp.Body != nil {
				resp.Body.Close()
			}
			continue
		}

		if resp.Body != nil {
			resp.Body.Close()
		}
		return data, nil
	}
}
