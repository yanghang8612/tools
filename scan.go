package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/status-im/keycard-go/hexutils"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var (
	txsCommand = cli.Command{
		Name:  "txs",
		Usage: "Query txs in the given net for the given account",
		Action: func(c *cli.Context) error {
			if c.NArg() < 2 {
				return errors.New("txs subcommand needs net and addr args")
			}
			var domain string
			if strings.Compare("main", c.Args().Get(0)) == 0 {
				domain = "apilist"
			} else if strings.Compare("nile", c.Args().Get(0)) == 0 {
				domain = "nileapi"
			} else {
				return errors.New("error net argument")
			}
			addr := c.Args().Get(1)
			start := "0"
			limit := "20"
			if c.NArg() > 2 {
				start = c.Args().Get(2)
			}
			if c.NArg() > 3 {
				limit = c.Args().Get(3)
			}
			data := doGet("https://" + domain +
				".tronscan.org/api/contracts/transaction?" +
				"sort=-timestamp&" +
				"count=true&" +
				"limit=" + limit +
				"&start=" + start +
				"&contract=" + addr)

			if data != nil {
				var txs Txs
				err := json.Unmarshal(data, &txs)
				if err != nil {
					return err
				}

				fmt.Printf("[Total]: %5d\n", txs.Total)
				fmt.Println("[Legend]: ✅ - [Success] ⚠️  - [Revert] ⏱  - [Out_Of_Time] ⚡️ - [Out_Of_Energy] 💢 - [Other]")
				cache := make(map[string]string)
				for i, tx := range txs.Data {
					fmt.Printf("%2d %s %s %s ", i+1, time.Unix(tx.Timestamp/1000, 0).Format("01-02 15:04:05"), tx.TxHash, tx.OwnAddress)
					switch tx.ContractRet {
					case "SUCCESS":
						fmt.Printf("✅ ")
					case "REVERT":
						fmt.Printf("⚠️  ")
					case "OUT_OF_TIME":
						fmt.Printf("⏱  ")
					case "OUT_OF_ENERGY":
						fmt.Printf("⚡️ ")
					default:
						fmt.Printf("💢 ")
					}
					if len(tx.CallData) >= 8 {
						if _, ok := cache[tx.CallData[:8]]; !ok {
							selector, _ := hex.DecodeString(tx.CallData[:8])
							method := queryMethod(selector)
							if len(method) != 0 {
								cache[tx.CallData[:8]] = method
							} else {
								cache[tx.CallData[:8]] = fmt.Sprintf("%x", selector)
							}
						}
						fmt.Print(cache[tx.CallData[:8]])
					}
					fmt.Println()
				}
			}
			return nil
		},
	}
	txCommand = cli.Command{
		Name:  "tx",
		Usage: "Query the tx in the given net for the given tx hash",
		Action: func(c *cli.Context) error {
			if c.NArg() < 2 {
				return errors.New("tx subcommand needs net and hash args")
			}
			var gridDomain, scanDomain string
			if strings.Compare("main", c.Args().Get(0)) == 0 {
				gridDomain = "api"
				scanDomain = "apilist"
			} else if strings.Compare("nile", c.Args().Get(0)) == 0 {
				gridDomain = "nile"
				scanDomain = "nileapi"
			} else {
				return errors.New("error net argument")
			}
			hash := c.Args().Get(1)
			url := "https://" + gridDomain + ".trongrid.io/wallet/gettransactioninfobyid"
			if reqData, err := json.Marshal(&TxHash{Value: hash}); err == nil {
				rspData := doPost(url, reqData)
				var gridTxInfo GridTxInfo
				err := json.Unmarshal(rspData, &gridTxInfo)
				if err != nil {
					return err
				}

				if len(gridTxInfo.ContractResult) != 0 {
					data, err := hex.DecodeString(gridTxInfo.ContractResult[0])
					if err != nil {
						return err
					}
					if len(data) == 0 {
						fmt.Println("[No return data]")
					} else {
						fmt.Println("[Result]:" + string(data))
						fmt.Println("In HEX:\t" + hexutils.BytesToHex(data))
						fmt.Println("In ASCII:\t" + string(data))
					}
				}
			}

			url = "https://" + scanDomain + ".tronscan.org/api/transaction-info?hash=" + hash
			rspData := doGet(url)
			var scanTxInfo ScanTxInfo
			err := json.Unmarshal(rspData, &scanTxInfo)
			if err != nil {
				return err
			}

			// print some details in ScanTxInfo
			fmt.Println("[From]: ", scanTxInfo.ContractData.OwnerAddress)
			fmt.Println("[To]: ", scanTxInfo.ContractData.ContractAddress)

			var method string
			callData := hexutils.HexToBytes(scanTxInfo.ContractData.Data)
			// scan tx info does not contain method, so we query by 4byte
			if strings.Compare("()", scanTxInfo.TriggerInfo.Method) == 0 {
				// make sure calldata >= 4
				if len(callData) >= 4 {
					method = queryMethod(callData[:4])
				}
			} else {
				method = scanTxInfo.TriggerInfo.Method
			}

			// we get the method signature, so try to abi.decode
			if len(method) != 0 {
				fmt.Println("[Method]: " + method)
				result := strings.FieldsFunc(method, func(r rune) bool {
					if r == '(' || r == ')' || r == ',' {
						return true
					}
					return false
				})
				args := abi.Arguments{}
				for _, param := range result[1:] {
					if strings.ContainsAny(param, " ") {
						param = strings.Split(param, " ")[0]
					}
					solType, _ := abi.NewType(param, "", nil)
					args = append(args, abi.Argument{Type: solType})
				}
				if res, err := args.UnpackValues(callData[4:]); err == nil {
					for i, r := range res {
						printSol(r, &args[i].Type, "Arg", i, 0)
					}
				}
			} else {
				fmt.Println("[Selector]: " + scanTxInfo.ContractData.Data[:8])
				fmt.Println("[DataWord]:")
			}
			return nil
		},
	}
)

type Txs struct {
	Total uint
	Data  []Tx
}

type Tx struct {
	OwnAddress  string
	ToAddress   string
	CallData    string `json:"call_data"`
	TxHash      string
	Timestamp   int64
	ContractRet string
}

type Rsp4Bytes struct {
	Count   uint
	Results []struct {
		Signature string `json:"text_signature"`
	}
}

type TxHash struct {
	Value string `json:"value"`
}

type GridTxInfo struct {
	ContractResult []string
}

type ScanTxInfo struct {
	ContractData struct {
		Data            string
		OwnerAddress    string `json:"owner_address"`
		ContractAddress string `json:"contract_address"`
	}
	TriggerInfo struct {
		Method    string
		parameter string
	} `json:"trigger_info"`
}

func doGet(url string) []byte {
	resp, err := http.Get(url)
	if err == nil && resp.StatusCode == 200 {
		defer resp.Body.Close()
		if body, err := ioutil.ReadAll(resp.Body); err == nil {
			return body
		}
	}
	return nil
}

func doPost(url string, data []byte) []byte {
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err == nil && resp.StatusCode == 200 {
		defer resp.Body.Close()
		if body, err := ioutil.ReadAll(resp.Body); err == nil {
			return body
		}
	}
	return nil
}

func queryMethod(selector []byte) string {
	var data []byte
	data = doGet(fmt.Sprintf("https://www.4byte.directory/api/v1/signatures/"+
		"?hex_signature=%x", selector))
	var rsp Rsp4Bytes
	err := json.Unmarshal(data, &rsp)
	if err == nil {
		if rsp.Count != 0 {
			return rsp.Results[rsp.Count-1].Signature
		}
	}
	return ""
}
