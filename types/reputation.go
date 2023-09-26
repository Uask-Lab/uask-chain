package types

// DefaultReputation is every user have by default
const DefaultReputation = 1

// question
const (
	VoteQuestionReputationNeed       = 3
	UpVoteQuestionReputationIncrease = 2
	DownVoteQuestionReputationReduce = 2
)

// answer
const (
	VoteAnswerReputationNeed       = 5
	PickUpAnswerReputationIncrease = 2
	UpVoteAnswerReputationIncrease = 2
	DownVoteAnswerReputationReduce = 2
)
