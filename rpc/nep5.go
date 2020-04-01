package rpc

import (
	"math/big"
	"math/rand"
	"squirrel/log"
	"time"
)

// ApplicationLogResponse is the struct of returning data from 'getapplicationlog' rpc call.
type ApplicationLogResponse struct {
	jsonRPCResponse
	Result *RawApplicationLogResult `json:"result"`
}

// RawApplicationLogResult is the inner struct of struct 'ApplicationLogResponse'.
type RawApplicationLogResult struct {
	TxID       string                       `json:"txid"`
	Executions []RawApplicationLogExecution `json:"executions"`
}

// RawApplicationLogExecution is the execution part of application log.
type RawApplicationLogExecution struct {
	Trigger       string             `json:"trigger"`
	Contract      string             `json:"contract"`
	VMState       string             `json:"vmstate"`
	GasConsumed   *big.Float         `json:"gas_consumed"`
	Stack         interface{}        `json:"stack"`
	Notifications []RawNotifications `json:"notifications"`
}

// RawNotifications is the inner struct of struct 'RawApplicationLogResult'.
type RawNotifications struct {
	Contract string    `json:"contract"`
	State    *RawState `json:"state"`
}

// GetApplicationLog returns application log of nep5 transaction.
func GetApplicationLog(blockIndex int, txID string) *RawApplicationLogResult {
	params := []interface{}{txID}
	const method = "getapplicationlog"
	args := getRPCRequestBody(method, params)

	respData := ApplicationLogResponse{}
	rpcCall(blockIndex, args, &respData)

	if respData.Result != nil {
		return respData.Result
	}

	retryTime := uint(0)
	delay := 0

	for {
		retryTime++
		if delay < 10*1000 {
			delay = rand.Intn(1<<retryTime) + 1000
		}

		log.Printf("Can not get application log of %s\n", txID)
		log.Printf("Delay for %d msecs and try to connect again. RetryTime=%d\n", delay, retryTime)

		time.Sleep(time.Duration(delay) * time.Millisecond)
		rpcCall(blockIndex, args, &respData)
		if respData.Result != nil {
			return respData.Result
		}
	}
}
