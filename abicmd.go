package main

import (
    "tools/log"
    "tools/net"
    "tools/util"

    "encoding/json"
    "errors"
    "fmt"
    "math/big"
    "reflect"
    "regexp"
    "sort"
    "strconv"
    "strings"

    "github.com/btcsuite/btcutil/base58"
    "github.com/ethereum/go-ethereum/accounts/abi"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/common/math"
    "github.com/ethereum/go-ethereum/crypto"
    "github.com/status-im/keycard-go/hexutils"
    "github.com/urfave/cli/v2"
)

type Contract struct {
    Address string `json:"contract_address"`
    ABI     struct {
        Entries []map[string]interface{} `json:"entrys"`
    } `json:"abi"`
}

var (
    callCommand = cli.Command{
        Name:  "call",
        Usage: "Interact with contract on TRON network (main or nile)",
        Action: func(c *cli.Context) error {
            if c.NArg() < 2 {
                return errors.New("call command needs at least net and contract address")
            }
            domain := c.Args().Get(0)
            if strings.Compare("main", domain) == 0 {
                domain = "api"
            } else if strings.Compare("nile", domain) == 0 {
                domain = "nile"
            } else {
                return errors.New("wrong net arg (main or nile)")
            }
            contractAddr := c.Args().Get(1)
            abiAddr := c.Args().Get(1)
            if c.NArg() > 2 {
                abiAddr = c.Args().Get(2)
            }
            resData := net.Get(fmt.Sprintf("https://%s.trongrid.io/wallet/getcontract?value=%s&visible=true", domain, abiAddr))
            var contract Contract
            if err := json.Unmarshal(resData, &contract); err == nil {
                if len(contract.Address) == 0 {
                    return errors.New("contract not exist, you may input wrong net")
                }
                for _, abi := range contract.ABI.Entries {
                    if _, ok := abi["stateMutability"]; ok {
                        abi["stateMutability"] = strings.ToLower(abi["stateMutability"].(string))
                    }
                    if _, ok := abi["type"]; ok {
                        abi["type"] = strings.ToLower(abi["type"].(string))
                    }
                }
                data, _ := json.Marshal(contract.ABI.Entries)
                contractABI, _ := abi.JSON(strings.NewReader(string(data)))

                // first sort key
                var keys []string
                for k := range contractABI.Methods {
                    keys = append(keys, k)
                }
                sort.Strings(keys)

                // second get each method according the sorted keys
                var methods []abi.Method
                for _, k := range keys {
                    methods = append(methods, contractABI.Methods[k])
                    fmt.Printf("%2d. %s\n", len(methods), contractABI.Methods[k].Sig)
                }

                // next ask user to input the method index he wants to call
                for {
                    fmt.Print("Which method you want to call: ")
                    var index string
                    fmt.Scanln(&index)
                    i, err := strconv.Atoi(index)
                    if i <= 0 || i >= len(methods) || err != nil {
                        fmt.Println("Input index error, try again.")
                        continue
                    }
                    method := methods[i-1]
                    fmt.Printf("You choose method: [%s]\n", strings.ReplaceAll(method.String(), "function ", ""))
                    args := make([]interface{}, 0)
                    if len(method.Inputs) > 0 {
                        fmt.Println("Please input arguments:")
                        for i, inputType := range method.Inputs {
                            if len(inputType.Name) == 0 {
                                fmt.Printf(" - %d: ", i)
                            } else {
                                fmt.Printf(" - %s: ", inputType.Name)
                            }
                            var input string
                            fmt.Scanln(&input)
                            if arg, err := pack(inputType.Type, input); err == nil {
                                args = append(args, arg)
                            }
                        }
                    }
                    calldata, err := method.Inputs.Pack(args...)
                    if err == nil {
                        var from string
                        if !method.IsConstant() {
                            fmt.Print("Please input from address (default zero address): ")
                            fmt.Scanln(&from)
                        }
                        if _, ok := utils.ToAddress(from); !ok {
                            from = "T9yD14Nj9j7xAB4dbGeiX9h8unkKHxuWwb"
                        }
                        res := net.Trigger(domain, contractAddr, from, method.Sig, hexutils.BytesToHex(calldata))
                        // print energy used
                        fmt.Println("[Energy Used]\n  - " + strconv.Itoa(int(res.EnergyUsed)))
                        // print constant result
                        if len(res.ConstantResult) > 0 && len(res.ConstantResult[0]) > 0 {
                            fmt.Println("[Return Data]")
                            results := make(map[string]interface{})
                            err = method.Outputs.UnpackIntoMap(results, common.FromHex(res.ConstantResult[0]))
                            if err != nil {
                                fmt.Println(err.Error())
                            } else {
                                for k, v := range results {
                                    if len(k) == 0 {
                                        fmt.Printf("  - %v\n", v)
                                    } else {
                                        fmt.Printf("  - %s: %v\n", k, v)
                                    }
                                }
                            }
                        }
                        // print logs
                        if len(res.Logs) != 0 {
                            fmt.Println("[Logs]")
                        }
                        for _, log := range res.Logs {
                            fmt.Println(log)
                        }
                        // print internal transactions
                        if len(res.InternalTxs) != 0 {
                            fmt.Println("[Internal Txs]")
                        }
                        for _, tx := range res.InternalTxs {
                            fmt.Println(tx)
                        }
                    } else {
                        fmt.Printf("Pack error: %s\n", err.Error())
                    }
                }
            }
            return nil
        },
    }
    abiPadCommand = cli.Command{
        Name:  "pad",
        Usage: "Pad num(in hex or dec) to 32bytes",
        Action: func(c *cli.Context) error {
            if c.NArg() != 1 {
                return errors.New("pad subcommand needs num arg")
            }
            arg := c.Args().Get(0)
            // input is in hex
            if argBytes, ok := utils.FromHex(arg); ok {
                words := len(argBytes)/32 + 1
                log.NewLog("32bytes in BE", common.LeftPadBytes(argBytes, words*32))
                log.NewLog("32bytes in LE", common.RightPadBytes(argBytes, words*32))
            } else {
                // otherwise input must be in dec
                if num, ok := utils.FromDec(arg); ok {
                    log.NewLog("origin hex", num.Bytes())
                    log.NewLog("padded hex", common.LeftPadBytes(num.Bytes(), 32))
                } else {
                    return errors.New("input is in dec, but cannot convert it")
                }
            }
            return nil
        },
    }
    abiSplitCommand = cli.Command{
        Name:  "split",
        Usage: "Spilt data to each 32bytes",
        Action: func(c *cli.Context) error {
            if c.NArg() != 1 {
                return errors.New("split subcommand only needs data arg")
            }
            arg0 := c.Args().Get(0)
            data, ok := utils.FromHex(arg0)
            if !ok {
                return errors.New("only accept input in hex")
            }
            if len(data)%32 == 4 {
                method := net.QueryMethod(data[:4])
                if len(method) != 0 {
                    fmt.Printf("[selector]: %x - %s\n", data[:4], method)
                } else {
                    fmt.Printf("[selector]: %x\n", data[:4])
                }
                data = data[4:]
            }
            if len(data)%32 != 0 {
                return errors.New("data must be 32*N")
            }
            fmt.Println("[each data word]:")
            format := "0x%02x: %x\n"
            if len(data) > 8*32 {
                format = "0x%03x: %x\n"
            }
            for i := 0; i < len(data)/32; i++ {
                fmt.Printf(format, i*32, data[i*32:i*32+32])
            }
            return nil
        },
    }
    abiUnpackCommand = cli.Command{
        Name:  "unpack",
        Usage: "Unpack data with given types",
        Action: func(c *cli.Context) error {
            if c.NArg() != 2 {
                return errors.New("unpack subcommand needs data and type args")
            }
            arg0, arg1 := c.Args().Get(0), c.Args().Get(1)
            data, ok := utils.FromHex(arg1)
            if !ok {
                return errors.New("only accept data in hex")
            }
            args := abi.Arguments{}
            for _, arg := range strings.Split(arg0, ",") {
                solType, _ := abi.NewType(arg, "", nil)
                args = append(args, abi.Argument{Type: solType})
            }
            if res, err := args.UnpackValues(data); err == nil {
                fmt.Printf("[unpack result]:\n")
                for i, r := range res {
                    printSol(r, &args[i].Type, "arg", i, 1)
                }
                return nil
            } else {
                return err
            }
        },
    }
    abi4bytesCommand = cli.Command{
        Name:  "4bytes",
        Usage: "Get 4bytes selector for given method or event",
        Action: func(c *cli.Context) error {
            if c.NArg() != 1 {
                return errors.New("4bytes subcommand needs func or event signature arg")
            }
            signature := c.Args().Get(0)
            // drop name for all params
            nameRegExp := regexp.MustCompile(`\s\w+([,)])`)
            signature = nameRegExp.ReplaceAllString(signature, "$1")
            // expand all int|uint to int256|uint256
            abbIntRegExp := regexp.MustCompile(`int([,)\s\[])`)
            signature = abbIntRegExp.ReplaceAllString(signature, "int256$1")
            // drop all whitespaces
            signature = strings.ReplaceAll(signature, " ", "")
            log.NewLog("abi valid", signature)
            selector := crypto.Keccak256([]byte(signature))[:4]
            log.NewLog("origin hex", selector)
            res := make([]byte, 32-len(selector))
            selector = append(selector, res...)
            log.NewLog("padded hex", selector)
            return nil
        },
    }
)

func pack(t abi.Type, v string) (interface{}, error) {
    switch t.T {
    case abi.AddressTy:
        if addrBytes, ok := utils.ToAddress(v); ok {
            addr := new(common.Address)
            addr.SetBytes(addrBytes)
            return addr, nil
        }
        return nil, errors.New("can`t pack address type")
    case abi.IntTy, abi.UintTy:
        var num math.HexOrDecimal256
        err := num.UnmarshalText([]byte(v))
        if err == nil {
            return (*big.Int)(&num), nil
        }
        return nil, errors.New("can`t pack int/uint type")
    case abi.FixedBytesTy, abi.BytesTy:
        if argBytes, ok := utils.FromHex(v); ok {
            if t.GetType().Kind() == reflect.Array {
                fixedBytes := reflect.New(t.GetType()).Elem()
                reflect.Copy(fixedBytes, reflect.ValueOf(argBytes))
                return fixedBytes.Interface(), nil
            } else {
                return argBytes, nil
            }
        }
    case abi.StringTy:
        return v, nil
    case abi.BoolTy:
        if strings.Compare(v, "true") == 0 {
            return true, nil
        }
        if strings.Compare(v, "false") == 0 {
            return false, nil
        }
        return nil, errors.New("bool type must be true or false")
    case abi.SliceTy, abi.ArrayTy:
        return nil, nil
    default:
        return nil, errors.New(fmt.Sprintf("%s type is not supported", t.String()))
    }
    return nil, nil
}

func printSol(param interface{}, paramTy *abi.Type, name string, index, offset int) {
    printSeparator(offset, "  ", "", "- ")
    switch paramTy.T {
    case abi.ArrayTy:
        fmt.Printf("[%s-%02d]: %s\n", name, index, paramTy.String())
        paramArray := reflect.ValueOf(param)
        for i := 0; i < paramArray.Len(); i++ {
            printSol(paramArray.Index(i).Interface(), paramTy.Elem, "array", i, offset+1)
        }
    case abi.SliceTy:
        fmt.Printf("[%s-%02d]: %s\n", name, index, paramTy.String())
        paramSlice := reflect.ValueOf(param)
        for i := 0; i < paramSlice.Len(); i++ {
            printSol(paramSlice.Index(i).Interface(), paramTy.Elem, "slice", i, offset+1)
        }
    case abi.BytesTy, abi.FixedBytesTy:
        fmt.Printf("[%s-%02d]: %s, %#x\n", name, index, paramTy.String(), param)
    case abi.AddressTy:
        fmt.Printf("[%s-%02d]: %s, %v - %s\n", name, index, paramTy.String(), param, base58.CheckEncode(param.(common.Address).Bytes(), 0x41))
    case abi.IntTy, abi.UintTy:
        fmt.Printf("[%s-%02d]: %s, %v", name, index, paramTy.String(), param)
        intWithDot := formatBigInt(param.(*big.Int))
        if strings.ContainsAny(intWithDot, ",") {
            fmt.Printf(" - %s", intWithDot)
        }
        if len(param.(*big.Int).String()) >= 6 {
            fmt.Printf(" (%d)", len(param.(*big.Int).String()))
        }
        fmt.Println()
    default:
        fmt.Printf("[%s-%02d]: %s, %v\n", name, index, paramTy.String(), param)
        //fmt.Printf("[Parameter-%d]: %T, %#x\n", index, param, param)
    }
}

func printSeparator(repeat int, symbol, prefix, suffix string) {
    fmt.Print(prefix)
    for i := 0; i < repeat; i++ {
        fmt.Print(symbol)
    }
    fmt.Print(suffix)
}

func formatBigInt(n *big.Int) string {
    var (
        text  = n.String()
        buf   = make([]byte, len(text)+len(text)/3)
        comma = 0
        i     = len(buf) - 1
    )
    for j := len(text) - 1; j >= 0; j, i = j-1, i-1 {
        c := text[j]

        switch {
        case c == '-':
            buf[i] = c
        case comma == 3:
            buf[i] = ','
            i--
            comma = 0
            fallthrough
        default:
            buf[i] = c
            comma++
        }
    }
    return string(buf[i+1:])
}
