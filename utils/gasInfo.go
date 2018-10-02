package utils

import (
	"encoding/json"
	"log"
	"math"
	"math/big"
	"net/http"
)

// SuggestedGasPrice returns the fast gas price that ethGasStation
// recommends for transactions to be mined on Ethereum blockchain.
func SuggestedGasPrice() *big.Int {
	request, _ := http.NewRequest("GET",
		"https://ethgasstation.info/json/ethgasAPI.json", nil)

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Printf("cannot connect to ethGasStationAPI: %v", err)
		return nil
	}

	if response.StatusCode != http.StatusOK {
		log.Printf("received status code %v from ethGasStation",
			response.StatusCode)
		return nil
	}

	type resp struct {
		Fast float64 `json:"fast"`
	}

	data := new(resp)

	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		log.Printf("cannot decode json response from ethGasStation: %v", err)
		return nil
	}
	return big.NewInt(int64(data.Fast * math.Pow10(8)))
}
