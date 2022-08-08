package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"tools/log"
	"tools/util"

	"github.com/urfave/cli/v2"
)

var (
	nowCommand = cli.Command{
		Name:  "now",
		Usage: "Convert time between datetime and timestamp",
		Action: func(c *cli.Context) error {
			// display now
			if c.NArg() == 0 {
				log.NewLog("in sec", time.Now().Unix())
				log.NewLog("in milli", time.Now().UnixMilli())
				log.NewLog("in datetime", time.Now())
			} else {
				arg := c.Args().Get(0)
				ts, err := strconv.Atoi(arg)
				// input str is timestamp
				if err == nil {
					log.NewLog("if sec", time.Unix(int64(ts), 0).Format("2006-01-02 15:04:05"))
					log.NewLog("if milli", time.Unix(int64(ts/1000), 0).Format("2006-01-02 15:04:05"))
				} else {
					// input str is date or time
					loc, _ := time.LoadLocation("Asia/Shanghai")
					formats := []string{"2006-01-02 15:04:05", "2006-01-02 15:04", "2006-01-02 15", "2006-01-02",
						"01-02 15:04:05", "01-02 15:04", "01-02", "15:04:05", "15:04"}
					for _, format := range formats {
						if dt, err := time.ParseInLocation(format, arg, loc); err == nil {
							if !strings.ContainsAny(format, "-") {
								dt = dt.AddDate(time.Now().Year(), int(time.Now().Month())-1, time.Now().Day())
							}
							if dt.Year() == 0 {
								dt = dt.AddDate(time.Now().Year(), 0, 0)
							}
							log.NewLog("in sec", dt.Unix())
							log.NewLog("in milli", dt.UnixMilli())
							log.NewLog("in datetime", dt)
						}
					}
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
	app.Copyright = "Copyright 2021-2022 Asuka, jeancky"
	app.Usage = "very useful tool kits for TRON and ethereum"
	app.CustomAppHelpTemplate = ""
	app.Action = func(c *cli.Context) error {
		if c.NArg() != 1 {
			return cli.ShowAppHelp(c)
		}
		arg := c.Args().Get(0)
		// arg is in hex
		if utils.ContainHexPrefix(arg) {
			data := utils.HexToBytes(arg)
			// 20 bytes, it should be eth addr
			if len(data) == 20 {
				return addrEncodeCommand.Action(c)
			}
			// 0 ~ 32 bytes, pad it to 32 bytes
			if len(data) < 32 {
				_ = vmPadCommand.Action(c)
			} else {
				// > 32 bytes, can be call data
				_ = vmSplitCommand.Action(c)
			}
			// output its value in decimal and try to display its readable string
			_ = hexStrCommand.Action(c)
		} else {
			// arg is not hex
			// check if arg is TRON addr
			if strings.HasPrefix(arg, "T") && len(arg) == 34 {
				return addrDecodeCommand.Action(c)
			}
			// check if arg is num in decimal
			num, err := strconv.Atoi(arg)
			if err == nil {
				if num >= 8 && num <= 256 && num%8 == 0 {
					_ = hexMaxCommand.Action(c)
				}
				// arg is decimal num
				_ = nowCommand.Action(c)
				// pad it to 32bytes
				_ = vmPadCommand.Action(c)
				// convert it to addr
				_ = addrNumCommand.Action(c)
			} else {
				// it may be a string
				// check if arg is date
				if strings.ContainsAny(arg, "-") || strings.ContainsAny(arg, ":") {
					return nowCommand.Action(c)
				}
				// check if arg is function or event
				if strings.ContainsAny(arg, "(") {
					return vm4bytesCommand.Action(c)
				}
				// otherwise, convert str to hex
				return hexStrCommand.Action(c)
			}
		}
		return nil
	}
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
				&hexAddrCommand,
				&hexIntCommand,
				&hexMaxCommand,
				&hexStrCommand,
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
	}
	app.After = func(c *cli.Context) error {
		log.FlushLogsToConsole()
		return nil
	}

	for _, cmd := range app.Commands {
		cmd.HideHelpCommand = true
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
