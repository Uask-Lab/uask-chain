package types

// DefaultReputation is every user have by default
const DefaultReputation = 1

// question
const (
	AddQuestionReputationNeed  = DefaultReputation
	VoteQuestionReputationNeed = 3

	UpVoteQuestionReputationIncrease = 2
	DownVoteQuestionReputationReduce = 2
)

// answer
const (
	AddAnswerReputationNeed  = 2
	VoteAnswerReputationNeed = 5

	PickUpAnswerReputationIncrease = 2
	UpVoteAnswerReputationIncrease = 2
	DownVoteAnswerReputationReduce = 2
)

// comment
const (
	AddCommentReputationNeed = 2
)
