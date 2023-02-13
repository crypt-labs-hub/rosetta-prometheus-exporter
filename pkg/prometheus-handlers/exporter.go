package prometheusexporter

import (
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"rosetta_exporter/pkg/config"
	rosettastatus "rosetta_exporter/pkg/rosetta-handlers"
)

var (
	blockchainInfo = prometheus.NewDesc(
		"blockchain_info",
		"blockchain info",
		nil,
		prometheus.Labels{
			"blockchain_name": "",
			"network_name":    "",
		},
	)

	rosettaInfo = prometheus.NewDesc(
		"rosetta_info",
		"Version of rosetta",
		nil, nil,
	)

	nodeInfo = prometheus.NewDesc(
		"node_info",
		"Version of the node",
		nil, nil,
	)

	currentBlockIndex = prometheus.NewDesc(
		"curr_block_index",
		"Index of the current block",
		nil, nil,
	)

	currentBlockTimestamp = prometheus.NewDesc(
		"curr_block_timestamp",
		"Timestamp of current block",
		nil, nil,
	)

	syncStatus = prometheus.NewDesc(
		"sync_status",
		"Sync Status",
		nil, nil,
	)
)

type Exporter struct {
	cfg *config.Config
}

func NewExporter(cfg *config.Config) *Exporter {
	return &Exporter{
		cfg: cfg,
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- blockchainInfo
	ch <- rosettaInfo
	ch <- nodeInfo
	ch <- currentBlockIndex
	ch <- currentBlockTimestamp
	ch <- syncStatus
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	// Get network status
	primaryNetwork, networkStatus, err := rosettastatus.GetStatus(e.cfg)
	if err != nil {

		log.Println(err)
		return
	}

	// set labels for the blockchain info
	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			"blockchain_info",
			"blockchain info",
			nil,
			prometheus.Labels{
				"blockchain_name": primaryNetwork.Blockchain,
				"network_name":    primaryNetwork.Network,
			},
		),
		prometheus.GaugeValue,
		1,
	)

	// set current block index
	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			"current_block_index",
			"current block index",
			nil,
			prometheus.Labels{
				"block_hash": networkStatus.CurrentBlockIdentifier.Hash,
			},
		),
		prometheus.GaugeValue,
		float64(networkStatus.CurrentBlockIdentifier.Index),
	)

	// set current block timestamp
	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			"current_block_timestamp",
			"current block timestamp",
			nil,
			nil,
		),
		prometheus.GaugeValue,
		float64(networkStatus.CurrentBlockTimestamp),
	)

	// convert sync status to integer
	syncStatus := 0
	if *networkStatus.SyncStatus.Synced {
		syncStatus = 1
	}

	// set current sync status
	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			"sync_status",
			"sync status",
			nil,
			nil,
		),
		prometheus.GaugeValue,
		float64(syncStatus),
	)

	// set current index
	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			"sync_current_index",
			"sync status: Current Index",
			nil,
			nil,
		),
		prometheus.GaugeValue,
		float64(*networkStatus.SyncStatus.CurrentIndex),
	)

	// set target index
	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			"sync_target_index",
			"sync status: Target Index",
			nil,
			nil,
		),
		prometheus.GaugeValue,
		float64(*networkStatus.SyncStatus.TargetIndex),
	)

}
