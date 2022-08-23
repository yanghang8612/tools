package main

import (
    "tools/log"
    "tools/util"

    "errors"
    "math/big"
    "strconv"
    "strings"

    "github.com/btcsuite/btcutil/base58"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/common/math"
    "github.com/urfave/cli/v2"
)

var (
    hexAddrCommand = cli.Command{
        Name:  "addr",
        Usage: "Convert addr between hex, TRON-addr and eth-addr",
        Action: func(c *cli.Context) error {
            if c.NArg() != 1 {
                return errors.New("addr subcommand only needs one arg")
            }
            arg0 := c.Args().Get(0)
            addr := new(common.Address)
            if len(arg0) == 34 && arg0[0] == 'T' {
                // TRON address in Base58
                addrBytes, _, err := base58.CheckDecode(arg0)
                if err != nil {
                    return err
                }
                addr.SetBytes(addrBytes)
            } else if len(arg0) == 42 && strings.HasPrefix(arg0, "41") {
                // TRON address in hexadecimal with 0x41 prefix
                if addrBytes, ok := utils.FromHex(arg0[2:]); ok {
                    addr.SetBytes(addrBytes)
                } else {
                    return errors.New("input is begin with 0x41, but can`t convert to hex")
                }
            } else if bigint, ok := math.ParseBig256(arg0); ok {
                // input is in dec or hex
                addr.SetBytes(bigint.Bytes())
            } else {
                return errors.New("input is not recognized, please append 0x prefix if in hex")
            }
            log.NewLog("eth addr", addr.String())
            log.NewLog("tron addr", base58.CheckEncode(addr.Bytes(), 0x41))
            return nil
        },
    }
    hexIntCommand = cli.Command{
        Name:  "int",
        Usage: "Convert num between dec and hex",
        Action: func(c *cli.Context) error {
            if c.NArg() != 1 {
                return errors.New("int subcommand only needs num arg")
            }
            arg0 := c.Args().Get(0)
            if bigint, ok := math.ParseBig256(arg0); ok {
                log.NewLog("in hex", bigint.Bytes())
                log.NewLog("in dec", bigint)
                return nil
            } else {
                return errors.New("only accept input in dec or hex")
            }
        },
    }
    hexMaxCommand = cli.Command{
        Name:  "max",
        Usage: "Get max value for the type like uint-x",
        Action: func(c *cli.Context) error {
            if c.NArg() != 1 {
                return errors.New("max subcommand only needs size arg")
            }
            arg0 := c.Args().Get(0)
            size, err := strconv.Atoi(arg0)
            if err != nil {
                return err
            }
            if size <= 0 {
                return errors.New("input uint size should be greater than 0")
            }
            maxValue := new(big.Int)
            maxValue.Lsh(big.NewInt(1), uint(size))
            maxValue.Sub(maxValue, big.NewInt(1))
            log.NewLog("max hex", maxValue.Bytes())
            log.NewLog("max dec", maxValue)
            return nil
        },
    }
    hexStrCommand = cli.Command{
        Name:  "str",
        Usage: "convert hex between str",
        Action: func(c *cli.Context) error {
            if c.NArg() != 1 {
                return errors.New("hex command only needs single arg")
            }
            arg0 := c.Args().Get(0)
            // check if input is in hex
            if argBytes, ok := utils.FromHex(arg0); ok {
                if len(argBytes) <= 32 {
                    log.NewLog("in decimal", new(big.Int).SetBytes(argBytes))
                }
                // special case, first byte is `backspace`
                if len(argBytes) > 0 && argBytes[0] == 0x08 {
                    argBytes = argBytes[1:]
                }
                log.NewLog("in ascii", utils.ToReadableASCII(argBytes))
            } else {
                // otherwise treat input as str
                log.NewLog("in hex", []byte(arg0))
            }
            return nil
        },
    }
)
