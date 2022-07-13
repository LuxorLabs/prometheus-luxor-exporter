package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func runServer() error {
	fmt.Printf("Listening on %s.\n", endpoint)
	var mainServeMux http.ServeMux
	mainServeMux.HandleFunc("/", handleOtherRequest)
	mainServeMux.HandleFunc("/miner", handleMinerScrapeRequest)
	if err := http.ListenAndServe(endpoint, &mainServeMux); err != nil {
		return fmt.Errorf("error while running main HTTP server: %s", err)
	}
	return nil
}

func handleOtherRequest(response http.ResponseWriter, request *http.Request) {
	if request.URL.Path == "/" {
		fmt.Fprintf(response, "\nPool IDs:\n")
		for poolID := range Pools {
			fmt.Fprintf(response, "- %s\n", poolID)
		}
		fmt.Fprintf(response, "\nMetrics paths:\n")
		fmt.Fprintf(response, "- Miner: /miner?pool=<pool>&target=<miner-address>\n")
	} else {
		message := fmt.Sprintf("404 - Page not found.\n")
		http.Error(response, message, 404)
	}
}

func handleMinerScrapeRequest(response http.ResponseWriter, request *http.Request) {
	if enableDebug {
		fmt.Printf("[DEBUG] Request: endpoint=%s from=%s to=%v\n", "miner", request.RemoteAddr, request.URL.String())
	}

	// Get pool
	var poolID string
	if values, ok := request.URL.Query()["pool"]; ok && len(values) > 0 && values[0] != "" {
		poolID = values[0]
	} else {
		http.Error(response, "400 - Missing pool.\n", 400)
		return
	}
	pool, poolOK := Pools[poolID]
	if !poolOK {
		http.Error(response, "404 - Pool not found.\n", 404)
		return
	}

	// Get apikey
	var apikey string
	if values, ok := request.URL.Query()["apikey"]; ok && len(values) > 0 && values[0] != "" {
		apikey = values[0]
	} else {
		http.Error(response, "400 - Missing miner address.\n", 400)
		return
	}

	// Get subaccount
	var subaccount string
	if values, ok := request.URL.Query()["subaccount"]; ok && len(values) > 0 && values[0] != "" {
		subaccount = values[0]
	} else {
		http.Error(response, "400 - Missing subaccount name.\n", 400)
		return
	}

	minerInfo := getMinersInfo(response, apikey)
	usdPrice, btcPrice, _ := GetPrices()

	pendingBalance := getPendingBalance(response, subaccount, apikey)
	statsData := getStatsData(minerInfo, usdPrice, btcPrice, subaccount, pendingBalance)
	workersData := getWorkerData(minerInfo, subaccount)
	registry := buildMinerRegistry(response, &pool, subaccount, &statsData, &workersData)

	if registry == nil {
		return
	}

	// Delegare final handling to Prometheus
	handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	handler.ServeHTTP(response, request)
}
