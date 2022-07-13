package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

var (
	lastUSDPrice, lastBTCPrice float64
)

// GetPrices is for getting the price in usd and btc.
func GetPrices() (usd float64, btc float64, err error) {
	params := url.Values{}
	// Yes, all these params seem to be required
	params.Set("ids", "ethereum")
	params.Set("vs_currencies", "btc,usd")
	requestStr := fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?%s", params.Encode())
	request, _ := http.NewRequest("GET", requestStr, nil)
	httpClient := http.Client{}
	response, err := httpClient.Do(request)
	if err != nil {
		return lastBTCPrice, lastUSDPrice, err
	}
	if response != nil {
		defer response.Body.Close()
	}
	if response.Body == nil || response.StatusCode >= 400 {
		return lastUSDPrice, lastBTCPrice, fmt.Errorf("unable to get coin prices. %d", response.StatusCode)
	}

	var priceData map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&priceData)
	if err != nil {
		return lastUSDPrice, lastBTCPrice, err
	}

	coinData := priceData["ethereum"]
	if coinData == nil {
		return lastUSDPrice, lastBTCPrice, fmt.Errorf("unable to get pricedata for coin. Pricedata: %v", priceData)
	}
	priceMap := coinData.(map[string]interface{})

	// Updating last known values so we can use them if the api fail.
	lastUSDPrice = priceMap["usd"].(float64)
	lastBTCPrice = priceMap["btc"].(float64)
	return priceMap["usd"].(float64), priceMap["btc"].(float64), nil
}
