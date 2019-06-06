// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package chainmeta

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/iotexproject/iotex-analytics/indexprotocol"
	"github.com/iotexproject/iotex-analytics/indexprotocol/blocks"
	"github.com/iotexproject/iotex-analytics/indexservice"
	"github.com/iotexproject/iotex-analytics/queryprotocol/chainmeta/chainmetautil"
	s "github.com/iotexproject/iotex-analytics/sql"
)

// Protocol defines the protocol of querying tables
type Protocol struct {
	indexer *indexservice.Indexer
}

// ChainMeta defines chain meta
type ChainMeta struct {
	MostRecentEpoch       string
	MostRecentBlockHeight string
	MostRecentTps         string
}

type blkInfo struct {
	Transfer               int
	Execution              int
	DepositToRewardingFund int
	ClaimFromRewardingFund int
	GrantReward            int
	PutPollResult          int
	Timestamp              int
}

// NewProtocol creates a new protocol
func NewProtocol(idx *indexservice.Indexer) *Protocol {
	return &Protocol{indexer: idx}
}

// MostRecentTPS get most tps
func (p *Protocol) MostRecentTPS(ranges uint64) (tps int, err error) {
	_, ok := p.indexer.Registry.Find(blocks.ProtocolID)
	if !ok {
		err = errors.New("blocks protocol is unregistered")
		return
	}
	if ranges <= 0 {
		err = errors.Wrap(err, "TPS block window should be greater than 0")
		return
	}
	db := p.indexer.Store.GetDB()
	_, tipHeight, err := chainmetautil.GetCurrentEpochAndHeight(p.indexer.Registry, p.indexer.Store)
	if err != nil {
		err = errors.Wrap(err, "failed to get most recent block height")
		return
	}
	blockLimit := ranges
	if tipHeight < ranges {
		blockLimit = tipHeight
	}
	start := tipHeight - blockLimit + 1
	end := tipHeight
	getQuery := fmt.Sprintf("SELECT transfer,execution,depositToRewardingFund,claimFromRewardingFund,grantReward,putPollResult,timestamp FROM %s WHERE block_height>=? AND block_height<=?",
		blocks.BlockHistoryTableName)
	stmt, err := db.Prepare(getQuery)
	if err != nil {
		err = errors.Wrap(err, "failed to prepare get query")
		return
	}
	rows, err := stmt.Query(start, end)
	if err != nil {
		err = errors.Wrap(err, "failed to execute get query")
		return
	}
	var blk blkInfo
	parsedRows, err := s.ParseSQLRows(rows, &blk)
	if err != nil {
		err = errors.Wrap(err, "failed to parse results")
		return
	}
	if len(parsedRows) == 0 {
		err = indexprotocol.ErrNotExist
		return
	}
	var numActions int
	startTime := parsedRows[0].(*blkInfo).Timestamp
	endTime := parsedRows[0].(*blkInfo).Timestamp
	for _, parsedRow := range parsedRows {
		blk := parsedRow.(*blkInfo)
		numActions += blk.Transfer + blk.Execution + blk.ClaimFromRewardingFund + blk.DepositToRewardingFund + blk.GrantReward + blk.PutPollResult
		if blk.Timestamp > startTime {
			startTime = blk.Timestamp
		}
		if blk.Timestamp < endTime {
			endTime = blk.Timestamp
		}
	}
	timeDuration := startTime - endTime
	if timeDuration < 1 {
		timeDuration = 1
	}
	tps = numActions / timeDuration
	return
}

// GetChainMeta gets chain meta
func (p *Protocol) GetChainMeta(ranges int) (chainMeta *ChainMeta, err error) {
	currentEpoch, tipHeight, err := chainmetautil.GetCurrentEpochAndHeight(p.indexer.Registry, p.indexer.Store)
	if err != nil {
		err = errors.Wrap(err, "failed to get most recent block height")
		return
	}
	tps, err := p.MostRecentTPS(uint64(ranges))
	if err != nil {
		err = errors.Wrap(err, "failed to get most recent TPS")
		return
	}
	chainMeta = &ChainMeta{
		fmt.Sprintf("%d", currentEpoch),
		fmt.Sprintf("%d", tipHeight),
		fmt.Sprintf("%d", tps),
	}
	return
}

// GetNumberOfActions gets number of actions
func (p *Protocol) GetNumberOfActions(startEpoch uint64, epochCount uint64) (numberOfActions string, err error) {
	db := p.indexer.Store.GetDB()

	currentEpoch, _, err := chainmetautil.GetCurrentEpochAndHeight(p.indexer.Registry, p.indexer.Store)
	if err != nil {
		err = errors.Wrap(err, "failed to get current epoch")
		return
	}
	if startEpoch > currentEpoch {
		err = errors.New("start epoch should not be greater than current epoch")
		return
	}

	endEpoch := startEpoch + epochCount - 1
	getQuery := fmt.Sprintf("SELECT SUM(transfer)+SUM(execution)+SUM(depositToRewardingFund)+SUM(claimFromRewardingFund)+SUM(grantReward)+SUM(putPollResult) FROM %s WHERE epoch_number>=? and epoch_number<=?", blocks.BlockHistoryTableName)
	stmt, err := db.Prepare(getQuery)
	if err != nil {
		err = errors.Wrap(err, "failed to prepare get query")
		return
	}

	if err = stmt.QueryRow(startEpoch, endEpoch).Scan(&numberOfActions); err != nil {
		err = errors.Wrap(err, "failed to execute get query")
		return
	}
	return
}