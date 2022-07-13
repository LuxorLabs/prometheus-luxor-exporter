package main

import (
	"strconv"
	"time"
)

func getStatsData(mi MinerInfo, usdPrice, btcPrice float64, user string, wallet Wallet) minerStatsAPIData {
	var minerStats minerStatsAPIData
	for idx := range mi.Data.Miners.Nodes {
		if mi.Data.Miners.Nodes[idx].User.Username == user {
			minerStats.Data.Timestamp = float64(time.Now().Unix())
			if float64(mi.Data.Miners.Nodes[idx].Details1H.UpdatedAt.Unix()) > minerStats.Data.LastSeenTimestamp {
				minerStats.Data.LastSeenTimestamp = float64(mi.Data.Miners.Nodes[idx].Details1H.UpdatedAt.Unix())
			}
			minerStats.Data.CurrentHashRate += toFloat(mi.Data.Miners.Nodes[idx].Details1H.Hashrate)
			minerStats.Data.AverageHashRate += toFloat(mi.Data.Miners.Nodes[idx].Details24H.Hashrate)
			minerStats.Data.ValidShares += toFloat(mi.Data.Miners.Nodes[idx].Details1H.ValidShares)
			minerStats.Data.InvalidShares += toFloat(mi.Data.Miners.Nodes[idx].Details1H.InvalidShares) +
				toFloat(mi.Data.Miners.Nodes[idx].Details1H.BadShares) +
				toFloat(mi.Data.Miners.Nodes[idx].Details1H.LowDiffShares) +
				toFloat(mi.Data.Miners.Nodes[idx].Details1H.DuplicateShares)
			minerStats.Data.StaleShares += toFloat(mi.Data.Miners.Nodes[idx].Details1H.StaleShares)
			revenuePerMinute := toFloat(mi.Data.Miners.Nodes[idx].Details1H.Revenue) / 60
			minerStats.Data.CoinsPerMinute += revenuePerMinute
			minerStats.Data.BTCPerMinute += revenuePerMinute * btcPrice
			minerStats.Data.USDPerMinute += revenuePerMinute * usdPrice
			if mi.Data.Miners.Nodes[idx].Details1H.Status == "Active" {
				minerStats.Data.ActiveWorkers++
			}

			unpaidBalance := float64(0)
			if len(wallet.Data.GetWallets.Nodes) > 0 {
				unpaidBalance = toFloat(wallet.Data.GetWallets.Nodes[0].PendingBalance)
			}
			minerStats.Data.UnpaidBalanceBaseUnits = unpaidBalance
			minerStats.Data.USDUnpaidBalance = unpaidBalance * usdPrice
		}
	}
	return minerStats
}

func getWorkerData(mi MinerInfo, user string) minerWorkersAPIData {
	workerStats := minerWorkersAPIData{
		Data: make([]minerWorkersAPIDataElement, 0, len(mi.Data.Miners.Nodes)),
	}

	for idx := range mi.Data.Miners.Nodes {
		if mi.Data.Miners.Nodes[idx].User.Username == user {
			workerStats.Data = append(workerStats.Data, minerWorkersAPIDataElement{
				Name:              mi.Data.Miners.Nodes[idx].WorkerName,
				Timestamp:         float64(time.Now().Unix()),
				LastSeenTimestamp: float64(mi.Data.Miners.Nodes[idx].Details1H.UpdatedAt.Unix()),
				CurrentHashRate:   toFloat(mi.Data.Miners.Nodes[idx].Details1H.Hashrate),
				ValidShares:       toFloat(mi.Data.Miners.Nodes[idx].Details1H.ValidShares),
				InvalidShares: toFloat(mi.Data.Miners.Nodes[idx].Details1H.InvalidShares) +
					toFloat(mi.Data.Miners.Nodes[idx].Details1H.BadShares) +
					toFloat(mi.Data.Miners.Nodes[idx].Details1H.LowDiffShares) +
					toFloat(mi.Data.Miners.Nodes[idx].Details1H.DuplicateShares),
				StaleShares: toFloat(mi.Data.Miners.Nodes[idx].Details1H.StaleShares),
			})
		}
	}
	return workerStats
}

func toFloat(val string) float64 {
	if val == "" || val == "null" {
		return 0
	}
	f, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return 0
	}
	return f
}
