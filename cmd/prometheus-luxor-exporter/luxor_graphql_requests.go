package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gitlab.com/luxorlabs/prometheus-luxor-exporter/util"
	"io/ioutil"
	"net/http"
	"time"
)

func getMinersInfo(response http.ResponseWriter, apikey string) MinerInfo {
	jsonData := map[string]string{
		"query": `
            {
			miners (first: 20000 filter: {miningProfileName: {equalTo: ETH}} orderBy: ID_ASC) {
            nodes { 
                workerName     
                user { username }  
                details1H {       badShares       duplicateShares       efficiency       hashrate       lowDiffShares       invalidShares       revenue       staleShares       status       updatedAt       validShares     }
				details24H { hashrate }
            }
        }
    }
        `,
	}
	jsonValue, _ := json.Marshal(jsonData)
	request, err := http.NewRequest("POST", "https://api.beta.luxor.tech/graphql", bytes.NewBuffer(jsonValue))
	request.Header.Set("X-Lux-Api-Key", apikey)
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: time.Second * 100}
	resp, err := client.Do(request)
	defer resp.Body.Close()
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}
	rawdata, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(rawdata))

	var data MinerInfo
	util.ParseJSON(&data, response, rawdata, false, enableDebug)

	return data
}

func getPendingBalance(response http.ResponseWriter, user, apikey string) Wallet {
	jsonData := map[string]string{
		"query": fmt.Sprintf(`
            {
			getWallets(first: 100 uname:"%s" filter: {currencyProfileName: {equalTo: ETH}}) {
            nodes { 
                pendingBalance
            }
        }
    }`, user),
	}
	jsonValue, _ := json.Marshal(jsonData)
	request, err := http.NewRequest("POST", "https://api.beta.luxor.tech/graphql", bytes.NewBuffer(jsonValue))
	request.Header.Set("X-Lux-Api-Key", apikey)
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: time.Second * 100}
	resp, err := client.Do(request)
	defer resp.Body.Close()
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}
	rawdata, _ := ioutil.ReadAll(resp.Body)

	var wallet Wallet
	util.ParseJSON(&wallet, response, rawdata, false, enableDebug)

	return wallet
}
