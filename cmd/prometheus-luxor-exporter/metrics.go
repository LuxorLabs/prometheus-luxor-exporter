package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"gitlab.com/luxorlabs/prometheus-luxor-exporter/util"
	"net/http"
)

// Builds a new registry for the miner endpoint, adds scraped data to it and returns it if successful or nil if not.
func buildMinerRegistry(response http.ResponseWriter, pool *Pool, username string, statsData *minerStatsAPIData, workersData *minerWorkersAPIData) *prometheus.Registry {
	registry := prometheus.NewRegistry()
	registry.MustRegister(prometheus.NewGoCollector())

	util.NewExporterMetric(registry, namespace, appVersion)

	// Note: Miner address isn't needed as it's the instance/target of the scrape.
	constLabels := prometheus.Labels{
		"pool":       pool.ID,
		"subaccount": username,
	}
	constLabelsWithCurrency := util.MergeLabels(constLabels, prometheus.Labels{
		"currency": string(pool.Currency),
	})

	// Miner info
	util.NewGauge(registry, namespace, "miner", "info", "Metadata about the miner.", util.MergeLabels(constLabels, prometheus.Labels{
		"pool_name":     pool.Name,
		"pool_currency": string(pool.Currency),
	})).Set(1)

	// Miner stats
	util.NewGauge(registry, namespace, "miner", "last_seen_seconds", "Delta between time of last statistics entry and when any workers from the miner was last seen (s).", constLabels).Set(statsData.Data.Timestamp - statsData.Data.LastSeenTimestamp)
	util.NewGauge(registry, namespace, "miner", "hashrate_current_hps", "Total current hash rate for a miner (H/s).", constLabels).Set(statsData.Data.CurrentHashRate)
	util.NewGauge(registry, namespace, "miner", "hashrate_average_hps", "Total average hash rate for a miner (H/s).", constLabels).Set(statsData.Data.AverageHashRate)
	util.NewGauge(registry, namespace, "miner", "shares_valid", "Total number of valid shares for a miner.", constLabels).Set(statsData.Data.ValidShares)
	util.NewGauge(registry, namespace, "miner", "shares_invalid", "Total number of invalid shares for a miner.", constLabels).Set(statsData.Data.InvalidShares)
	util.NewGauge(registry, namespace, "miner", "shares_stale", "Total number of stale shares for a miner.", constLabels).Set(statsData.Data.StaleShares)
	util.NewGauge(registry, namespace, "miner", "workers_active", "Number of active workers.", constLabels).Set(statsData.Data.ActiveWorkers)
	util.NewGauge(registry, namespace, "miner", "income_coins", "Mined coins per second.", constLabelsWithCurrency).Set(statsData.Data.CoinsPerMinute / 60)
	util.NewGauge(registry, namespace, "miner", "income_usd", "Mined coins per second (converted to USD).", constLabels).Set(statsData.Data.USDPerMinute / 60)
	util.NewGauge(registry, namespace, "miner", "income_btc", "Mined coins per second (converted to BTC).", constLabels).Set(statsData.Data.BTCPerMinute / 60)
	util.NewGauge(registry, namespace, "miner", "balance_unpaid_coins", "Unpaid balance for a miner.", constLabelsWithCurrency).Set(statsData.Data.UnpaidBalanceBaseUnits)
	util.NewGauge(registry, namespace, "miner", "balance_unpaid_usd", "Unpaid balance for a miner.", constLabels).Set(statsData.Data.UnpaidBalanceBaseUnits)
	// Deprecated
	util.NewGauge(registry, namespace, "miner", "income_minute_coins", "(Deprecated) Mined coins per minute.", constLabelsWithCurrency).Set(statsData.Data.CoinsPerMinute)
	util.NewGauge(registry, namespace, "miner", "income_minute_usd", "(Deprecated) Mined coins per minute (converted to USD).", constLabels).Set(statsData.Data.USDPerMinute)
	util.NewGauge(registry, namespace, "miner", "income_minute_btc", "(Deprecated) Mined coins per minute (converted to BTC).", constLabels).Set(statsData.Data.BTCPerMinute)

	// Worker stats
	workerLabels := make(prometheus.Labels)
	workerLabels["worker"] = ""
	workerLastSeenMetric := util.NewGaugeVec(registry, namespace, "worker", "last_seen_seconds", "Delta between time of last statistics entry and when the miner was last seen (s).", constLabels, workerLabels)
	workerCurrentHashRateMetric := util.NewGaugeVec(registry, namespace, "worker", "hashrate_current_hps", "Current hash rate for a worker (H/s).", constLabels, workerLabels)
	workerValidSharesMetric := util.NewGaugeVec(registry, namespace, "worker", "shares_valid", "Number of valid shared for a worker.", constLabels, workerLabels)
	workerInvalidSharesMetric := util.NewGaugeVec(registry, namespace, "worker", "shares_invalid", "Number of invalid shared for a worker.", constLabels, workerLabels)
	workerStaleSharesMetric := util.NewGaugeVec(registry, namespace, "worker", "shares_stale", "Number of stale shared for a worker.", constLabels, workerLabels)
	for _, element := range workersData.Data {
		labels := make(prometheus.Labels)
		labels["worker"] = element.Name
		workerLastSeenMetric.With(labels).Set(element.Timestamp - element.LastSeenTimestamp)
		workerCurrentHashRateMetric.With(labels).Set(element.CurrentHashRate)
		workerValidSharesMetric.With(labels).Set(element.ValidShares)
		workerInvalidSharesMetric.With(labels).Set(element.InvalidShares)
		workerStaleSharesMetric.With(labels).Set(element.StaleShares)
	}

	return registry
}
