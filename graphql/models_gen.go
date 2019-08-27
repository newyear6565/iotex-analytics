// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package graphql

type Account struct {
	ActiveAccounts  []string         `json:"activeAccounts"`
	Alias           *Alias           `json:"alias"`
	OperatorAddress *OperatorAddress `json:"operatorAddress"`
}

type Action struct {
	ByDates *ActionList `json:"byDates"`
}

type ActionInfo struct {
	ActHash   string `json:"actHash"`
	BlkHash   string `json:"blkHash"`
	TimeStamp int    `json:"timeStamp"`
	ActType   string `json:"actType"`
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
	Amount    string `json:"amount"`
}

type ActionList struct {
	Exist   bool          `json:"exist"`
	Actions []*ActionInfo `json:"actions"`
	Count   int           `json:"count"`
}

type Alias struct {
	Exist     bool   `json:"exist"`
	AliasName string `json:"aliasName"`
}

type Bookkeeping struct {
	Exist              bool                  `json:"exist"`
	RewardDistribution []*RewardDistribution `json:"rewardDistribution"`
	Count              int                   `json:"count"`
}

type BucketInfo struct {
	VoterEthAddress string `json:"voterEthAddress"`
	WeightedVotes   string `json:"weightedVotes"`
}

type BucketInfoList struct {
	EpochNumber int           `json:"epochNumber"`
	BucketInfo  []*BucketInfo `json:"bucketInfo"`
	Count       int           `json:"count"`
}

type BucketInfoOutput struct {
	Exist          bool              `json:"exist"`
	BucketInfoList []*BucketInfoList `json:"bucketInfoList"`
}

type CandidateMeta struct {
	EpochNumber        int    `json:"epochNumber"`
	TotalCandidates    int    `json:"totalCandidates"`
	ConsensusDelegates int    `json:"consensusDelegates"`
	TotalWeightedVotes string `json:"totalWeightedVotes"`
	VotedTokens        string `json:"votedTokens"`
}

type Chain struct {
	MostRecentEpoch       int              `json:"mostRecentEpoch"`
	MostRecentBlockHeight int              `json:"mostRecentBlockHeight"`
	MostRecentTps         int              `json:"mostRecentTPS"`
	NumberOfActions       *NumberOfActions `json:"numberOfActions"`
}

type Contract struct {
	Hash      string `json:"hash"`
	Timestamp string `json:"timestamp"`
	From      string `json:"from"`
	To        string `json:"to"`
	Quantity  string `json:"quantity"`
}

type Delegate struct {
	Reward       *Reward           `json:"reward"`
	Productivity *Productivity     `json:"productivity"`
	Bookkeeping  *Bookkeeping      `json:"bookkeeping"`
	BucketInfo   *BucketInfoOutput `json:"bucketInfo"`
	Staking      *StakingOutput    `json:"staking"`
}

type DelegateAmount struct {
	DelegateName string `json:"delegateName"`
	Amount       string `json:"amount"`
}

type EpochRange struct {
	StartEpoch int `json:"startEpoch"`
	EpochCount int `json:"epochCount"`
}

type Hermes struct {
	Exist              bool                  `json:"exist"`
	HermesDistribution []*HermesDistribution `json:"hermesDistribution"`
}

type HermesDistribution struct {
	DelegateName        string                `json:"delegateName"`
	RewardDistribution  []*RewardDistribution `json:"rewardDistribution"`
	StakingIotexAddress string                `json:"stakingIotexAddress"`
	VoterCount          int                   `json:"voterCount"`
	WaiveServiceFee     bool                  `json:"waiveServiceFee"`
	Refund              string                `json:"refund"`
}

type NumberOfActions struct {
	Exist bool `json:"exist"`
	Count int  `json:"count"`
}

type OperatorAddress struct {
	Exist           bool   `json:"exist"`
	OperatorAddress string `json:"operatorAddress"`
}

type Pagination struct {
	Skip  int `json:"skip"`
	First int `json:"first"`
}

type Productivity struct {
	Exist              bool   `json:"exist"`
	Production         string `json:"production"`
	ExpectedProduction string `json:"expectedProduction"`
}

type Reward struct {
	Exist           bool   `json:"exist"`
	BlockReward     string `json:"blockReward"`
	EpochReward     string `json:"epochReward"`
	FoundationBonus string `json:"foundationBonus"`
}

type RewardDistribution struct {
	VoterEthAddress   string `json:"voterEthAddress"`
	VoterIotexAddress string `json:"voterIotexAddress"`
	Amount            string `json:"amount"`
}

type RewardSources struct {
	Exist                 bool              `json:"exist"`
	DelegateDistributions []*DelegateAmount `json:"delegateDistributions"`
}

type StakingInformation struct {
	EpochNumber  int    `json:"epochNumber"`
	TotalStaking string `json:"totalStaking"`
	SelfStaking  string `json:"selfStaking"`
}

type StakingOutput struct {
	Exist       bool                  `json:"exist"`
	StakingInfo []*StakingInformation `json:"stakingInfo"`
}

type Voting struct {
	VotingMeta    *VotingMeta    `json:"votingMeta"`
	RewardSources *RewardSources `json:"rewardSources"`
}

type VotingMeta struct {
	Exist         bool             `json:"exist"`
	CandidateMeta []*CandidateMeta `json:"candidateMeta"`
}
