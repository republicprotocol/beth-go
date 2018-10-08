package netutils

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// SendInfuraRequest will send a request to infura and return the unmarshalled data
// back to the caller. It will retry until a valid response is returned, or
// until the context times out.
func SendInfuraRequest(ctx context.Context, url string, request string) (body []byte, err error) {

	sleepDurationMs := time.Duration(1000)

	// Retry until a valid response is returned or until context times out
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		if body, err = func() ([]byte, error) {
			// Create a new http POST request
			req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(request)))
			if err != nil {
				return nil, err
			}

			// Send http POST request
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				return nil, err
			}

			// Decode response body
			return func() ([]byte, error) {
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
		}(); err == nil {
			break
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(sleepDurationMs * time.Millisecond):

		}

		// Increase delay for next round but saturate at 30s
		sleepDurationMs = time.Duration(float64(sleepDurationMs) * 1.6)
		if sleepDurationMs > 30000 {
			sleepDurationMs = 30000
		}
	}
	return
}
