package main

import (
	"tools/net"
	"tools/util"

	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/status-im/keycard-go/hexutils"
	"github.com/urfave/cli/v2"
)

var (
	txsCommand = cli.Command{
		Name:  "txs",
		Usage: "Query txs in the given net for the given account",
		Action: func(c *cli.Context) error {
			if c.NArg() < 2 {
				return errors.New("txs subcommand needs net and addr args")
			}
			var scanDomain string
			if strings.Compare("main", c.Args().Get(0)) == 0 {
				scanDomain = "apilist"
			} else if strings.Compare("nile", c.Args().Get(0)) == 0 {
				scanDomain = "nileapi"
			} else {
				return errors.New("error net argument")
			}
			addr := c.Args().Get(1)
			start := 0
			total := 20
			if c.NArg() == 3 {
				total, _ = strconv.Atoi(c.Args().Get(2))
			}
			if c.NArg() == 4 {
				start, _ = strconv.Atoi(c.Args().Get(2))
				total, _ = strconv.Atoi(c.Args().Get(3))
			}
			fmt.Println("[Legend]: ✅ - [Success] ⚠️  - [Revert] ⏱  - [Out_Of_Time] ⚡️ - [Out_Of_Energy] 💢 - [Other]")
			for i := 0; i < total; i += 50 {
				data := net.Get("https://" + scanDomain +
					".tronscan.org/api/transaction?" +
					"sort=-timestamp&" +
					"count=true&" +
					"limit=50" +
					"&start=" + strconv.Itoa(start+i) +
					"&address=" + addr)

				if data != nil {
					var txs Txs
					err := json.Unmarshal(data, &txs)
					if err != nil {
						return err
					}

					// fmt.Printf("[Total]: %5d\n", txs.Total)
					// total = txs.Total
					cache := make(map[string]string)
					for j, tx := range txs.Data {
						fmt.Printf("%"+strconv.Itoa(len(strconv.Itoa(total)))+"d %s %s %s ",
							i+j+1,
							time.Unix(tx.Timestamp/1000, 0).Format("2006-01-02 15:04:05"),
							tx.Hash,
							tx.OwnerAddress)
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
						if len(tx.TriggerInfo.Data) >= 8 {
							if _, ok := cache[tx.TriggerInfo.Data[:8]]; !ok {
								selector, _ := hex.DecodeString(tx.TriggerInfo.Data[:8])
								method := net.QueryMethod(selector)
								if len(method) != 0 {
									cache[tx.TriggerInfo.Data[:8]] = method
								} else {
									cache[tx.TriggerInfo.Data[:8]] = fmt.Sprintf("%x", selector)
								}
							}
							fmt.Print(cache[tx.TriggerInfo.Data[:8]])
						}
						fmt.Println()
					}
				}
			}
			return nil
		},
	}
	txCommand = cli.Command{
		Name:  "tx",
		Usage: "Query the tx in the given net for the given tx hash",
		Action: func(c *cli.Context) error {
			if c.NArg() < 1 {
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
				rspData := net.Post(url, reqData)
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
						fmt.Println("[Return data]:")
						fmt.Println("  - In HEX: " + hexutils.BytesToHex(data))
						if len(data) == 32 {
							fmt.Println("  - In INT: " + big.NewInt(0).SetBytes(data).String())
						}
						fmt.Println("  - In ASCII: " + utils.ToReadableASCII(data))
					}
				}
			}

			url = "https://" + scanDomain + ".tronscan.org/api/transaction-info?hash=" + hash
			rspData := net.Get(url)
			var scanTxInfo ScanTxInfo
			err := json.Unmarshal(rspData, &scanTxInfo)
			if err != nil {
				return err
			}

			// print some details in ScanTxInfo
			fmt.Println("[From]:", scanTxInfo.ContractData.OwnerAddress)
			fmt.Println("[To]:", scanTxInfo.ContractData.ContractAddress)

			var method string
			callData := hexutils.HexToBytes(scanTxInfo.ContractData.Data)
			// scan tx info does not contain method, so we query by 4byte
			if strings.Compare("()", scanTxInfo.TriggerInfo.Method) == 0 {
				// make sure calldata >= 4
				if len(callData) >= 4 {
					method = net.QueryMethod(callData[:4])
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
						printSol(r, &args[i].Type, "Arg", i, 1)
					}
				}
			} else if len(scanTxInfo.ContractData.Data) >= 8 {
				fmt.Println("[Selector]: " + scanTxInfo.ContractData.Data[:8])
				fmt.Println("[DataWord]:")
			} else {
				fmt.Println("[Selector]: none")
			}
			return nil
		},
	}
)

type Txs struct {
	Total int
	Data  []Tx
}

type Tx struct {
	OwnerAddress string `json:"ownerAddress"`
	ToAddress    string `json:"toAddress"`
	Hash         string `json:"hash"`
	Timestamp    int64  `json:"timestamp"`
	ContractRet  string `json:"contractRet"`
	TriggerInfo  struct {
		Data       string `json:"data"`
		MethodName string `json:"methodName"`
	} `json:"trigger_info"`
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
