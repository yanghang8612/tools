package main

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"tools/util"

	"github.com/status-im/keycard-go/hexutils"
	"github.com/urfave/cli/v2"
)

var (
	nowCommand = cli.Command{
		Name:  "now",
		Usage: "Convert time between datetime and timestamp",
		Action: func(c *cli.Context) error {
			if c.NArg() == 0 {
				fmt.Println(time.Now().Unix())
				fmt.Println(time.Now().UnixMilli())
				fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
			} else {
				arg := c.Args().Get(0)
				ts, err := strconv.Atoi(arg)
				if err == nil {
					fmt.Println(time.Unix(int64(ts), 0).Format("2006-01-02 15:04:05"))
				} else {
					loc, _ := time.LoadLocation("Asia/Shanghai")
					var dt time.Time
					var err error
					dt, err = time.ParseInLocation("2006-01-02 15:04:05", arg, loc)
					if err != nil {
						dt, err = time.ParseInLocation("2006-01-02", arg, loc)
					}
					if err != nil {
						fmt.Println("Date format err.")
					} else {
						fmt.Println(dt.Unix())
					}
				}
			}
			return nil
		},
	}

	hexCommand = cli.Command{
		Name:  "hex",
		Usage: "Convert num between decimal and hexadecimal",
		Action: func(c *cli.Context) error {
			if c.NArg() != 1 {
				return errors.New("hex command needs one arg")
			}
			arg := c.Args().Get(0)
			// input is num in hex
			if utils.ContainHexPrefix(arg) {
				argBytes := hexutils.HexToBytes(utils.DropHexPrefix(arg))
				fmt.Printf("[in decimal] - %d\n", new(big.Int).SetBytes(argBytes))
				// special case, first byte is `backspace`
				if len(argBytes) > 0 && argBytes[0] == 0x08 {
					argBytes = argBytes[1:]
				}
				fmt.Printf("[in ascii]   - %s\n", string(bytes.ToValidUTF8(argBytes, nil)))
			} else {
				// otherwise input must be in dec
				if num, ok := new(big.Int).SetString(arg, 10); ok {
					fmt.Printf("[in hex] - 0x%x\n", num.Bytes())
				} else {
					return errors.New("input type is in dec, but cannot covert it")
				}
			}
			return nil
		},
	}
)

func main() {
	app := cli.NewApp()
	app.Name = filepath.Base(os.Args[0])
	app.HideHelp = true
	app.Copyright = "Copyright 2021-2021 Asuka"
	app.Usage = "very useful tool kits for Asuka"
	app.CustomAppHelpTemplate = ""
	app.Commands = []*cli.Command{
		{
			Name:  "db",
			Usage: "Database related commands",
			Subcommands: []*cli.Command{
				&dbCountCommand,
				&dbGetCommand,
				&dbRootCommand,
				&dbPrintCommand,
				&dbDiffCommand,
			},
		},
		{
			Name:  "addr",
			Usage: "Address related commands",
			Subcommands: []*cli.Command{
				&addrNumCommand,
				&addrDecodeCommand,
				&addrEncodeCommand,
			},
		},
		{
			Name:  "hex",
			Usage: "Hex related commands",
			Subcommands: []*cli.Command{
				&addrCommand,
				&intCommand,
				&maxCommand,
				&strCommand,
			},
		},
		{
			Name:  "vm",
			Usage: "EVM related commands",
			Subcommands: []*cli.Command{
				&vmPadCommand,
				&vmSplitCommand,
				&vmUnpackCommand,
				&vm4bytesCommand,
			},
		},
		{
			Name:  "scan",
			Usage: "TronScan related commands",
			Subcommands: []*cli.Command{
				&txsCommand,
				&txCommand,
				&speedCommand,
				&transferCommand,
			},
		},
		{
			Name:  "eth",
			Usage: "ETH JSON-RPC related commands",
			Subcommands: []*cli.Command{
				&logsCommand,
			},
		},
		&nowCommand,
		&hexCommand,
	}

	for _, cmd := range app.Commands {
		cmd.HideHelpCommand = true
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
