package main

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/yu-org/yu/common"
	"github.com/yu-org/yu/core/keypair"
	"github.com/yu-org/yu/example/client/callchain"
	"os"
	"strings"
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

	strs := strings.Split(action, ".")
	tripod, writing = strs[0], strs[1]

	switch tripod {
	case "question":
		info := &types.QuestionAddRequest{
			Title:   titleOrId,
			Content: content,
			Tags:    nil,
		}
		params, err = json.Marshal(info)
		if err != nil {
			fmt.Println("marshal ask err: ", err)
			os.Exit(1)
		}
	case "answer":
		info := &types.AnswerAddRequest{
			QID:     titleOrId,
			Content: content,
		}
		params, err = json.Marshal(info)
		if err != nil {
			fmt.Println("marshal answer err: ", err)
			os.Exit(1)
		}
	case "comment":
		info := &types.CommentAddRequest{
			AID:     titleOrId,
			Content: content,
		}
		params, err = json.Marshal(info)
		if err != nil {
			fmt.Println("marshal comment err: ", err)
			os.Exit(1)
		}
	}

	callchain.CallChainByWriting(callchain.Http, priv, pub, &common.WrCall{
		TripodName:  tripod,
		WritingName: writing,
		Params:      string(params),
		LeiPrice:    0,
	})
}
