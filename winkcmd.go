package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/status-im/keycard-go/hexutils"
	"math/big"
	"os"
	"strings"
)

const (
	URL = "https://api.trongrid.io/wallet/triggerconstantcontract"
)

func wink01() {
	var winkList []Wink
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := input.Text()

		if line == "query" {
			break
		}

		splitArray := strings.Split(line, " ")
		winkList = append(winkList, Wink{Name: splitArray[0], Proxy: splitArray[1]})
	}

	for _, v := range winkList {
		queryWink(&v)
		//testSubmit(&v)
	}

	//var accList []string
	//input := bufio.NewScanner(os.Stdin)
	//for input.Scan() {
	//	line := input.Text()
	//
	//	if line == "query" {
	//		break
	//	}
	//
	//	accList = append(accList, line)
	//}
	//
	////queryTRX(accList)
	//queryWINK(accList)
}

func wink02(w *Wink) {
	//w.Aggregator = toTronAddr(query(w.Proxy, "aggregator()", ""))
	currentRound := toUint64(query(w.Proxy, "latestRound()", ""))

	res := query(w.Proxy, "submit(uint256,int256)",
		fmt.Sprintf("%X", uint256.NewInt(currentRound+1).Bytes32())+fmt.Sprintf("%X", uint256.NewInt(60000).Bytes32()))
	if len(res) != 0 {
		fmt.Println(w.Proxy + " submit failed by " + res)
	}
}

func queryTRX(addrs []string) {
	for _, addr := range addrs {
		repData := doGet("https://apilist.tronscan.org/api/accountv2?address=" + addr)
		var acc AccountV2
		err := json.Unmarshal(repData, &acc)
		if err == nil {
			fmt.Printf("%s trx balance is %d\n", addr, acc.Balance)
		}
	}
}

func queryWINK(addrs []string) {
	for _, addr := range addrs {
		ethAddr, _, _ := base58.CheckDecode(addr)
		fmt.Printf("%s wink balance is %d\n", addr,
			toUint64(query("TLa2f6VPqDgRE67v1736s7bJ8Ray5wYjU7", "balanceOf(address)",
				hexutils.BytesToHex(common.BytesToAddress(ethAddr).Hash().Bytes()))))
	}
}

func queryWink(w *Wink) {
	//fmt.Println("Start to query " + w.Proxy)

	w.Description = abiDecode("string",
		query(w.Proxy, "description()", "")).(string)

	w.Decimals = toUint64(query(w.Proxy, "decimals()", ""))
	w.Aggregator = toTronAddr(query(w.Proxy, "aggregator()", ""))
	w.ProxyOwner = toTronAddr(query(w.Proxy, "owner()", ""))

	oracles := abiDecode("address[]",
		query(w.Aggregator, "getOracles()", "")).([]common.Address)
	for _, oracle := range oracles {
		oracleTron := toTronAddr(hexutils.BytesToHex(oracle.Bytes()))
		w.Oracles = append(w.Oracles, oracleTron)
		adminTron := toTronAddr(query(w.Aggregator, "getAdmin(address)",
			hexutils.BytesToHex(oracle.Hash().Bytes())))
		w.Admins = append(w.Admins, adminTron)
	}

	w.AggregatorOwner = toTronAddr(query(w.Aggregator, "owner()", ""))
	w.AggregatorValidator = toTronAddr(query(w.Aggregator, "validator()", ""))
	w.PaymentAmount = toUint256(query(w.Aggregator, "paymentAmount()", ""))
	w.MaxSubmissionCount = toUint64(query(w.Aggregator, "maxSubmissionCount()", ""))
	w.MinSubmissionCount = toUint64(query(w.Aggregator, "minSubmissionCount()", ""))
	w.RestartDelay = toUint64(query(w.Aggregator, "restartDelay()", ""))
	w.Timeout = toUint64(query(w.Aggregator, "timeout()", ""))

	fmt.Println(w)
}

func query(addr, selector, param string) string {
	reqData, _ := json.Marshal(&Query{
		OwnerAddress:     "TGArstQjuME6fjBmEXVMdkGZNufxEDT6QB",
		ContractAddress:  addr,
		FunctionSelector: selector,
		Parameter:        param,
		Visible:          true,
	})
	rspData := doPost(URL, reqData)
	var result Response
	_ = json.Unmarshal(rspData, &result)
	if !result.RpcResult.TriggerResult {
		fmt.Println(addr + " trigger failed.")
	}
	if len(result.Result) > 0 {
		return result.Result[0]
	}
	return "no return"
}

func abiDecode(typeStr, dataStr string) interface{} {
	args := abi.Arguments{}
	solType, _ := abi.NewType(typeStr, "", nil)
	args = append(args, abi.Argument{Type: solType})
	resultList, _ := args.UnpackValues(hexutils.HexToBytes(dataStr))
	return resultList[0]
}

func toTronAddr(ethAddressHex string) string {
	return base58.CheckEncode(common.HexToAddress(ethAddressHex).Bytes(), 0x41)
}

func toEthAddr(tronAddress string) string {
	ethAddr, _, _ := base58.CheckDecode(tronAddress)
	return hexutils.BytesToHex(common.BytesToAddress(ethAddr).Hash().Bytes())
}

func toUint64(hexData string) uint64 {
	return uint256.NewInt(0).SetBytes(hexutils.HexToBytes(hexData)).Uint64()
}

func toUint256(hexData string) *uint256.Int {
	return uint256.NewInt(0).SetBytes(hexutils.HexToBytes(hexData))
}

func toBigInt(hexData string) *big.Int {
	return big.NewInt(0).SetBytes(hexutils.HexToBytes(hexData))
}

type Query struct {
	OwnerAddress     string `json:"owner_address"`
	ContractAddress  string `json:"contract_address"`
	FunctionSelector string `json:"function_selector"`
	Parameter        string `json:"parameter"`
	Visible          bool   `json:"visible"`
}

type Response struct {
	Result    []string `json:"constant_result"`
	RpcResult struct {
		TriggerResult bool `json:"result"`
	} `json:"result"`
}

type Wink struct {
	Name                string
	Proxy               string
	Aggregator          string
	Description         string
	Decimals            uint64
	ProxyOwner          string
	Oracles             []string
	Admins              []string
	PendingAdmin        []string
	AggregatorOwner     string
	AggregatorValidator string
	PaymentAmount       *uint256.Int
	MaxSubmissionCount  uint64
	MinSubmissionCount  uint64
	RestartDelay        uint64
	Timeout             uint64
}

func (w *Wink) String() string {
	return fmt.Sprintf("%s,%s,%s,%s,%d,%s,%s,%s,%s,%s,%d,%d,%d,%d,%d", w.Name, w.Proxy, w.Aggregator, w.Description,
		w.Decimals, w.ProxyOwner, w.Oracles, w.Admins, w.AggregatorOwner, w.AggregatorValidator,
		w.PaymentAmount, w.MaxSubmissionCount, w.MinSubmissionCount, w.RestartDelay, w.Timeout)
}

type AccountV2 struct {
	Balance uint64
}
