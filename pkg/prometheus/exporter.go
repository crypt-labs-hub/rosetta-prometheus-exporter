package prometheusexporter

import (
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"rosetta_exporter/pkg/config"
	"rosetta_exporter/pkg/rosetta"
	"time"
)

type Exporter struct {
	// cfg configuration passed to the exporter
	cfg *config.Config
	// rh rosetta handler to connect with the rosetta endpoints
	rh *rosettahandlers.RosettaHandler
	// blockchainInfo the details of the blockchain including name and network
	blockchainInfo *prometheus.Desc
	// rosettaInfo the details about rosetta including version
	rosettaInfo *prometheus.Desc
	// nodeInfo node level details including version
	nodeInfo *prometheus.Desc
	// currentBlockIndex the latest block available at the node
	currentBlockIndex *prometheus.Desc
	// currentBlockTimestamp the timestamp on the latest block on the node
	currentBlockTimestamp *prometheus.Desc
	// syncStatus the difference between current block and max available block
	syncStatus *prometheus.Desc
}

// NewExporter constructor for exporter
// initializes every descriptor and returns a pointer to the exporter
func NewExporter(cfg *config.Config) *Exporter {
	// Get the rosetta handler
	rh, err := rosettahandlers.NewRosettaHandler(cfg)
	if err != nil {
		log.Panicf("Rosetta handler initialize error::%s", err)
		return nil
	}
	// return pointer to exporter
	return &Exporter{
		cfg: cfg,
		rh:  rh,
		blockchainInfo: prometheus.NewDesc(
			"blockchain_info",
			"blockchain info",
			nil,
			prometheus.Labels{
				"blockchain_name": "",
				"network_name":    "",
			},
		),
		rosettaInfo: prometheus.NewDesc(
			"rosetta_info",
			"Version of rosetta",
			nil, nil,
		),
		nodeInfo: prometheus.NewDesc(
			"node_info",
			"Version of the node",
			nil, nil,
		),
		currentBlockIndex: prometheus.NewDesc(
			"curr_block_index",
			"Index of the current block",
			nil, nil,
		),
		currentBlockTimestamp: prometheus.NewDesc(
			"curr_block_timestamp",
			"Timestamp of current block",
			nil, nil,
		),
		syncStatus: prometheus.NewDesc(
			"sync_status",
			"Sync Status",
			nil, nil,
		),
	}
}

// Describe set descriptors for all prometheus desc channels
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.blockchainInfo
	ch <- e.rosettaInfo
	ch <- e.nodeInfo
	ch <- e.currentBlockIndex
	ch <- e.currentBlockTimestamp
	ch <- e.syncStatus
}

// Collect implment the required collection function for all collectors
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	// Get network status
	networkStatus, err := e.rh.GetStatus()
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
				"blockchain_name": e.rh.PrimaryNetwork.Blockchain,
				"network_name":    e.rh.PrimaryNetwork.Network,
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

	// Get current block
	block, err := e.rh.GetBlock(networkStatus.CurrentBlockIdentifier)
	if err != nil {
		log.Println(err)
		return
	}

	// set number of txs
	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			"block_tx_count",
			"block transaction count",
			nil,
			prometheus.Labels{
				"block_hash": block.BlockIdentifier.Hash,
			},
		),
		prometheus.GaugeValue,
		float64(len(block.Transactions)),
	)

	// Get blocks for sample size
	blockCount := 0
	txCount := 0
	earliestTimestamp := time.UnixMilli(block.Timestamp)
	currTimestamp := time.UnixMilli(block.Timestamp)
	parentBlockId := block.ParentBlockIdentifier
	sampleSize := e.cfg.GetSampleSize()
	for blockCount < sampleSize {
		// Get parent block
		parentBlock, err := e.rh.GetBlock(parentBlockId)
		if err != nil {
			log.Println(err)
			return
		}
		// change to parent block
		parentBlockId = parentBlock.ParentBlockIdentifier
		// set timestamp
		earliestTimestamp = time.UnixMilli(parentBlock.Timestamp)
		// increase tx count
		txCount += len(parentBlock.Transactions)
		// increment block count
		blockCount++
	}

	// calculate time between max and min block in sample
	blockTimeRangeInSec := currTimestamp.Sub(earliestTimestamp).Seconds()

	// calculate block rate
	blockRate := float64(blockCount) / blockTimeRangeInSec

	// set block rate
	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			"block_rate",
			"block rate (blocks/sec)",
			nil,
			nil,
		),
		prometheus.GaugeValue,
		blockRate,
	)

	// TODO: tx endpoint might have other_txs that are not included
	// calculate tx rate
	txRate := float64(txCount) / blockTimeRangeInSec

	// set block rate
	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			"transaction_rate",
			"transaction rate(tx/sec)",
			nil,
			nil,
		),
		prometheus.GaugeValue,
		txRate,
	)

}
