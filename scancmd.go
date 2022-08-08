package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"
	utils "tools/util"

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
			var domain string
			if strings.Compare("main", c.Args().Get(0)) == 0 {
				domain = "apilist"
			} else if strings.Compare("nile", c.Args().Get(0)) == 0 {
				domain = "nileapi"
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
				data := doGet("https://" + domain +
					".tronscan.org/api/contracts/transaction?" +
					"sort=-timestamp&" +
					"count=true&" +
					"limit=50" +
					"&start=" + strconv.Itoa(start+i) +
					"&contract=" + addr)

				if data != nil {
					var txs Txs
					err := json.Unmarshal(data, &txs)
					if err != nil {
						return err
					}

					//fmt.Printf("[Total]: %5d\n", txs.Total)
					//total = txs.Total
					cache := make(map[string]string)
					for j, tx := range txs.Data {
						fmt.Printf("%"+strconv.Itoa(len(strconv.Itoa(total)))+"d %s %s %s ",
							i+j+1,
							time.Unix(tx.Timestamp/1000, 0).Format("2006-01-02 15:04:05"),
							tx.TxHash,
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
			rspData := doGet(url)
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
	speedCommand = cli.Command{
		Name:  "speed",
		Usage: "Query pool speed for given token",
		Action: func(c *cli.Context) error {
			if c.NArg() < 3 {
				return errors.New("speed subcommand needs net, hash and decimal args")
			}
			pool := c.Args().Get(0)
			token := c.Args().Get(1)
			totalStolen := big.NewInt(0)
			totalReward := big.NewInt(0)
			decFloat, _ := new(big.Float).SetString("1e" + c.Args().Get(2))
			decInteger, _ := decFloat.Int(new(big.Int))

			var preTransfer TokenTransfer
			var preSpeed *big.Int

			start := 0
			skip := 0
			total := 10_000_000
			for i := 0; i < total; i += 50 {
				data := doGet("https://apilist.tronscan.org/api/token_trc20/transfers?" +
					"sort=timestamp&count=true&limit=50" +
					"&start=" + strconv.Itoa(start+i) +
					"&fromAddress=" + pool +
					"&tokens=" + token +
					"&relatedAddress=" + pool)

				if data != nil {
					var transfers Transfers
					err := json.Unmarshal(data, &transfers)
					if err != nil {
						return err
					}

					//fmt.Printf("[Total]: %5d\n", txs.Total)
					total = transfers.Total
					for _, transfer := range transfers.TokenTransfers {
						reward, _ := big.NewInt(0).SetString(transfer.Quant, 10)
						totalReward = totalReward.Add(totalReward, reward)
						if preTransfer.BlockTs != 0 {
							diff := preTransfer.BlockTs - transfer.BlockTs
							if diff != 0 {
								speed, _ := big.NewInt(0).SetString(preTransfer.Quant, 10)
								speed = speed.Div(speed, big.NewInt(diff/1000))
								speed = speed.Mul(speed, big.NewInt(86400))
								speed = speed.Div(speed, decInteger)
								if skip == 0 && preSpeed != nil && preSpeed.Cmp(speed) != 0 {
									skip = 2
									stolen := big.NewInt(0).Abs(preSpeed.Sub(preSpeed, speed))
									stolen = stolen.Mul(stolen, decInteger)
									stolen = stolen.Div(stolen, big.NewInt(86400))
									stolen = stolen.Mul(stolen, big.NewInt(diff/1000))
									fmt.Printf("%s %s %d stolen - %d\n",
										time.Unix(preTransfer.BlockTs/1000, 0).Format("2006-01-02 15:04:05"),
										preTransfer.TransactionId,
										speed,
										stolen)
									totalStolen = totalStolen.Add(totalStolen, stolen)
								} else if skip != 0 {
									skip -= 1
								}
								//if preSpeed != nil && speed.Cmp(preSpeed) != 0 {
								//	fmt.Printf("%s %s cur - %d pre - %d\n",
								//		time.Unix(preTransfer.BlockTs/1000, 0).Format("2006-01-02 15:04:05"),
								//		preTransfer.TransactionId,
								//		speed,
								//		preSpeed)
								//}
								preSpeed = speed
							}
						}
						preTransfer = transfer
					}
				}
			}
			fmt.Printf("Total Stolen: %d\n", totalStolen.Div(totalStolen, decInteger))
			fmt.Printf("Total Reward: %d\n", totalReward.Div(totalReward, decInteger))
			return nil
		},
	}
	transferCommand = cli.Command{
		Name:  "transfer",
		Usage: "Query all transfer records for given pool",
		Action: func(c *cli.Context) error {
			if c.NArg() < 3 {
				return errors.New("speed subcommand needs net and hash args")
			}
			pool := c.Args().Get(0)
			token := c.Args().Get(1)
			decFloat, _ := new(big.Float).SetString("1e" + c.Args().Get(2))
			decInt, _ := decFloat.Int(new(big.Int))
			//totalOut := big.NewInt(0)
			start := 0
			total := 10_000_000
			//fmt.Print("Start to fetch out records")
			//for i := 0; i < total; i += 50 {
			//	data := doGet("https://apilist.tronscan.org/api/token_trc20/transfers?" +
			//		"sort=timestamp&count=true&limit=50" +
			//		"&start=" + strconv.Itoa(start+i) +
			//		"&fromAddress=" + pool +
			//		"&tokens=" + token +
			//		"&relatedAddress=" + pool)
			//
			//	if data != nil {
			//		var transfers Transfers
			//		err := json.Unmarshal(data, &transfers)
			//		if err != nil {
			//			return err
			//		}
			//
			//		//fmt.Printf("[Total]: %5d\n", txs.Total)
			//		total = transfers.Total
			//		for _, transfer := range transfers.TokenTransfers {
			//			amount, _ := big.NewInt(0).SetString(transfer.Quant, 10)
			//			totalOut = totalOut.Add(totalOut, amount)
			//		}
			//	}
			//	fmt.Print(".")
			//}
			//fmt.Print("\nFetch out records done.\n")
			totalIn := big.NewInt(0)
			start = 0
			total = 10_000_000
			fmt.Print("Start to fetch in records")
			for i := 0; i < total; i += 50 {
				data := doGet("https://apilist.tronscan.org/api/token_trc20/transfers?" +
					"sort=timestamp&count=true&limit=50" +
					"&start=" + strconv.Itoa(start+i) +
					"&toAddress=" + pool +
					"&tokens=" + token +
					"&relatedAddress=" + pool)

				if data != nil {
					var transfers Transfers
					err := json.Unmarshal(data, &transfers)
					if err != nil {
						return err
					}

					//fmt.Printf("[Total]: %5d\n", txs.Total)
					total = transfers.Total
					for _, transfer := range transfers.TokenTransfers {
						amount, _ := big.NewInt(0).SetString(transfer.Quant, 10)
						totalIn = totalIn.Add(totalIn, amount)
					}
				}
				fmt.Print(".")
			}
			fmt.Print("\nFetch in records done.\n")
			//fmt.Printf("Total out: %d\n", totalOut.Div(totalOut, decInteger))
			fmt.Printf("Total in: %d\n", totalIn.Div(totalIn, decInt))
			rewardRate := toUint64(query(pool, "rewardRate()", ""))
			lastUpdateTime := toUint64(query(pool, "lastUpdateTime()", ""))
			balance := toBigInt(query(token, "balanceOf(address)", toEthAddr(pool)))
			endTime := big.NewInt(1655812800)
			duration := endTime.Sub(endTime, big.NewInt(int64(lastUpdateTime)))
			rewardNotClaim := duration.Mul(duration, big.NewInt(int64(rewardRate)))
			fmt.Printf("Current token balance: %d\n", balance.Div(balance, decInt))
			fmt.Printf("Reward not claimed: %d\n", rewardNotClaim.Div(rewardNotClaim, decInt))
			return nil
		},
	}
)

type Transfers struct {
	Total          int
	TokenTransfers []TokenTransfer `json:"token_transfers"`
}

type TokenTransfer struct {
	TransactionId string `json:"transaction_id"`
	BlockTs       int64  `json:"block_ts"`
	Quant         string
	FromAddress   string `json:"from_address"`
	ToAddress     string `json:"to_address"`
}

type Txs struct {
	Total int
	Data  []Tx
}

type Tx struct {
	OwnerAddress string
	ToAddress    string
	CallData     string `json:"call_data"`
	TxHash       string
	Timestamp    int64
	ContractRet  string
	TriggerInfo  struct {
		MethodName string
	} `json:"trigger_info"`
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
