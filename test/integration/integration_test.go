package integration

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	tc "github.com/testcontainers/testcontainers-go/modules/compose"
	"github.com/yu-org/yu/common"
	"github.com/yu-org/yu/core/keypair"
	"github.com/yu-org/yu/example/client/callchain"
	"testing"
	"time"
	"uask-chain/types"
)

func startDockerCompose(t *testing.T) {
	compose, err := tc.NewDockerCompose("./docker-compose.yml")
	assert.NoError(t, err, "NewDockerComposeAPI()")
	t.Cleanup(func() {
		assert.NoError(t, compose.Down(context.Background(), tc.RemoveOrphans(true), tc.RemoveImagesLocal), "compose.Down()")
	})
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	assert.NoError(t, compose.Up(ctx, tc.Wait(true)), "compose.Up()")
}

var (
	askPub, askPriv         = keypair.GenSrKeyWithSecret([]byte("asker"))
	answerPub, answerPriv   = keypair.GenSrKeyWithSecret([]byte("answer"))
	commentPub, commentPriv = keypair.GenSrKeyWithSecret([]byte("comment"))
)

func TestQuestion(t *testing.T) {
	startDockerCompose(t)

	err := writeQuestion("AddQuestion", &types.QuestionAddRequest{
		Title:     "What is Uask",
		Content:   []byte("What is Uask, what can it do?"),
		Tags:      nil,
		Timestamp: time.Now().String(),
	})
	assert.NoError(t, err)
}

func TestAnswer(t *testing.T) {

}

func TestComment(t *testing.T) {

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
