package eth

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
)

// ResponseSyncing respons of eth_syncing
type ResponseSyncing struct {
	KnownStates   int64 `json:"knownStates"`
	PulledStates  int64 `json:"pulledStates"`
	StartingBlock int64 `json:"startingBlock"`
	CurrentBlock  int64 `json:"currentBlock"`
	HighestBlock  int64 `json:"highestBlock"`
}

// Syncing returns sync status or bool
//  - there seems 2 different responses
func (e *Ethereum) Syncing() (*ResponseSyncing, bool, error) {
	e.logger.Info("Syncing()")

	var (
		resIF  interface{}
		resMap map[string]string
	)

	err := e.client.CallContext(e.ctx, &resIF, "eth_syncing")
	if err != nil {
		return nil, false, errors.Wrap(err, "fail to call client.CallContext(eth_syncing)")
	}
	//log.Println(resMap)
	//log.Println(resMap["knownStates"])   //0x6621b
	//log.Println(resMap["pulledStates"])  //0x62d5e
	//log.Println(resMap["startingBlock"]) //0x0
	//log.Println(resMap["currentBlock"])  //0x38e85
	//log.Println(resMap["highestBlock"])  //0x5dcc9f
	//return nil, false, nil

	//booleanにキャストできるか
	bRes, ok := resIF.(bool)
	if !ok {
		//map型にキャストできるか=>できないので、取得し直したほうがよさそう
		//resMap, ok = resIF.(map[string]string)
		err := e.client.CallContext(e.ctx, &resMap, "eth_syncing")
		if err != nil {
			return nil, false, errors.Wrap(err, "fail to call client.CallContext(eth_syncing)")
		}

		//map型のそれぞれの値をint64に変換する
		knownStates, err := hexutil.DecodeBig(resMap["knownStates"])
		if err != nil {
			return nil, false, errors.New("response is invalid")
		}
		pulledStates, err := hexutil.DecodeBig(resMap["pulledStates"])
		if err != nil {
			return nil, false, errors.New("response is invalid")
		}
		startingBlock, err := hexutil.DecodeBig(resMap["startingBlock"])
		if err != nil {
			return nil, false, errors.New("response is invalid")
		}
		currentBlock, err := hexutil.DecodeBig(resMap["currentBlock"])
		if err != nil {
			return nil, false, errors.New("response is invalid")
		}
		highestBlock, err := hexutil.DecodeBig(resMap["highestBlock"])
		if err != nil {
			return nil, false, errors.New("response is invalid")
		}

		resSync := ResponseSyncing{
			KnownStates:   knownStates.Int64(),
			PulledStates:  pulledStates.Int64(),
			StartingBlock: startingBlock.Int64(),
			CurrentBlock:  currentBlock.Int64(),
			HighestBlock:  highestBlock.Int64(),
		}

		return &resSync, false, nil
	}
	return nil, bRes, nil
}
