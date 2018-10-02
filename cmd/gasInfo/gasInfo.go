package main

import (
	"encoding/json"
	"errors"
	"log"
	"math"
	"math/big"
	"net/http"
)

type resp struct {
	Fast string `json:"fast"`
}

func main() {
	print(GetGasPrice())
}

// GetGasPrice
func GetGasPrice() error {
	request, _ := http.NewRequest("GET", "https://ethgasstation.info/json/ethgasAPI.json", nil)

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Println("here")
		return err
	}

	type resp struct {
		Fast float64 `json:"fast"`
	}

	log.Println(response.Status)

	data := new(resp)
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		log.Println("nope")
		return err
	}
	// gasPrice, err := strconv.Atoi(fmt.Sprintf("%", data["fast"]))
	log.Println(big.NewInt(int64(data.Fast * math.Pow10(8))))

	return errors.New("cannot get gas price")
}
