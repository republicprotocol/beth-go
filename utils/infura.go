package utils

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// SendRequest will send a request to infura and return the unmarshalled data
// back to the caller. It will retry until a valid response is returned, or
// until the context times out.
func SendRequest(ctx context.Context, url string, request string) ([]byte, error) {

	// Retry until a valid response is returned or until context times out
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(5 * time.Millisecond):
		}

		// Create a new http POST request
		req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(request)))
		if err != nil {
			continue
		}

		// Send http POST request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			continue
		}

		// Decode response body
		body, err := func() ([]byte, error) {
			defer resp.Body.Close()

			// Check status
			if resp.StatusCode != http.StatusOK {
				return nil, fmt.Errorf("unexpected status %v", resp.StatusCode)
			}
			// Check body
			if resp.Body != nil {
				return ioutil.ReadAll(resp.Body)
			}
			return nil, fmt.Errorf("response body is nil")
		}()
		if err != nil {
			fmt.Printf("[error] (infura) %v", err)
			continue
		}
		return body, nil
	}
}
