package integration

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/yu-org/yu/common"
	"github.com/yu-org/yu/core/keypair"
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
	askPub, askPriv         = keypair.GenSrKeyWithSecret([]byte("asker"))
	answerPub, answerPriv   = keypair.GenSrKeyWithSecret([]byte("answer"))
	commentPub, commentPriv = keypair.GenSrKeyWithSecret([]byte("comment"))
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
	q1Content = []byte("What is Uask, what can it do?")
	q2Title   = "what is Ethereum"
	q2Content = []byte("What is Ethereum, what it blockchain")

	q1UpTitle   = "What is the Uask chain"
	q1UpContent = []byte("What can Uask do? how can I run it?")

	answer   = []byte("It is a question and answer appchain")
	answerUp = []byte("Uask is a question and answer appchain!")

	comment   = []byte("I agree with you")
	commentUp = []byte("I don't agree with you")
)

func TestUask(t *testing.T) {
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

	t.Run("DeleteQuestion", testDeleteQuestion)
	t.Run("DeleteAnswer", testDeleteAnswer)
	t.Run("DeleteComment", testDeleteComment)

	//sub.CloseSub()
	//stopDockerCompose()
}

func testAddQuestion(t *testing.T) {
	assert.NoError(t, writeQuestion("AddQuestion", &types.QuestionAddRequest{
		Title:     q1Title,
		Content:   q1Content,
		Timestamp: time.Now().Unix(),
	}))

	qid1 = getIdfromEvent(t, resultCh)

	assert.NoError(t, writeQuestion("AddQuestion", &types.QuestionAddRequest{
		Title:     q2Title,
		Content:   q2Content,
		Timestamp: time.Now().Unix(),
	}))

	qid2 = getIdfromEvent(t, resultCh)
}

func testListQuestions(t *testing.T) {
	bytes, err := readQuestion("ListQuestions", map[string]int{"pageSize": 2, "page": 1})
	assert.NoError(t, err)
	var qs []*types.QuestionInfo
	t.Logf("bytes = %s", bytes)
	assert.NoError(t, json.Unmarshal(bytes, &qs))
	assert.Equal(t, qs[0].ID, qid2)
	assert.Equal(t, qs[1].ID, qid1)
}

func testUpdateQuestion(t *testing.T) {
	assert.NoError(t, writeQuestion("UpdateQuestion", &types.QuestionUpdateRequest{
		ID: qid1,
		QuestionAddRequest: types.QuestionAddRequest{
			Title:     q1UpTitle,
			Content:   q1UpContent,
			Timestamp: time.Now().Unix(),
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
		QID:       qid1,
		Content:   answer,
		Timestamp: time.Now().Unix(),
	}))

	aid = getIdfromEvent(t, resultCh)
}

func testUpdateAnswer(t *testing.T) {
	assert.NoError(t, writeAnswer("UpdateAnswer", &types.AnswerUpdateRequest{
		ID: aid,
		AnswerAddRequest: types.AnswerAddRequest{
			QID:       qid1,
			Content:   answerUp,
			Timestamp: time.Now().Unix(),
		},
	}))
	dealResult(t, resultCh)
}

func testAddComment(t *testing.T) {
	assert.NoError(t, writeComment("AddComment", &types.CommentAddRequest{
		AID:       aid,
		Content:   comment,
		Timestamp: time.Now().Unix(),
	}))

	cid = getIdfromEvent(t, resultCh)
}

func testUpdateComment(t *testing.T) {
	assert.NoError(t, writeComment("UpdateComment", &types.CommentUpdateRequest{
		ID: cid,
		CommentAddRequest: types.CommentAddRequest{
			AID:       aid,
			Content:   commentUp,
			Timestamp: time.Now().Unix(),
		},
	}))
	dealResult(t, resultCh)
}

func testGetQuestion(t *testing.T) {
	qbyt, err := readQuestion("GetQuestion", map[string]string{"id": qid1})
	assert.NoError(t, err, "get question")
	q := new(types.QuestionInfo)
	assert.NoError(t, json.Unmarshal(qbyt, q))
	assert.Equal(t, q1UpTitle, q.Title)
}

func testGetAnswer(t *testing.T) {
	abyt, err := readAnswer("GetAnswer", map[string]string{"id": aid})
	assert.NoError(t, err, "get answer")
	a := new(types.AnswerInfo)
	assert.NoError(t, json.Unmarshal(abyt, a))
	assert.Equal(t, answerUp, a.Content)
}

func testGetComment(t *testing.T) {
	cbyt, err := readComment("GetComment", map[string]string{"id": cid})
	assert.NoError(t, err, "get comment")
	c := new(types.CommentInfo)
	assert.NoError(t, json.Unmarshal(cbyt, c))
	assert.Equal(t, commentUp, c.Content)
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
	return writeToUask("question", wrName, askPriv, askPub, params)
}

func writeAnswer(wrName string, params interface{}) error {
	return writeToUask("answer", wrName, answerPriv, answerPub, params)
}

func writeComment(wrName string, params interface{}) error {
	return writeToUask("comment", wrName, commentPriv, commentPub, params)
}

func writeToUask(tripodName, wrName string, priv keypair.PrivKey, pub keypair.PubKey, params interface{}) error {
	byt, err := json.Marshal(params)
	if err != nil {
		return err
	}
	callchain.CallChainByWriting(callchain.Http, priv, pub, &common.WrCall{
		TripodName:  tripodName,
		WritingName: wrName,
		Params:      string(byt),
	})
	return nil
}

func readQuestion(rdName string, params interface{}) ([]byte, error) {
	return readFromUask("question", rdName, params)
}

func readAnswer(rdName string, params interface{}) ([]byte, error) {
	return readFromUask("answer", rdName, params)
}

func readComment(rdName string, params interface{}) ([]byte, error) {
	return readFromUask("comment", rdName, params)
}

func readFromUask(tripodName, rdName string, params interface{}) ([]byte, error) {
	byt, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	return callchain.CallChainByReading(callchain.Http, &common.RdCall{
		TripodName:  tripodName,
		ReadingName: rdName,
		Params:      string(byt),
	}), nil
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
	assert.Equal(t, result.EventType, res.Type)
}
