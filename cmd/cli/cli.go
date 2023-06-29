package main

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/yu-org/yu/common"
	"github.com/yu-org/yu/core/keypair"
	"github.com/yu-org/yu/example/client/callchain"
	"os"
	"time"
	"uask-chain/types"
)

func main() {
	pub, priv := keypair.GenSrKeyWithSecret([]byte("uask-chain"))

	action := os.Args[1]
	titleOrId := os.Args[2]
	content := os.Args[3]

	logrus.SetLevel(logrus.DebugLevel)
	fmt.Printf("%s %s %s \n", action, titleOrId, content)

	var (
		tripod  string
		writing string
		params  []byte

		err error
	)

	//url := "localhost:5001"
	//hash, err := api.NewShell(url).Add(bytes.NewReader([]byte(content)))
	//if err != nil {
	//	panic(err)
	//}

	switch action {
	case "ask":
		info := &types.QuestionAddRequest{
			Title:       titleOrId,
			Content:     []byte(content),
			Tags:        nil,
			Timestamp:   time.Now().String(),
			Recommender: common.Address{},
		}
		params, err = json.Marshal(info)
		if err != nil {
			fmt.Println("marshal ask err: ", err)
			os.Exit(1)
		}
		tripod = "question"
		writing = "AddQuestion"
	case "answer":
		info := &types.AnswerAddRequest{
			QID:         titleOrId,
			Content:     []byte(content),
			Timestamp:   time.Now().String(),
			Recommender: common.Address{},
		}
		params, err = json.Marshal(info)
		if err != nil {
			fmt.Println("marshal answer err: ", err)
			os.Exit(1)
		}
		tripod = "answer"
		writing = "AddAnswer"
	case "comment":
		info := &types.CommentAddRequest{
			AID:       titleOrId,
			CID:       titleOrId,
			Content:   []byte(content),
			Timestamp: time.Now().String(),
		}
		params, err = json.Marshal(info)
		if err != nil {
			fmt.Println("marshal comment err: ", err)
			os.Exit(1)
		}
		tripod = "comment"
		writing = "AddComment"
	}

	callchain.CallChainByExec(callchain.Http, priv, pub, &common.WrCall{
		TripodName:  tripod,
		WritingName: writing,
		Params:      string(params),
		LeiPrice:    0,
	})
}
