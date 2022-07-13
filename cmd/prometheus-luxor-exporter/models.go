package main

import (
	"strings"
	"time"
)

// CurrencySymbol - The symbol of a currency, e.g. ETH for Ethereum.
type CurrencySymbol string

// Currency - Currency info.
type Currency struct {
	Currency CurrencySymbol
	// E.g. number of weis in 1 ether.
	BaseUnitsPerUnit float64
}

// Currencies - List of currencies used by the supported pools.
var Currencies = map[CurrencySymbol]Currency{
	CurrencySymbolEthereum: {CurrencySymbolEthereum, 1e18},
}

// Pool - Pool info.
type Pool struct {
	ID       string
	Name     string
	Currency CurrencySymbol
	// URL must not have a trailing slash.
	APIURL string
}

// Pools - List of supported pools.
var Pools = map[string]Pool{
	"luxor": {"luxor", "Luxor", CurrencySymbolEthereum, "https://api.beta.luxor.tech/graphql"},
}

type baseAPIData struct {
	Status string `json:"status"`
}

type minerStatsAPIData struct {
	baseAPIData
	Data struct {
		Timestamp              float64 `json:"time"`
		LastSeenTimestamp      float64 `json:"lastSeen"`
		CurrentHashRate        float64 `json:"currentHashrate"`
		AverageHashRate        float64 `json:"averageHashrate"`
		ValidShares            float64 `json:"validShares"`
		InvalidShares          float64 `json:"invalidShares"`
		StaleShares            float64 `json:"staleShares"`
		ActiveWorkers          float64 `json:"activeWorkers"`
		CoinsPerMinute         float64 `json:"coinsPerMin"`
		BTCPerMinute           float64 `json:"btcPerMin"`
		USDPerMinute           float64 `json:"usdPerMin"`
		UnpaidBalanceBaseUnits float64 `json:"unpaid"`
		USDUnpaidBalance       float64 `json:"usdUnpaid"`
	} `json:"data"`
}

type GraphQLTime struct {
	time.Time
}

const ctLayout = "2006-01-02T15:04:05.99"

func (ct *GraphQLTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		ct.Time = time.Time{}
		return
	}
	ct.Time, err = time.Parse(ctLayout, s)
	return
}

var nilTime = (time.Time{}).UnixNano()

func (ct *GraphQLTime) IsSet() bool {
	return ct.UnixNano() != nilTime
}

type MinerInfo struct {
	Data struct {
		Miners struct {
			Nodes []struct {
				Deleted           bool   `json:"deleted"`
				MiningProfileName string `json:"miningProfileName"`
				User              struct {
					Username string `json:"username"`
					RowID    int    `json:"rowId"`
				} `json:"user"`
				RowID      int    `json:"rowId"`
				WorkerName string `json:"workerName"`
				Details1H  struct {
					BadShares       string      `json:"badShares"`
					DuplicateShares string      `json:"duplicateShares"`
					Efficiency      string      `json:"efficiency"`
					Hashrate        string      `json:"hashrate"`
					LowDiffShares   string      `json:"lowDiffShares"`
					InvalidShares   string      `json:"invalidShares"`
					Revenue         string      `json:"revenue"`
					StaleShares     string      `json:"staleShares"`
					Status          string      `json:"status"`
					UpdatedAt       GraphQLTime `json:"updatedAt"`
					ValidShares     string      `json:"validShares"`
				} `json:"details1H"`
				Details24H struct {
					Hashrate string `json:"hashrate"`
				} `json:"details24H"`
			} `json:"nodes"`
		} `json:"miners"`
	} `json:"data"`
}

type minerWorkersAPIData struct {
	baseAPIData
	Data []minerWorkersAPIDataElement `json:"data"`
}

type minerWorkersAPIDataElement struct {
	Name              string  `json:"worker"`
	Timestamp         float64 `json:"time"`
	LastSeenTimestamp float64 `json:"lastSeen"`
	CurrentHashRate   float64 `json:"currentHashrate"`
	ValidShares       float64 `json:"validShares"`
	InvalidShares     float64 `json:"invalidShares"`
	StaleShares       float64 `json:"staleShares"`
}

type Wallet struct {
	Data struct {
		GetWallets struct {
			Nodes []struct {
				PendingBalance string `json:"pendingBalance"`
			} `json:"nodes"`
		} `json:"getWallets"`
	} `json:"data"`
}
