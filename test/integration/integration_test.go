package integration

import (
	"crypto/ecdsa"
	"encoding/json"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/yu-org/yu/common"
	"github.com/yu-org/yu/core/result"
	"github.com/yu-org/yu/example/client/callchain"
	"os/exec"
	"testing"
	"time"
	"uask-chain/types"
)

func startDockerCompose(t *testing.T) {
	cmd := exec.Command("docker-compose", "up", "-d")
	assert.NoError(t, cmd.Run())
}

func stopDockerCompose() {
	exec.Command("docker-compose", "stop").Run()
}

var (
	askPriv, answerPriv, commentPriv *ecdsa.PrivateKey
)

var (
	qid1 string
	qid2 string
	aid  string
	cid  string

	resultCh = make(chan *result.Result)
)

var (
	q1Title   = "What is Uask"
	q1Content = "What is Uask, what can it do?"
	q2Title   = "what is Ethereum"
	q2Content = "What is Ethereum, what it blockchain"

	q1UpTitle   = "What is the Uask chain"
	q1UpContent = "What can Uask do? how can I run it?"

	answer   = "It is a question and answer appchain"
	answerUp = "Uask is a question and answer appchain!"

	comment   = "I agree with you"
	commentUp = "I don't agree with you"
)

func TestUask(t *testing.T) {
	var err error
	askPriv, err = crypto.HexToECDSA("80d7c5951de1ce731f90a8ee4e9221fa5075d41f199de56b4ab044089bd748af")
	assert.NoError(t, err)
	t.Logf("asker address: %x", crypto.PubkeyToAddress(askPriv.PublicKey))

	answerPriv, err = crypto.HexToECDSA("287d070e3c8df886998cdc76bf3a1cc6e171717f8b3546b9266e40fe3a8359c3")
	assert.NoError(t, err)
	t.Logf("answerer address: %x", crypto.PubkeyToAddress(answerPriv.PublicKey))

	commentPriv, err = crypto.HexToECDSA("2e16b4cc95fb40caa5295180b8dce8a72cc17ff1adb357a37b841eb6eac032c8")
	assert.NoError(t, err)
	t.Logf("commenter address: %x", crypto.PubkeyToAddress(commentPriv.PublicKey))

	startDockerCompose(t)

	time.Sleep(5 * time.Second)
	sub, err := callchain.NewSubscriber()
	assert.NoError(t, err)

	go sub.SubEvent(resultCh)

	t.Run("AddQuestion", testAddQuestion)
	t.Run("ListQuestions", testListQuestions)
	t.Run("UpdateQuestion", testUpdateQuestion)
	t.Run("GetQuestion", testGetQuestion)
	t.Run("SearchQuestion", testSearchQuestion)

	t.Run("AddAnswer", testAddAnswer)
	t.Run("UpdateAnswer", testUpdateAnswer)
	t.Run("GetAnswer", testGetAnswer)

	t.Run("AddComment", testAddComment)
	t.Run("UpdateComment", testUpdateComment)
	t.Run("GetComment", testGetComment)

	t.Run("testPickUp", testPickUp)
	t.Run("testDrop", testDrop)

	t.Run("testUpVoteQuestion", testUpVoteQuestion)
	t.Run("testDownVoteQuestion", testDownVoteQuestion)
	t.Run("testUpVoteAnswer", testUpVoteAnswer)
	t.Run("testDownVoteAnswer", testDownVoteAnswer)

	t.Run("DeleteQuestion", testDeleteQuestion)
	t.Run("DeleteAnswer", testDeleteAnswer)
	t.Run("DeleteComment", testDeleteComment)

	time.Sleep(5 * time.Second)

	//sub.CloseSub()
	//stopDockerCompose()
}

func testAddQuestion(t *testing.T) {
	assert.NoError(t, writeQuestion("AddQuestion", &types.QuestionAddRequest{
		Title:   q1Title,
		Content: q1Content,
	}))

	qid1 = getIdfromEvent(t, resultCh)

	assert.NoError(t, writeQuestion("AddQuestion", &types.QuestionAddRequest{
		Title:   q2Title,
		Content: q2Content,
	}))

	qid2 = getIdfromEvent(t, resultCh)
}

func testUpVoteQuestion(t *testing.T) {
	assert.NoError(t, voteQuestion("UpVote", map[string]string{"id": qid1}))
}

func testDownVoteQuestion(t *testing.T) {
	assert.NoError(t, voteQuestion("DownVote", map[string]string{"id": qid2}))
}

func testListQuestions(t *testing.T) {
	qs, err := readQuestion("ListQuestions", map[string]string{"pageSize": "2", "page": "1"})
	assert.NoError(t, err)
	assert.Equal(t, qs.([]any)[0].(map[string]any)["id"], qid2)
	assert.Equal(t, qs.([]any)[1].(map[string]any)["id"], qid1)
}

func testUpdateQuestion(t *testing.T) {
	assert.NoError(t, writeQuestion("UpdateQuestion", &types.QuestionUpdateRequest{
		ID: qid1,
		QuestionAddRequest: types.QuestionAddRequest{
			Title:   q1UpTitle,
			Content: q1UpContent,
		},
	}))
	dealResult(t, resultCh)
}

func testSearchQuestion(t *testing.T) {
	resp, err := readQuestion("SearchQuestion", map[string]string{"phrase": "Uask"})
	assert.NoError(t, err, "search question")
	t.Logf("search quesion result: %s", resp)
}

func testAddAnswer(t *testing.T) {
	assert.NoError(t, writeAnswer("AddAnswer", &types.AnswerAddRequest{
		QID:     qid1,
		Content: answer,
	}))

	aid = getIdfromEvent(t, resultCh)
}

func testUpVoteAnswer(t *testing.T) {
	assert.NoError(t, voteAnswer("UpVote", map[string]string{"id": aid}))
}

func testDownVoteAnswer(t *testing.T) {
	assert.NoError(t, voteAnswer("DownVote", map[string]string{"id": aid}))
}

func testPickUp(t *testing.T) {
	assert.NoError(t, writeToUask("answer", "PickUp", askPriv, map[string]string{"id": aid}))
}

func testDrop(t *testing.T) {
	assert.NoError(t, writeToUask("answer", "Drop", askPriv, map[string]string{"id": aid}))
}

func testUpdateAnswer(t *testing.T) {
	assert.NoError(t, writeAnswer("UpdateAnswer", &types.AnswerUpdateRequest{
		ID: aid,
		AnswerAddRequest: types.AnswerAddRequest{
			QID:     qid1,
			Content: answerUp,
		},
	}))
	dealResult(t, resultCh)
}

func testAddComment(t *testing.T) {
	assert.NoError(t, writeComment("AddComment", &types.CommentAddRequest{
		AID:     aid,
		Content: comment,
	}))

	cid = getIdfromEvent(t, resultCh)
}

func testUpdateComment(t *testing.T) {
	assert.NoError(t, writeComment("UpdateComment", &types.CommentUpdateRequest{
		ID: cid,
		CommentAddRequest: types.CommentAddRequest{
			AID:     aid,
			Content: commentUp,
		},
	}))
	dealResult(t, resultCh)
}

func testGetQuestion(t *testing.T) {
	q, err := readQuestion("GetQuestion", map[string]string{"id": qid1})
	assert.NoError(t, err, "get question")
	assert.Equal(t, q1UpTitle, q.(map[string]any)["title"])
}

func testGetAnswer(t *testing.T) {
	a, err := readAnswer("GetAnswer", map[string]string{"id": aid})
	assert.NoError(t, err, "get answer")
	assert.Equal(t, answerUp, a.(map[string]any)["content"])
}

func testGetComment(t *testing.T) {
	c, err := readComment("GetComment", map[string]string{"id": cid})
	assert.NoError(t, err, "get comment")
	assert.Equal(t, commentUp, c.(map[string]any)["content"])
}

func testDeleteQuestion(t *testing.T) {
	assert.NoError(t, writeQuestion("DeleteQuestion", map[string]string{"id": qid1}))
	dealResult(t, resultCh)
}

func testDeleteAnswer(t *testing.T) {
	assert.NoError(t, writeAnswer("DeleteAnswer", map[string]string{"id": aid}))
	dealResult(t, resultCh)
}

func testDeleteComment(t *testing.T) {
	assert.NoError(t, writeComment("DeleteComment", map[string]string{"id": cid}))
	dealResult(t, resultCh)
}

// helper funcs

func writeQuestion(wrName string, params interface{}) error {
	return writeToUask("question", wrName, askPriv, params)
}

func voteQuestion(wrName string, params any) error {
	return writeToUask("question", wrName, commentPriv, params)
}

func writeAnswer(wrName string, params interface{}) error {
	return writeToUask("answer", wrName, answerPriv, params)
}

func voteAnswer(wrName string, params any) error {
	return writeToUask("answer", wrName, commentPriv, params)
}

func writeComment(wrName string, params interface{}) error {
	return writeToUask("comment", wrName, commentPriv, params)
}

func writeToUask(tripodName, wrName string, priv *ecdsa.PrivateKey, params interface{}) error {
	byt, err := json.Marshal(params)
	if err != nil {
		return err
	}
	callchain.CallChainByWriting(priv, &common.WrCall{
		TripodName: tripodName,
		FuncName:   wrName,
		Params:     string(byt),
	})
	return nil
}

func readQuestion(rdName string, params map[string]string) (any, error) {
	return readFromUask("question", rdName, params)
}

func readAnswer(rdName string, params map[string]string) (any, error) {
	return readFromUask("answer", rdName, params)
}

func readComment(rdName string, params map[string]string) (any, error) {
	return readFromUask("comment", rdName, params)
}

func readUser(rdName string, params map[string]string) (any, error) {
	return readFromUask("user", rdName, params)
}

func readFromUask(tripodName, rdName string, params map[string]string) (any, error) {
	bytes := callchain.CallChainByReading(&common.RdCall{
		TripodName: tripodName,
		FuncName:   rdName,
	}, params)
	resp := new(types.Response)
	err := json.Unmarshal(bytes, resp)
	if err != nil {
		return nil, err
	}
	return resp.Payload, nil
}

func getIdfromEvent(t *testing.T, resCh chan *result.Result) string {
	res := <-resCh
	assert.Equal(t, result.EventType, res.Type)
	m := make(map[string]string)
	assert.NoError(t, res.Event.DecodeJsonValue(&m))
	return m["id"]
}

func dealResult(t *testing.T, resCh chan *result.Result) {
	res := <-resCh
	if res.Type == result.ErrorType {
		t.Error(res.String())
	} else {
		t.Log(res.String())
	}
}
