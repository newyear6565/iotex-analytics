package hermes2

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/iotexproject/iotex-analytics/indexprotocol"
	"github.com/iotexproject/iotex-analytics/indexprotocol/accounts"
	"github.com/iotexproject/iotex-analytics/indexprotocol/actions"
	"github.com/iotexproject/iotex-analytics/indexprotocol/votings"
	"github.com/iotexproject/iotex-analytics/indexservice"
	s "github.com/iotexproject/iotex-analytics/sql"
)

const (
	// SelectCountByDelegateName selects the count of Hermes distribution by delegate name
	SelectCountByDelegateName = selectCount + fromJoinedTables + delegateFilter
	// SelectCountByVoterAddress selects the count of Hermes distribution by voter address
	SelectCountByVoterAddress = selectCount + fromJoinedTables + voterFilter

	fromJoinedTables = "FROM (SELECT * FROM %s WHERE epoch_number >= ? AND epoch_number <= ? AND `from` in (%s)) " +
		"AS t1 INNER JOIN (SELECT * FROM %s WHERE epoch_number >= ? AND epoch_number <= ?) AS t2 ON t1.action_hash = t2.action_hash "
	timeOrdering                           = "ORDER BY `timestamp` desc limit ?,?"
	fromTable                              = "FROM %s "
	selectVoter                            = "SELECT `to`, from_epoch, to_epoch, amount, t1.action_hash, `timestamp` "
	delegateFilter                         = "WHERE delegate_name = ? "
	selectHermesDistributionByDelegateName = selectVoter + fromJoinedTables + delegateFilter + timeOrdering
	delegateFilterWithEpochRange           = "WHERE delegate_name = ? AND epoch_number >= ? AND epoch_number <= ? "
	selectDelegate                         = "SELECT delegate_name, from_epoch, to_epoch, amount, t1.action_hash, `timestamp` "
	voterFilter                            = "WHERE `to` = ? "
	selectHermesDistributionByVoterAddress = selectDelegate + fromJoinedTables + voterFilter + timeOrdering
	selectDistributionRatio                = "SELECT block_reward_percentage AS block_reward_ratio, epoch_reward_percentage as epoch_reward_ratio, foundation_bonus_percentage as foundation_bonus_ratio, epoch_number "
	selectDistributionRatioByDelegateName  = selectDistributionRatio + fromTable + delegateFilterWithEpochRange
	selectCount                            = "SELECT COUNT(*),IFNULL(SUM(amount),0) "
	selectHermesMeta                       = "SELECT COUNT(DISTINCT delegate_name), COUNT(DISTINCT `to`), IFNULL(SUM(amount),0) " + fromJoinedTables
)

// HermesArg defines Hermes request parameters
type HermesArg struct {
	StartEpoch int
	EpochCount int
	Offset     uint64
	Size       uint64
}

// VoterInfo defines voter information
type VoterInfo struct {
	VoterAddress string
	FromEpoch    uint64
	ToEpoch      uint64
	Amount       string
	ActionHash   string
	Timestamp    string
}

// Ratio defines delegate reward distribution ratio
type Ratio struct {
	BlockRewardRatio     float64
	EpochRewardRatio     float64
	FoundationBonusRatio float64
	EpochNumber          int
}

// DelegateInfo defines delegate information
type DelegateInfo struct {
	DelegateName string
	FromEpoch    uint64
	ToEpoch      uint64
	Amount       string
	ActionHash   string
	Timestamp    string
}

// Protocol defines the protocol of querying tables
type Protocol struct {
	indexer      *indexservice.Indexer
	hermesConfig indexprotocol.HermesConfig
}

// NewProtocol creates a new protocol
func NewProtocol(idx *indexservice.Indexer, cfg indexprotocol.HermesConfig) *Protocol {
	return &Protocol{
		indexer:      idx,
		hermesConfig: cfg,
	}
}

// GetHermes2ByDelegate gets Hermes voter list by delegate name
func (p *Protocol) GetHermes2ByDelegate(arg HermesArg, delegateName string) ([]*VoterInfo, error) {
	db := p.indexer.Store.GetDB()
	getQuery := fmt.Sprintf(selectHermesDistributionByDelegateName, accounts.BalanceHistoryTableName, strings.Join(wrapperQueryValue(p.hermesConfig.MultiSendContractAddressList), ","), actions.HermesContractTableName)
	stmt, err := db.Prepare(getQuery)
	if err != nil {
		return nil, errors.Wrap(err, "failed to prepare get query")
	}
	defer stmt.Close()
	endEpoch := arg.StartEpoch + arg.EpochCount - 1
	rows, err := stmt.Query(arg.StartEpoch, endEpoch, arg.StartEpoch, endEpoch,
		delegateName, arg.Offset, arg.Size)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute get query")
	}

	var voterInfo VoterInfo
	parsedRows, err := s.ParseSQLRows(rows, &voterInfo)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse results")
	}
	if len(parsedRows) == 0 {
		return nil, indexprotocol.ErrNotExist
	}

	voterInfoList := make([]*VoterInfo, 0)
	for _, parsedRow := range parsedRows {
		voterInfoList = append(voterInfoList, parsedRow.(*VoterInfo))
	}

	return voterInfoList, nil
}

// GetHermes2ByVoter gets Hermes delegate list by voter name
func (p *Protocol) GetHermes2ByVoter(arg HermesArg, voterAddress string) ([]*DelegateInfo, error) {
	db := p.indexer.Store.GetDB()
	getQuery := fmt.Sprintf(selectHermesDistributionByVoterAddress, accounts.BalanceHistoryTableName, strings.Join(wrapperQueryValue(p.hermesConfig.MultiSendContractAddressList), ","), actions.HermesContractTableName)
	stmt, err := db.Prepare(getQuery)
	if err != nil {
		return nil, errors.Wrap(err, "failed to prepare get query")
	}
	defer stmt.Close()

	endEpoch := arg.StartEpoch + arg.EpochCount - 1
	rows, err := stmt.Query(arg.StartEpoch, endEpoch, arg.StartEpoch, endEpoch,
		voterAddress, arg.Offset, arg.Size)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute get query")
	}

	var delegateInfo DelegateInfo
	parsedRows, err := s.ParseSQLRows(rows, &delegateInfo)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse results")
	}
	if len(parsedRows) == 0 {
		return nil, indexprotocol.ErrNotExist
	}

	delegateInfoList := make([]*DelegateInfo, 0)
	for _, parsedRow := range parsedRows {
		delegateInfoList = append(delegateInfoList, parsedRow.(*DelegateInfo))
	}

	return delegateInfoList, nil
}

// GetHermes2Ratio gets Hermes distribution ratio list by delegate name
func (p *Protocol) GetHermes2Ratio(arg HermesArg, delegateName string) ([]*Ratio, error) {

	db := p.indexer.Store.GetDB()
	getQuery := fmt.Sprintf(selectDistributionRatioByDelegateName, votings.VotingResultTableName)
	stmt, err := db.Prepare(getQuery)
	if err != nil {
		return nil, errors.Wrap(err, "failed to prepare get query")
	}
	defer stmt.Close()

	endEpoch := arg.StartEpoch + arg.EpochCount - 1
	rows, err := stmt.Query(delegateName, arg.StartEpoch, endEpoch)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute get query")
	}

	var distributionRatioInfo Ratio
	parsedRows, err := s.ParseSQLRows(rows, &distributionRatioInfo)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse results")
	}
	if len(parsedRows) == 0 {
		return nil, indexprotocol.ErrNotExist
	}
	distributionRatioList := make([]*Ratio, 0)
	for _, parsedRow := range parsedRows {
		distributionRatioList = append(distributionRatioList, parsedRow.(*Ratio))
	}
	return distributionRatioList, nil
}

// GetHermes2Count gets the count of Hermes distributions
func (p *Protocol) GetHermes2Count(arg HermesArg, selectQuery string, filter string) (count int, total string, err error) {
	db := p.indexer.Store.GetDB()
	getQuery := fmt.Sprintf(selectQuery, accounts.BalanceHistoryTableName, strings.Join(wrapperQueryValue(p.hermesConfig.MultiSendContractAddressList), ","), actions.HermesContractTableName)
	stmt, err := db.Prepare(getQuery)
	if err != nil {
		err = errors.Wrap(err, "failed to prepare get query")
		return
	}
	defer stmt.Close()

	endEpoch := arg.StartEpoch + arg.EpochCount - 1
	if err = stmt.QueryRow(arg.StartEpoch, endEpoch, arg.StartEpoch, endEpoch,
		filter).Scan(&count, &total); err != nil {
		err = errors.Wrap(err, "failed to execute get query")
		return
	}
	return
}

// GetHermes2Meta gets the hermes meta info
func (p *Protocol) GetHermes2Meta(startEpoch int, epochCount int) (numberOfDelegates int,
	numberOfRecipients int, totalRewardsDistributed string, err error) {
	endEpoch := startEpoch + epochCount - 1
	db := p.indexer.Store.GetDB()
	getQuery := fmt.Sprintf(selectHermesMeta, accounts.BalanceHistoryTableName, strings.Join(wrapperQueryValue(p.hermesConfig.MultiSendContractAddressList), ","), actions.HermesContractTableName)
	stmt, err := db.Prepare(getQuery)
	if err != nil {
		err = errors.Wrap(err, "failed to prepare get query")
		return
	}
	defer stmt.Close()
	if err = stmt.QueryRow(startEpoch, endEpoch, startEpoch, endEpoch).
		Scan(&numberOfDelegates, &numberOfRecipients, &totalRewardsDistributed); err != nil {
		err = errors.Wrap(err, "failed to execute get query")
		return
	}
	return
}

func wrapperQueryValue(queryValues []string) []string {
	ret := make([]string, len(queryValues))
	for index, str := range queryValues {
		ret[index] = "'" + str + "'"
	}
	return ret
}
