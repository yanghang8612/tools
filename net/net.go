package net

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/btcsuite/btcutil/base58"
	"github.com/ethereum/go-ethereum/common"
)

const (
	Endpoint    = "https://%s.trongrid.io/"
	TriggerPath = "wallet/triggerconstantcontract"
)

var appClient = &http.Client{
	Timeout: 3 * time.Second,
}

type TriggerRequest struct {
	OwnerAddress     string `json:"owner_address"`
	ContractAddress  string `json:"contract_address"`
	FunctionSelector string `json:"function_selector"`
	Parameter        string `json:"parameter"`
	Visible          bool   `json:"visible"`
}

type Log struct {
	Address string
	Data    string
	Topics  []string
}

func (l *Log) String() string {
	addrBytes, versionByte, _ := base58.CheckDecode(l.Address)
	combined := make([]byte, 0)
	combined = append(combined, versionByte)
	combined = append(combined, addrBytes...)
	tronAddr := base58.CheckEncode(combined, 0x41)
	topics := "[\n"
	for _, topic := range l.Topics {
		topics += "\t\t0x" + topic + ",\n"
	}
	topics += "\t]"
	return fmt.Sprintf("{\n\taddress: %s,\n\tdata: 0x%s,\n\ttopics: %s\n},", tronAddr, l.Data, topics)
}

type InternalTx struct {
	From string `json:"caller_address"`
	To   string `json:"transferTo_address"`
	Note string
}

func (tx *InternalTx) String() string {
	return fmt.Sprintf("{\n\tfrom: %s,\n\tto: %s,\n\ttype: %s\n}", tx.From, tx.To, string(common.FromHex(tx.Note)))
}

type TriggerResponse struct {
	Result struct {
		Result  bool
		Message string
	}
	EnergyUsed     uint64   `json:"energy_used"`
	ConstantResult []string `json:"constant_result"`
	Logs           []*Log
	InternalTxs    []*InternalTx `json:"internal_transactions"`
}

func Trigger(net, addr, from, selector, params string) *TriggerResponse {
	reqData, _ := json.Marshal(&TriggerRequest{
		OwnerAddress:     from,
		ContractAddress:  addr,
		FunctionSelector: selector,
		Parameter:        params,
		Visible:          true,
	})
	resData := Post(fmt.Sprintf(Endpoint, net)+TriggerPath, reqData)
	var triggerResponse TriggerResponse
	if err := json.Unmarshal(resData, &triggerResponse); err == nil {
		return &triggerResponse
	}
	return nil
}

type Rsp4Bytes struct {
	Count   uint
	Results []struct {
		Signature string `json:"text_signature"`
	}
}

func QueryMethod(selector []byte) string {
	var data []byte
	data = Get(fmt.Sprintf("https://www.4byte.directory/api/v1/signatures/?hex_signature=%x", selector))
	var rsp Rsp4Bytes
	err := json.Unmarshal(data, &rsp)
	if err == nil {
		if rsp.Count != 0 {
			return rsp.Results[rsp.Count-1].Signature
		}
	}

	data = Get(fmt.Sprintf("https://raw.githubusercontent.com/ethereum-lists/4bytes/master/signatures/%x", selector[:4]))
	if !strings.Contains(string(data), "404") {
		return string(data)
	}
	return ""
}

func Get(url string) []byte {
	resp, err := appClient.Get(url)
	if err == nil && resp.StatusCode == 200 {
		defer resp.Body.Close()
		if body, err := io.ReadAll(resp.Body); err == nil {
			return body
		}
	}
	return nil
}

func HighGet(url string, res interface{}) error {
	rspData := Get(url)
	err := json.Unmarshal(rspData, res)
	if err != nil {
		return err
	}
	return nil
}

func Post(url string, data []byte) []byte {
	resp, err := appClient.Post(url, "application/json", bytes.NewBuffer(data))
	if err == nil && resp.StatusCode == 200 {
		defer resp.Body.Close()
		if body, err := io.ReadAll(resp.Body); err == nil {
			return body
		}
	}
	return nil
}

func HighPost(url string, req interface{}, res interface{}) error {
	reqData, err := json.Marshal(&req)
	if err != nil {
		return err
	}
	rspData := Post(url, reqData)
	err = json.Unmarshal(rspData, res)
	if err != nil {
		return err
	}
	return nil
}
