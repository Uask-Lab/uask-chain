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

var (
	askPub, askPriv         = keypair.GenSrKeyWithSecret([]byte("asker"))
	answerPub, answerPriv   = keypair.GenSrKeyWithSecret([]byte("answer"))
	commentPub, commentPriv = keypair.GenSrKeyWithSecret([]byte("comment"))
)

func TestUask(t *testing.T) {
	startDockerCompose(t)

	resultCh := make(chan result.Result)
	go callchain.SubEvent(resultCh)

	// add question
	assert.NoError(t, writeQuestion("AddQuestion", &types.QuestionAddRequest{
		Title:     "What is Uask",
		Content:   []byte("What is Uask, what can it do?"),
		Timestamp: time.Now().String(),
	}))

	qid := getIdfromEvent(t, resultCh)

	// update question
	assert.NoError(t, writeQuestion("UpdateQuestion", &types.QuestionUpdateRequest{
		ID: qid,
		QuestionAddRequest: types.QuestionAddRequest{
			Title:     "What is the Uask",
			Content:   []byte("What can Uask do? how can I run it?"),
			Timestamp: time.Now().String(),
		},
	}))
	<-resultCh

	// search question
	resp, err := readQuestion("searchQuestion", map[string]string{"phrase": "Uask"})
	assert.NoError(t, err, "search question")
	t.Logf("search quesion result: %s", resp)

	// add answer
	assert.NoError(t, writeAnswer("AddAnswer", &types.AnswerAddRequest{
		QID:       qid,
		Content:   []byte("It is a question and answer appchain"),
		Timestamp: time.Now().String(),
	}))

	aid := getIdfromEvent(t, resultCh)

	// update answer
	assert.NoError(t, writeAnswer("UpdateAnswer", &types.AnswerUpdateRequest{
		ID: aid,
		AnswerAddRequest: types.AnswerAddRequest{
			QID:       qid,
			Content:   []byte("Uask is a question and answer appchain!"),
			Timestamp: time.Now().String(),
		},
	}))
	<-resultCh

	// add comment

}

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
	return callchain.CallChainByReading(callchain.Http, &common.Rdcall{
		TripodName:  tripodName,
		ReadingName: rdName,
		Params:      string(byt),
	}), nil
}

func getIdfromEvent(t *testing.T, resCh chan result.Result) string {
	res := <-resCh
	assert.Equal(t, result.EventType, res.Type())
	m := make(map[string]string)
	assert.NoError(t, res.(*result.Event).DecodeJsonValue(&m))
	return m["id"]
}
