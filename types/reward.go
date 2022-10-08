package types

import "math/big"

type RewardRequest struct {
	// question id
	QID string `json:"qid"`
	// answer id -> reward
	Rewards map[string]*big.Int `json:"rewards"`
}

type RewardScheme struct {
	// question id
	QID string `json:"qid"`
	// answer id
	AID string `json:"aid"`

	Reward *big.Int
}
