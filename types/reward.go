package types

type RewardRequest struct {
	// question id
	QID string `json:"qid"`
	// answer id -> reward
	Rewards map[string]uint64 `json:"rewards"`
}

type RewardScheme struct {
	// question id
	QID string `json:"qid"`
	// answer id
	AID string `json:"aid"`

	Reward uint64
}
