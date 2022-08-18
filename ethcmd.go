package main

import (
    utils "tools/util"

    "encoding/json"
    "errors"
    "fmt"
    "math/big"
    "strconv"
    "strings"

    "github.com/ethereum/go-ethereum/common/hexutil"
    "github.com/urfave/cli/v2"
)

var (
    logsCommand = cli.Command{
        Name:  "logs",
        Usage: "Query eth logs with given address, from block and topics, `page` logs at a query",
        Action: func(c *cli.Context) error {
            if c.NArg() != 4 {
                return errors.New("logs subcommand needs address, from block, topics and page args")
            }
            address := c.Args().Get(0)
            fromBlock, _ := strconv.Atoi(c.Args().Get(1))
            topics := strings.Split(c.Args().Get(2), ",")
            page, _ := strconv.Atoi(c.Args().Get(3))
            latestBlockNumber := int(getLatestBlockNumber())
            if latestBlockNumber != -1 {
                logs := make([]Log, 0)
                bar := utils.NewBar(0, latestBlockNumber-fromBlock)
                bar.Load()
                for i := fromBlock; i < latestBlockNumber; i += page + 1 {
                    var params []GetLogsParam
                    params = append(params, GetLogsParam{
                        Address:   address,
                        FromBlock: hexutil.EncodeUint64(uint64(i)),
                        ToBlock:   hexutil.EncodeUint64(uint64(i + page)),
                        Topics:    topics})

                    responseData := sendJsonRPCRequest("eth_getLogs", params)
                    var response GetLogsRPCResponse
                    err := json.Unmarshal(responseData, &response)
                    if err != nil {
                        return err
                    }

                    for _, log := range response.Result {
                        logs = append(logs, log)
                    }
                    bar.Add(page)
                }
                fmt.Print("\n")
                for _, log := range logs {
                    logData, _ := json.Marshal(log)
                    fmt.Println(string(logData))
                }
            }
            return nil
        },
    }
)

type RPCReq struct {
    JsonRPC string      `json:"jsonrpc"`
    Method  string      `json:"method"`
    Params  interface{} `json:"params"`
    Id      uint64      `json:"id"`
}

func sendJsonRPCRequest(method string, params interface{}) []byte {
    req := RPCReq{JsonRPC: "2.0",
        Method: method,
        Params: params,
        Id:     1}
    if reqData, err := json.Marshal(&req); err == nil {
        return doPost("http://47.90.254.215:8545/", reqData)
    }
    return nil
}

func getLatestBlockNumber() int64 {
    responseData := sendJsonRPCRequest("eth_blockNumber", nil)
    var response SingleStringRPCResponse
    err := json.Unmarshal(responseData, &response)
    if err != nil {
        return -1
    }

    blockNumberData, _ := hexutil.Decode(response.Result)
    return new(big.Int).SetBytes(blockNumberData).Int64()
}

type SingleStringRPCResponse struct {
    JsonRPC string `json:"jsonrpc"`
    Result  string `json:"result"`
    Id      uint64 `json:"id"`
}

type GetLogsParam struct {
    Address   string   `json:"address"`
    FromBlock string   `json:"fromBlock"`
    ToBlock   string   `json:"toBlock"`
    Topics    []string `json:"topics"`
}

type GetLogsRPCResponse struct {
    JsonRPC string `json:"jsonrpc"`
    Result  []Log  `json:"result"`
    Id      uint64 `json:"id"`
}

type Log struct {
    Address          string   `json:"address"`
    Topics           []string `json:"topics"`
    Data             string   `json:"data"`
    BlockNumber      string   `json:"blockNumber"`
    TransactionHash  string   `json:"transactionHash"`
    TransactionIndex string   `json:"transactionIndex"`
    BlockHash        string   `json:"blockHash"`
    LogIndex         string   `json:"logIndex"`
    Removed          bool     `json:"removed"`
}
